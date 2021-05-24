package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/teris-io/shortid"
	"net/http"
	"os"
	"tesseract_go/model"
)

// AsyncImage Response with task ID and async image process
func AsyncImage(c *gin.Context) {
	//Read data from POST request
	rawData, _ := c.GetRawData()

	//Generate unique task ID for this request
	taskID, err := shortid.Generate()
	if err != nil {
		panic(err)
	}

	//Concurrent image recognition implemented with goroutine
	go model.ImageProcess(rawData, taskID)

	//Response with generated task ID
	result := gin.H{
		"task_id": taskID,
	}
	c.JSON(http.StatusOK, result)

}

// RetrieveText Retrieve the OCR text with task ID
func RetrieveText(c *gin.Context) {

	//Get task ID from POST request
	rawData, _ := c.GetRawData()
	taskID := model.GetID(rawData)
	var result gin.H

	//Get OCR text from redis database with task ID
	val, err := model.GetFromRedis(taskID)
	if err == redis.Nil {
		result = gin.H{
			"task_id": nil,
		}
	} else if err != nil {
		panic(err)
	} else {
		result = gin.H{
			"task_id": val,
		}
	}

	//Response with task ID and OCR text
	c.JSON(http.StatusOK, result)
}

// SyncImage Sync image process and gives back text
func SyncImage(c *gin.Context) {
	//Read image data
	rawData, _ := c.GetRawData()

	//Convert b64 data to image and send it to tesseract server
	//Get OCR text
	model.GetImage(rawData, "ImageForSync.png")
	oriText := model.GetText("ImageForSync.png")
	err := os.Remove("ImageForSync.png")
	if err != nil {
		return 
	}

	//Response with OCR text
	result := gin.H{
		"text": oriText,
	}
	c.JSON(http.StatusOK, result)
}