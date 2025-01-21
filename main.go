package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/controllers"
	"github.com/natanhp/yangnder/scheduler"
)

func init() {
	config.Connect()
}

func main() {
	r := gin.Default()
	controllers.UserRoutes(r)
	controllers.SwipeRoutes(r)

	go func() {
		for {
			scheduler.Start()
		}
	}()

	r.Run()
}
