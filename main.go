package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natanhp/yangnder/config"
	users "github.com/natanhp/yangnder/controllers"
)

func init() {
	config.Connect()
}

func main() {
	r := gin.Default()
	users.Routes(r)
	r.Run()
}
