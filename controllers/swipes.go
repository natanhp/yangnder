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
	swipe.POST("/right", Right)
	swipe.POST("/left", Left)
}

func Right(c *gin.Context) {
	var swipe models.RSwipe
	c.ShouldBindJSON(&swipe)
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	var swiper models.User
	config.DB.First(&swiper, userID)

	if isSwiped(userID, swipe.RSwipeID, swiper, c) {
		return
	}

	swipe.UserID = userID

	config.DB.Create(&swipe)
	config.DB.Model(&swiper).Update("swipe_num", swiper.SwipeNum-1)

	c.JSON(201, gin.H{
		"data": swipe,
		"user": swiper,
	})
}

func isSwiped(userID uint, swipedID uint, swiper models.User, c *gin.Context) bool {
	var existingUser models.User
	config.DB.First(&existingUser, swipedID)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{
			"error": "User not found",
		})

		return true
	}

	if swiper.SwipeNum <= 0 {
		c.JSON(400, gin.H{
			"error": "Out of swipes",
		})

		return true
	}

	var existingRSwipe models.RSwipe
	config.DB.Where("user_id = ? AND r_swipe_id = ?", userID, swipedID).First(&existingRSwipe)

	if existingRSwipe.UserID != 0 {
		c.JSON(400, gin.H{
			"error": "Already swiped",
		})

		return true
	}

	var existingLSwipe models.LSwipe
	config.DB.Where("user_id = ? AND l_swipe_id = ?", userID, swipedID).First(&existingLSwipe)

	if existingLSwipe.UserID != 0 {
		c.JSON(400, gin.H{
			"error": "Already swiped",
		})

		return true
	}
	return false
}

func Left(c *gin.Context) {
	var swipe models.LSwipe
	c.ShouldBindJSON(&swipe)
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	var swiper models.User
	config.DB.First(&swiper, userID)

	if isSwiped(userID, swipe.LSwipeID, swiper, c) {
		return
	}

	swipe.UserID = userID
	swipe.DeleteOn = time.Now().AddDate(0, 0, 1)

	config.DB.Create(&swipe)
	config.DB.Model(&swiper).Update("swipe_num", swiper.SwipeNum-1)

	c.JSON(201, gin.H{
		"data": swipe,
		"user": swiper,
	})
}
