package mapping

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

type PercolateDocument struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetPercolate(filename string) (queryJSON string, err error) {

	filePath := "public/" + filename
	if _, err = os.Stat(filePath); err != nil {
		filePath = "1000-dataima.json"
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var bulks []PercolateDocument

	byteFile, _ := io.ReadAll(file)

	result := gjson.ParseBytes(byteFile)
	result.ForEach(func(_, value gjson.Result) bool {
		title := value.Get("title").String()
		content := value.Get("content").String()

		bulks = append(bulks, PercolateDocument{
			Title:   title,
			Content: content,
		})
		return true // continue iteration
	})

	query := `{
	  "size": 0,
	  "query": {
	    "percolate": {
	      "field": "query",
	      "documents": %s
	    }
	  },
	  "aggs": {
	    "tags": {
	      "terms": {
	        "field": "tags.keyword",
	        "size": 10
	      }
	    }
	  },
		"highlight" : {
		 "fields" : {
			"content":{}
}
}
	}`

	documentsJSON, err := json.Marshal(bulks)
	if err != nil {
		log.Println(err)
		return
	}
	queryJSON = fmt.Sprintf(query, string(documentsJSON))
	return
}
