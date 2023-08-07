package model

import (
	"time"

	"git.blackeye.id/Aldi.Rismawan/centrotil/db/umongo"
	"go.mongodb.org/mongo-driver/bson"
)

type MetadataResponse struct {
	Status        bool   `json:"status"`
	Message       string `json:"message"`
	TimeExecution string `json:"timeExecution"`

	Pagination *umongo.PaginationResponse `json:"pagination" bson:"-"`
}
type Response struct {
	Metadata MetadataResponse `json:"metaData"`
	Data     interface{}      `json:"data"`
}
type Response_Data_Upsert struct {
	ID string `json:"id"`
}

type Request struct {
	Request_Pagination
	Request_Search
}
type Range struct {
	Field string `json:"field" example:"updatedAt"`
	Start int64  `json:"start" example:"1646792565000"`
	End   int64  `json:"end" example:"1646792565000"`
}
type Request_Search struct {
	Range *Range `json:"range"`
}

func (o Request_Search) BaseHandle(filter bson.M, rangeField string) (res bson.M) {
	if requestRange := o.Range; requestRange != nil && requestRange.Field != "" {
		rangeField = o.Range.Field
	}
	if rangeField == "" {
		rangeField = "updatedAt"
	}
	res = filter

	if o.Range != nil {
		if o.Range.Start == 0 && o.Range.End == 0 {
			timeNow := time.Now()
			o.Range.End = timeNow.UnixMilli()
			o.Range.Start = timeNow.AddDate(0, 0, -7).UnixMilli()
		}
		filter[rangeField] = bson.M{
			"$gte": o.Range.Start, "$lt": o.Range.End,
		}
	}

	return
}

func (o Request_Search) Handle_RequestSearch(filter bson.M) (res bson.M) {
	return o.BaseHandle(filter, "")
}

type Metadata struct {
	CreatedAt int64  `json:"createdAt" bson:"createdAt,omitempty" swaggerignore:"true"`
	UpdatedAt int64  `json:"updatedAt" bson:"updatedAt,omitempty" swaggerignore:"true"`
	CreatedBy string `json:"createdBy" bson:"createdBy,omitempty"`
}
type MetadataWithID struct {
	IdDocument string `json:"id" bson:"_id"`
	Metadata   `bson:",inline"`
}
type Request_Pagination struct {
	OrderBy string `json:"orderBy" example:"createdAt"`
	Order   string `json:"order" example:"ASC"`

	Page int64 `example:"1" json:"page"`
	Size int64 `example:"11" json:"size"`
}
type PaginationResponse struct {
	Size          int   `json:"size"`
	TotalPages    int64 `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}
