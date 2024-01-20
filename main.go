package main

import (
	"speech/api"
	"speech/service"

	"github.com/gin-gonic/gin"
)

func main() {
	go service.Speek()

	r := gin.Default()
	r.POST("/", api.Speech)
	r.Run() // listen and serve on 0.0.0.0:8080
}
