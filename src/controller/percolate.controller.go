package controller

import (
	"gin-elastic-percolator/src/model"
	"gin-elastic-percolator/src/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type PercolateController struct {
	router  *gin.RouterGroup
	service *service.PercolateService
}

func NewPercolateController(router *gin.RouterGroup) (o *PercolateController) {
	o = &PercolateController{
		router:  router,
		service: service.NewPercolateService(),
	}

	layerMaps := router.Group("percolate")
	layerMaps.POST("/add-query", o.AddQuery)
	layerMaps.POST("/search", o.GetPercolate_Data)

	return
}

// @Tags Percolator
// @Accept json
// @Param parameter body model.Percolate true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /percolate/add-query [post]
func (o *PercolateController) AddQuery(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Percolate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		log.Println(err)
		return
	}
	resp = o.service.AddQuery(param)
}

// @Tags Percolator
// @Accept multipart/form-data
// @Param file formData file true "this is json or csv file"
// @Produce json
// @Success 200 {object} object{data=model.Percolate_Data,meta_data=model.MetadataResponse} "OK"
// @Router /percolate/search [post]
func (o *PercolateController) GetPercolate_Data(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	file, _ := ctx.FormFile("file")
	log.Println(file.Filename)

	dst := "public/" + file.Filename
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		log.Println(err)
		return
	}

	resp.Data, resp.Metadata.Message = o.service.GetPercolate_Data(file.Filename)
}
