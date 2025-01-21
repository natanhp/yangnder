package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natanhp/yangnder/auth"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func PremiumRoutes(route *gin.Engine) {
	premium := route.Group("/premiums", auth.AuthenticateMiddleware)
	premium.PATCH("/buy", BuyPremium)
}

func BuyPremium(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))
	var user models.User

	config.DB.First(&user, userID)

	if user.IsVerified {
		c.JSON(400, gin.H{
			"error": "User already premium",
		})

		return
	}

	config.DB.Model(&user).Update("is_verified", true)

	c.JSON(200, gin.H{
		"data": user,
	})
}
