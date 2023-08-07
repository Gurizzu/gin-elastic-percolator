package urequest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/logrusorgru/aurora"
)

func HTTPRequest(method, url, bodyParams string) (bodyResp []byte, err error) {
	client := http.Client{}
	log.Println(aurora.Cyan(method), aurora.Yellow(url))

	var bodyParamsBuffer *bytes.Buffer = nil
	if bodyParams != "" {
		log.Println(aurora.Yellow(bodyParams))
	}
	bodyParamsBuffer = bytes.NewBuffer([]byte(bodyParams))
	request, err := http.NewRequest(method, url, bodyParamsBuffer)

	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		log.Println("error status code not 200")
		err = fmt.Errorf("error %v from request", resp.StatusCode)
		return
	}

	bodyResp, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
