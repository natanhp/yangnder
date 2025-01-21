package auth

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(id uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
	})

	tokenString, err := claims.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthenticateMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}

	token = token[7:]

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}
