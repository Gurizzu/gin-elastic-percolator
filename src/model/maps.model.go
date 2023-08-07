package model

import (
	"encoding/json"
	"fmt"

	"git.blackeye.id/Aldi.Rismawan/centrotil/db/umongo"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type LayerMaps struct {
	MetadataWithID `bson:",inline"`

	DatasetId      string         `json:"datasetId" bson:"datasetId"`
	LocationConfig LocationConfig `json:"locationConfig" bson:"locationConfig"`
	Filters        []Filters      `json:"filters" bson:"filters"`
	LayerConfig    LayerConfig    `json:"layerConfig" bson:"layerConfig"`
}

type LocationConfig struct {
	LocationLevel  string           `json:"locationLevel" bson:"locationLevel"`
	LocationFilter []LocationFilter `json:"locationFilter" bson:"locationFilter"`
}

type LocationFilter struct {
	Level string `json:"level" bson:"level"`
	Value string `json:"value" bson:"value"`
	Field string `json:"field" bson:"field"`
	Alias string `json:"alias" bson:"alias"`
}

type Filters struct {
	Field string `json:"field" bson:"field"`
	Value string `json:"value" bson:"value"`
}

type Tooltips struct {
	Field string `json:"field" bson:"field"`
	Alias string `json:"alias" bson:"alias"`
}

type LayerConfig struct {
	Name          string     `json:"name" bson:"name"`
	LocationField string     `json:"locationField" bson:"locationField"`
	Type          string     `json:"type" bson:"type"`
	ValueField    string     `json:"valueField" bson:"valueField"`
	Tooltips      []Tooltips `json:"tooltips" bson:"tooltips"`
	Color         string     `json:"color" bson:"color"`
	Opacity       float32    `json:"opacity" bson:"opacity"`
	Icon          string     `json:"icon" bson:"icon"`
}

type DatasetMaps struct {
	Name      string `json:"name" bson:"name"`
	TableName string `json:"tableName" bson:"tableName"`
}

type LayerMaps_View struct {
	LayerMaps   `bson:",inline"`
	DatasetMaps DatasetMaps `json:"dataset" bson:"dataset"`
}

type LayerMaps_Data struct {
	JoinKey string   `json:"join_key"`
	Series  []Series `json:"series"`
	Props   []Props  `json:"props"`
}

type Series struct {
	Key      string            `json:"key"`
	Value    float64           `json:"value"`
	Tooltips *map[string]int64 `json:"tooltips,omitempty"`
}

type Props struct {
	ID           string                 `json:"_id"`
	Lan          float64                `json:"lan"`
	Lon          float64                `json:"lon"`
	CityCode     string                 `json:"city_code"`
	Geometry     map[string]interface{} `json:"geometry"`
	Properties   map[string]interface{} `json:"properties"`
	ProvinceCode string                 `json:"province_code"`
	Type         string                 `json:"type"`
}

type Req_LayerMaps_Data struct {
	LayerID string `json:"id" form:"id"`
	Year    int    `json:"year" form:"year"`
}

type Req_AllLayerMaps struct {
	Search     string `json:"search" form:"search"`
	DatasetId  string `json:"datasetId" form:"datasetId"`
	OPDId      string `json:"opdId" form:"opdId"`
	RegionId   string `json:"regionId" form:"regionId"`
	CategoryId string `json:"categoryId" form:"categoryId"`
	umongo.Request
}

func (o Req_AllLayerMaps) ToString() string {
	bytesData, _ := json.Marshal(o)
	return string(bytesData)
}

func (o LayerMaps) ToString() string {
	bytesData, _ := json.Marshal(o)
	return string(bytesData)
}

func (o LayerMaps) ToESBody() (request string) {
	request = fmt.Sprintf(`{"size":0,"query":{"bool":{"must":[{"match_all":{}}]}},"aggs":{"location":{"terms":{"field":"%s","size":100}}}}`, o.LayerConfig.LocationField)

	if valueField := o.LayerConfig.ValueField; valueField != "" {
		request, _ = sjson.Set(request, "aggs.location.aggs.value", gjson.Parse(fmt.Sprintf(`{"sum":{"field":"%s"}}`, valueField)).Value())
	}

	for _, tooltip := range o.LayerConfig.Tooltips {
		request, _ = sjson.Set(request, fmt.Sprintf("aggs.location.aggs.%s", tooltip.Alias), gjson.Parse(fmt.Sprintf(`{"sum":{"field":"%s"}}`, tooltip.Field)).Value())
	}

	return
}
