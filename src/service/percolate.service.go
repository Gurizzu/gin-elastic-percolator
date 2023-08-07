package service

import (
	"encoding/json"
	"fmt"
	"gin-elastic-percolator/src/config"
	"gin-elastic-percolator/src/model"
	"gin-elastic-percolator/src/utils/db/uelasticsearch"
	"gin-elastic-percolator/src/utils/mapping"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

type PercolateService struct {
	esUtil *uelasticsearch.ElasticClient
}

func NewPercolateService() *PercolateService {
	fmt.Println(fmt.Sprintf("http://%s", os.Getenv(config.ENV_ES_ADDR)), os.Getenv(config.ENV_ES_USER), os.Getenv(config.ENV_ES_PASS))
	return &PercolateService{
		esUtil: uelasticsearch.NewESConnection_ENV(),
	}
}

func arrayToJSONString(arr []string) string {
	strArr := make([]string, len(arr))
	for i, item := range arr {
		strArr[i] = fmt.Sprintf(`"%s"`, item)
	}
	return "[" + strings.Join(strArr, ",") + "]"
}

func (o *PercolateService) AddQuery(param model.Percolate) (resp model.Response) {
	query := `{
	  "query": {
	    "query_string": {
	      "query": "%s"
	    }
	  },
      "tags" : %s
	}`

	queryJSON := fmt.Sprintf(query, param.Query, arrayToJSONString(param.Tags))
	fmt.Println(queryJSON)
	bytesRes, err := o.esUtil.Insert("news-percolate", queryJSON)
	if err != nil {
		log.Println(aurora.Red(err))
		return
	}

	resp.Data = bytesRes
	return

}

func (o *PercolateService) GetPercolate_Data(filename string) (data model.Percolate_Data, errMessage string) {
	query, err := mapping.GetPercolate(filename)
	if err != nil {
		log.Println(aurora.Red(err))
		errMessage = "FAIL"
		return
	}

	bytesRes, err := o.esUtil.Search("news-percolate", query)
	if err != nil {
		log.Println(aurora.Red(err))
		errMessage = "FAIL"
	}
	//var res interface{}
	_ = json.Unmarshal([]byte(gjson.GetBytes(bytesRes, "aggregations").String()), &data.Result)

	return
}
