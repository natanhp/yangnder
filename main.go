package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/controllers"
)

func init() {
	config.Connect()
}

func main() {
	r := gin.Default()
	controllers.UserRoutes(r)
	controllers.SwipeRoutes(r)
	r.Run()
}
