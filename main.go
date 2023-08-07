package main

import (
	"fmt"
	"gin-elastic-percolator/docs"
	"gin-elastic-percolator/src/config"
	"gin-elastic-percolator/src/controller"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/logrusorgru/aurora"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	fmt.Println("Using timezone:", aurora.Green(time.Now().Location().String()))
}

func main() {
	log.Println(config.ENV_KEY_PORT)
	log.Println(os.Getenv(config.ENV_KEY_PORT))
	appPort := ":" + os.Getenv(config.ENV_KEY_PORT)

	rapidocPath := "docs/index.html"
	if os.Getenv("DEBUG_MODE") != "true" {
		gin.SetMode(gin.ReleaseMode)
		rapidocPath = "index.html"
	}
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(cors.Default())
	// router.Use(static.Serve("/", static.LocalFile("docs", true)))

	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "SUMUT LAYER MAPS"
	docs.SwaggerInfo.Description = "SUMUT LAYER MAPS DALAM PETA API DOCUMENTATION"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Printf("Rapidoc Path: %s\n", rapidocPath)
	router.StaticFile("/", rapidocPath)

	apiV1 := router.Group(basePath)
	apiV1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Nuwhofev"}) })

	controller.NewPercolateController(apiV1)

	log.Println(aurora.Green(
		fmt.Sprintf("http://localhost%s/swagger/index.html", appPort),
	))

	log.Fatalln(router.Run(appPort))
}
