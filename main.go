package main

import (
	"log"
	"project_shopping_tour/service_content/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const port string = ":7701"

func main() {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE", "GET"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin,access-control-allow-headers"},
	}))

	router.GET("/api/debug", controller.ControllerDebug)
	router.POST("/api/insert/content", controller.ControllerInsertContent)
	router.GET("/api/read/content", controller.ControllerReadContent)
	router.POST("/api/update/content", controller.ControllerUpdateContent)
	router.POST("/api/delete/content", controller.ControllerDeleteContent)

	err := router.Run(port)
	if err != nil {
		log.Println("Service admin CRUD fail to start" + err.Error())
	}
	log.Println("Service admin CRUD start at port" + port + "debug => " + "http://localhost:" + port + "/debug")
}
