package controller

import (
	"context"
	"log"
	"net/http"
	"project_shopping_tour/service_content/model"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func ControllerReadContent(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "product"
	const BUCKET = "padtravel"
	var arrayOfReplyContent []model.ModelContentReplay
	var arrayOfContent []model.ModelContent
	ctx := context.Background()
	clientDatastore, err := datastore.NewClient(ctx, PROJECTID)
	if err != nil {
		log.Println("err cannot create client datastore => ", err)
	}

	keys, err := clientDatastore.GetAll(ctx, datastore.NewQuery(KIND), &arrayOfContent)
	if err != nil {
		log.Println("err in ControllerGetAllPath => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": err})
	}

	for _, key := range keys {
		log.Println(key)
	}

	c.JSON(http.StatusOK, gin.H{"reply": arrayOfReplyContent})
}
