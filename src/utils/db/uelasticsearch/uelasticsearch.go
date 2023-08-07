package uelasticsearch

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type ElasticClient struct {
	addr string
	conf elasticsearch.Config
}

func NewESConnection(addr, user, pass string) *ElasticClient {
	o := ElasticClient{}
	o.conf = elasticsearch.Config{
		Addresses: []string{addr},
		Username:  user,
		Password:  pass,
		//Logger: &estransport.ColorLogger{
		//	Output: io.Discard,
		//},
	}
	o.addr = addr
	return &o
}

func NewESConnection_ENV() *ElasticClient {
	return NewESConnection(
		"http://localhost:9200/",
		"",
		"",
		//fmt.Sprintf("http://%s", os.Getenv(config.ENV_ES_ADDR)), os.Getenv(config.ENV_ES_USER), os.Getenv(config.ENV_ES_PASS),
	)
}

func (o *ElasticClient) connect() (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(o.conf)
}

func (o *ElasticClient) Insert(index, body string) (bytesBody []byte, err error) {
	esClient, err := o.connect()
	if err != nil {
		log.Println(err)
		return
	}

	timestamp := time.Now().UnixNano()
	hashId := md5.Sum([]byte(fmt.Sprintf("%d%s", timestamp, body)))

	//var buf bytes.Buffer
	//_ = json.NewEncoder(&buf).Encode(body)
	res, err := esClient.Index(
		index,
		strings.NewReader(body),
		esClient.Index.WithDocumentID(hex.EncodeToString(hashId[:])),
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if bytesBody, err = io.ReadAll(res.Body); err != nil {
		log.Println(err)
		return
	}

	return

}

func (o *ElasticClient) Search(index, body string) (bytesBody []byte, err error) {
	if index == "" && !gjson.Get(body, "pit").Exists() {
		err = errors.New("please input index or pit id")
		return
	}

	esClient, err := o.connect()
	if err != nil {
		log.Println(err)
		return
	}

	esSearch := esClient.Search
	listEsReq := []func(*esapi.SearchRequest){
		esSearch.WithContext(context.Background()),
		esSearch.WithPretty(),
	}

	if index != "" {
		listEsReq = append(listEsReq, esSearch.WithIndex(index))
		body, _ = sjson.Delete(body, "pit")
	} else {
		body, _ = sjson.Delete(body, "from")
	}

	if body != "" {
		listEsReq = append(listEsReq, esSearch.WithBody(strings.NewReader(body)))
	}

	if !gjson.Get(body, "search_after").Exists() {
		log.Println(aurora.Sprintf(
			aurora.Cyan("Requesting to ES %s \nGET %s/_search \n%s"),
			o.addr, index, body,
		))
	}

	res, err := esClient.Search(listEsReq...)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	if bytesBody, err = io.ReadAll(res.Body); err != nil {
		log.Println(err)
		return
	}

	if res.IsError() {
		errObj := gjson.ParseBytes(bytesBody)
		// Print the response status and error information.
		log.Printf("[%s] %s: %s\n",
			res.Status(),
			errObj.Get("error.type").String(),
			errObj.Get("error.reason").String(),
		)
		err = errors.New(errObj.Get("error.type").String())
		return
	}

	return
}

func (o *ElasticClient) OpenPIT(index string) (pitID string, err error) {
	esClient, err := o.connect()
	if err != nil {
		log.Println(err)
		return
	}

	var esOptions []func(*esapi.OpenPointInTimeRequest)
	esPIT := esClient.OpenPointInTime
	esOptions = append(esOptions,
		esPIT.WithContext(context.Background()),
	)

	res, err := esClient.OpenPointInTime([]string{index}, "5m", esOptions...)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	var bytesBody []byte
	if bytesBody, err = io.ReadAll(res.Body); err != nil {
		log.Println(err)
		return
	}

	if res.IsError() {
		errObj := gjson.ParseBytes(bytesBody)
		// Print the response status and error information.
		log.Printf("[%s] %s: %s\n",
			res.Status(),
			errObj.Get("error.type").String(),
			errObj.Get("error.reason").String(),
		)
		err = errors.New(errObj.Get("error.type").String())
		return
	}

	pitID = gjson.GetBytes(bytesBody, "id").String()

	return
}

func (o *ElasticClient) ClosePIT(pitID string) (err error) {
	esClient, err := o.connect()
	if err != nil {
		log.Println(err)
		return
	}

	jsonPit := gjson.Parse(pitID)
	if !jsonPit.IsObject() {
		pitID = fmt.Sprintf(`{"id": "%s"}`, pitID)
	}

	esPIT := esClient.ClosePointInTime
	esOptions := []func(*esapi.ClosePointInTimeRequest){
		esPIT.WithContext(context.Background()),
		esPIT.WithBody(strings.NewReader(pitID)),
	}

	res, err := esClient.ClosePointInTime(esOptions...)
	defer res.Body.Close()

	var bytesBody []byte
	if bytesBody, err = io.ReadAll(res.Body); err != nil {
		log.Println(err)
		return
	}

	if res.IsError() {
		errObj := gjson.ParseBytes(bytesBody)
		// Print the response status and error information.
		log.Printf("[%s] %s: %s\n",
			res.Status(),
			errObj.Get("error.type").String(),
			errObj.Get("error.reason").String(),
		)
		err = errors.New(errObj.Get("error.type").String())
		return
	}

	return
}
