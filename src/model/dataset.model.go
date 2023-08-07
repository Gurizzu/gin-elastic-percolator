package model

type Data struct {
	MetadataWithID `bson:",inline"`

	Metadata    *Data_Metadata `json:"metadata,omitempty" bson:"metadata"`
	Schemas     *[]Schema      `json:"schemas,omitempty" bson:"schemas"`
	Config      *Config        `json:"config,omitempty" bson:"config"`
	TableName   string         `json:"tableName,omitempty" bson:"tableName,omitempty" swaggerignore:"true"`
	FileName    string         `json:"fileName" bson:"fileName"`
	TotalViewed int64          `json:"totalViewed" bson:"totalViewed,omitempty" swaggerignore:"true"`
}

type Data_Metadata struct {
	DataName           string   `json:"dataName" bson:"dataName"`
	Indicator          string   `json:"indicator" bson:"indicator"`
	RegionID           string   `json:"regionId" bson:"regionId"`
	CategoriesID       []string `json:"categoriesId" bson:"categoriesId"`
	DataSource         []string `json:"dataSource" bson:"dataSource"`
	ExternalDataSource []string `json:"externalDataSource" bson:"externalDataSource"`
	IndicatorCode      string   `json:"indicatorCode" bson:"indicatorCode"`
	DataUnit           []string `json:"dataUnit" bson:"dataUnit"`
	DataPeriod         string   `json:"dataPeriod" bson:"dataPeriod"`
	StartPeriod        int64    `json:"startPeriod" bson:"startPeriod"`
	EndPeriod          int64    `json:"endPeriod" bson:"endPeriod"`
	Media              string   `json:"media" bson:"media"`
	Description        string   `json:"description" bson:"description"`
}

type Config struct {
	PublicStatus      bool              `json:"publicStatus" bson:"publicStatus"`
	DataVisualivation DataVisualivation `json:"dataVisualivation" bson:"dataVisualivation"`
	CategoryData      string            `json:"categoryData" bson:"categoryData"`
	Downloadable      bool              `json:"downloadable" bson:"downloadable"`
}

type DataVisualivation struct {
	Table bool `json:"table" bson:"table"`
	Graph bool `json:"graph" bson:"graph"`
	Maps  bool `json:"maps" bson:"maps"`
}

type Schema struct {
	Name          string `json:"name" bson:"name"`
	DataType      string `json:"dataType" bson:"dataType"`
	Alias         string `json:"alias" bson:"alias"`
	DataTypeAlias string `json:"dataTypeAlias" bson:"dataTypeAlias"`
	NameKeyword   string `json:"nameKeyword,omitempty" bson:"nameKeyword,omitempty" swaggerIgnore:"true"`
}

type Data_Table struct {
	RowData []map[string]interface{} `json:"rowData"`
	Schemas []Schema                 `json:"schemas"`
}

type Data_View struct {
	Data `bson:",inline"`
}
