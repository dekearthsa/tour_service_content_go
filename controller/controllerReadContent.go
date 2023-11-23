package controller

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"project_shopping_tour/service_content/model"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func ControllerReadContent(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "content"
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

	for idx, el := range arrayOfContent {
		var arrayBase64 []string
		for _, item := range el.ImagePath {
			client, err := storage.NewClient(ctx)
			buckets := client.Bucket(BUCKET)
			rc, err := buckets.Object(item).NewReader(ctx)
			if err != nil {
				log.Println("err when fetch out storage ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "err when fetch out storage"})
			}
			byteFile, err := io.ReadAll(rc)
			defer rc.Close()
			if err != nil {
				log.Println("err read file from bucket", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "err read file from bucket"})
			}
			str := base64.StdEncoding.EncodeToString(byteFile)
			arrayBase64 = append(arrayBase64, str)
		}
		dataOut := model.ModelContentReplay{
			Title:     arrayOfContent[idx].Title,
			Content:   arrayOfContent[idx].Content,
			ImgBase64: arrayBase64,
		}
		arrayOfReplyContent = append(arrayOfReplyContent, dataOut)
	}

	c.JSON(http.StatusOK, gin.H{"reply": arrayOfReplyContent})
}
