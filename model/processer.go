package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

//Address of OCR server and options
const (
	OCRServer = "http://tesseract_server:8884/tesseract"
	OCROptions = "{\"languages\":[\"eng\"]}"
)

// ImageProcess Get ORC text and save it in redis database
func ImageProcess(rawData []byte, taskID string) {

	//Set redis client

	//Convert b64 data to image and send it to tesseract server
	//Get OCR text
	GetImage(rawData, taskID + ".png")
	oriText := GetText(taskID + ".png")

	//Send result to redis database
	err := SendToRedis(taskID, oriText)
	if err != nil {
		panic(err)
	}

	//Delete generated image
	os.Remove(taskID + ".png")
}

// GetImage Generate image based on base64 data
func GetImage(bytedata []byte, imageName string) {

	//Convert b64 data to image
	var m map[string]interface{}
	_ = json.Unmarshal(bytedata, &m)
	b64code := fmt.Sprintf("%v", m["image_data"])
	dec, err := base64.StdEncoding.DecodeString(b64code)
	if err != nil {
		panic(err)
	}

	//save image data in current directory
	f, err := os.Create(imageName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil{
		panic(err)
	}
}

// GetID Get task ID from GET request
func GetID(bytedata []byte) (res string){

	var m map[string]interface{}
	_ = json.Unmarshal(bytedata, &m)
	taskID := fmt.Sprintf("%v", m["task_id"])

	return taskID
}

func GetText(filename string) (res string){

	//Set OCR server address
	remoteURL := OCRServer

	//prepare the reader instances based on OCR server API
	values := map[string]io.Reader{
		"file": mustOpen(filename),
		"options": strings.NewReader(OCROptions),
	}

	//Send image and OCR options to OCR server
	//and Get response
	res, err := Upload(remoteURL, values)
	if err != nil {
		panic(err)
	}
	return res
}


func Upload( url string, values map[string]io.Reader) (res string, err error) {

	// Prepare data that will be submitted to OCR server.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return
		}

	}
	w.Close()

	// Submit it to server.
	resp, err := http.Post(url, w.FormDataContentType() ,&b)
	if err != nil {
		return
	}

	//Get response from OCR server and return OCR text only
	var data map[string]interface{}
	err1 := json.NewDecoder(resp.Body).Decode(&data)
	if err1 != nil {
		return
	}
	tmp := data["data"]
	output := tmp.(map[string]interface{})
	resText := output["stdout"].(string)
	return resText, err

}

//Error test for open Image file
func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
