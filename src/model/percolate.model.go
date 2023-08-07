package model

type Percolate struct {
	Query string   `json:"query"`
	Tags  []string `json:"tags"`
}

type Percolate_Data struct {
	Result interface{} `json:"data"`
}
