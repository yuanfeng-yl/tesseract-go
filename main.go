package main

import (
	"github.com/gin-gonic/gin"
	"tesseract_go/controller"
)

// GoPort Address of OCR server, redis server and my app port
const GoPort = ":5000"

func main() {
	r := gin.Default()

	//async image process
	r.POST("/image", controller.AsyncImage)

	//Retrieve the OCR text with task ID
	r.GET("/image", controller.RetrieveText)

	//sync image process
	r.POST("/image-sync", controller.SyncImage)

	r.Run(GoPort)
}