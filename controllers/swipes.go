package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natanhp/yangnder/auth"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func SwipeRoutes(route *gin.Engine) {
	swipe := route.Group("/swipes", auth.AuthenticateMiddleware)
	swipe.POST("/right", right)
	swipe.POST("/left", left)
}

func right(c *gin.Context) {
	var swipe models.RSwipe
	c.ShouldBindJSON(&swipe)
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	var existingUser models.User
	config.DB.First(&existingUser, uint(userID))

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{
			"error": "User not found",
		})

		return
	}

	if existingUser.SwipeNum <= 0 {
		c.JSON(400, gin.H{
			"error": "Out of swipes",
		})

		return
	}

	var existingSwipe models.RSwipe
	config.DB.Where("user_id = ? AND r_swipe_id = ?", userID, swipe.RSwipeID).First(&existingSwipe)

	if existingSwipe.UserID != 0 {
		c.JSON(400, gin.H{
			"error": "Already swiped",
		})

		return
	}

	swipe.UserID = userID

	config.DB.Create(&swipe)
	config.DB.Model(&existingUser).Update("swipe_num", existingUser.SwipeNum-1)

	c.JSON(201, gin.H{
		"data": swipe,
		"user": existingUser,
	})
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
