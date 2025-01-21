package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func SwipeRoutes(route *gin.Engine) {
	swipe := route.Group("/swipes")
	swipe.POST("/right", right)
	swipe.POST("/left", left)
}

func right(c *gin.Context) {
	var swipe models.RSwipe
	c.ShouldBindJSON(&swipe)
	// Todo: Get user from token

	var existingUser models.User
	config.DB.First(&existingUser, swipe.UserID)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{
			"error": "User not found",
		})
	}

	if existingUser.SwipeNum <= 0 {
		c.JSON(400, gin.H{
			"error": "Out of swipes",
		})
	}

	config.DB.Create(&swipe)
	config.DB.Model(&existingUser).Update("swipe_num", existingUser.SwipeNum-1)
}

func left(c *gin.Context) {
	var swipe models.LSwipe
	c.ShouldBindJSON(&swipe)

	// Todo: Get user from token

	var existingUser models.User
	config.DB.First(&existingUser, swipe.UserID)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{
			"error": "User not found",
		})
	}

	if existingUser.SwipeNum <= 0 {
		c.JSON(400, gin.H{
			"error": "Out of swipes",
		})

	}

	swipe.DeleteOn = time.Now().AddDate(0, 0, 1)

	config.DB.Create(&swipe)
	config.DB.Model(&existingUser).Update("swipe_num", existingUser.SwipeNum-1)
}
