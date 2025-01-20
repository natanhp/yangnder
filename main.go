package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natanhp/yangnder/config"
)

func init() {
	config.Connect()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "peng",
		})
	})
	r.Run()
}
