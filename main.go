package main

import (
	"github.com/DenrianWeiss/catballCoin/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/task/count", handler.GetTaskCount)

	r.POST("/task/add/:key", handler.AddTask)

	_ = r.Run()
}
