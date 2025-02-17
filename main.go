package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"my-go-user-center/docs"
	"my-go-user-center/src/config"
	"my-go-user-center/src/job"
	"my-go-user-center/src/router"
)

//@contact.name   API Support
//@contact.url    http://www.swagger.io/support
//@contact.email  support@swagger.io

func main() {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	log.Info("The project is start!")
	config.InitConfig()
	go job.InitJob()
	config.InitMySQL()
	config.InitRedis()
	r := router.InitRouterAndServe() // router.Router()
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(viper.GetString("port.server"))
}
