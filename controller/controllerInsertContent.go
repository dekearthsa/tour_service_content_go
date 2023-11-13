package controller

import (
	"context"
	"io"
	"log"
	"net/http"
	"project_shopping_tour/service_content/model"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func ControllerInsertContent(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "content"
	const BUCKET = "padtravel"

	var timing = time.Now().UnixNano()
	var arrayImagePath []string

	ctx := context.Background()
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("err MultipartForm => ", err)
	}

	files := form.File["images"]
	title := form.Value["title"]
	content := form.Value["content"]

	clientDatastore, err := datastore.NewClient(ctx, PROJECTID)
	if err != nil {
		log.Println("err in create client datastore => ", err)
	}

	key := datastore.IncompleteKey(KIND, nil)

	for _, file := range files {
		size := file.Size
		if size >= 5000000 {
			log.Println("error file to big.")
			c.JSON(http.StatusRequestHeaderFieldsTooLarge, gin.H{"Status": "file must less than 5MB."})
		}

		src, err := file.Open()
		if err != nil {
			log.Println("err open file => ", err)
		}
		defer src.Close()

		imagePath := title[0] + "_" + file.Filename + strconv.Itoa(int(timing))
		arrayImagePath = append(arrayImagePath, imagePath)
		clientStorage, err := storage.NewClient(ctx)
		if err != nil {
			log.Println("err in create client cloud storage => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Status": "internal error."})
		}
		bucket := clientStorage.Bucket(BUCKET)
		wc := bucket.Object(imagePath).NewWriter(ctx)
		_, err = io.Copy(wc, src)
		if err != nil {
			log.Println("err when write in cloud storage => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Status": "can't writer object"})
		}
		err = wc.Close()
		if err != nil {
			log.Println("err when close cloud storage => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "cannot close cloud storage"})
		}
	}

	payload := model.ModelContent{
		Title:     title[0],
		Content:   content[0],
		ImagePath: arrayImagePath,
	}

	if _, err := clientDatastore.Put(ctx, key, &payload); err != nil {
		log.Print("cannot write in datastore => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "cannot write in datastore"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
