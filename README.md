### Introduction
This is a gin based backend that implement 
an OCR system with tesseract.
##### Framework
I use golang as my programming language 
and gin as my web Framework because golang has
good support for concurrent work.
##### OCR
I use a lightweight HTTP [tesseract-server](https://github.com/hertzg/tesseract-server)
from github.
##### Database
I use redis as my database for better scalability 
and performance with [go-redis](https://github.com/go-redis/redis).
##### Usability
I use [docker](https://www.docker.com) to build my application 
and use docker-compose to integrate database, 
OCR server and my application for the usability 
of my application.
### Tested on:  
OS: MacOS Big Sur 11.2.3 (Apple M1 chip)  
docker Engine version: 20.10.6  
docker-compose version: 1.29.1  

### Get Start:
cd to the directory of this project and run

```
$ docker-compose build
$ docker-compose up
```
Get the recognition result:
```
$ curl -XPOST "http://localhost:5000/image-sync" -d '{"image_data": "<b64 encoded image>"}'
```
Get a picture ID (for later query):
```
curl -XPOST "http://localhost:5000/image" -d '{"image_data": "<b64 encoded image>"}'
```
Get the recognition result of with:
```
curl -XGET "http://localhost:5000/image" \
    -d '{"task_id": "<picture id as received from POST /image>"}'
```


Note: The image_data in the request should be pure base64 encoded image


