package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natanhp/yangnder/auth"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func UserRoutes(route *gin.Engine) {
	user := route.Group("/users")
	user.GET("", auth.AuthenticateMiddleware, FindAll)
	user.POST("/register", Create)
	user.GET("/detail/:id", auth.AuthenticateMiddleware, FindOne)
	user.POST("/upload-photo", auth.AuthenticateMiddleware, UploadPhoto)
	user.POST("/login", Login)
}

func FindAll(c *gin.Context) {
	var users []models.User
	claims := c.MustGet("claims").(jwt.MapClaims)
	id := uint(claims["sub"].(float64))

	query := `
		SELECT 
		id, 
		email, 
		name, 
		dob, 
		desc, 
		photo 
		FROM 
		users 
		WHERE 
		id != ? 
		AND id NOT IN (
			SELECT 
			r_swipe_id 
			FROM 
			r_swipes 
			WHERE 
			user_id = ? 
			UNION 
			SELECT 
			l_swipe_id 
			FROM 
			l_swipes 
			WHERE 
			user_id = ?
		)
	`
	err := config.DB.Raw(query, id, id, id).Scan(&users).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to fetch users",
		})
		return

	}

	c.JSON(200, gin.H{
		"data": users,
	})
}

func FindOne(c *gin.Context) {
	var user models.User
	config.DB.First(&user, c.Param("id"))

	c.JSON(200, gin.H{
		"data": user,
	})
}

func Create(c *gin.Context) {
	var user models.User
	c.ShouldBindJSON(&user)

	var existingUser models.User
	config.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		c.JSON(400, gin.H{
			"error": "Email already taken",
		})
		return
	}

	argonParams := &argon2id.Params{
		Memory:      19 * 1024,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   16,
	}

	hash, err := argon2id.CreateHash(user.Password, argonParams)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create user",
		})
	}

	user.Password = hash
	user.SwipeNum = 10
	user.NextSwipeReset = time.Now().AddDate(0, 0, 1)

	config.DB.Create(&user)

	user.Password = ""

	c.JSON(201, gin.H{
		"data": user,
	})
}

func UploadPhoto(c *gin.Context) {
	var user models.User
	claims := c.MustGet("claims").(jwt.MapClaims)
	id := uint(claims["sub"].(float64))

	config.DB.First(&user, id)

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Photo is required",
		})
		return
	}

	fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	os.Remove("photos/" + user.Photo)
	err = c.SaveUploadedFile(file, "photos/"+fileName)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save photo",
		})
		return
	}

	user.Photo = fileName
	config.DB.Save(&user)

	c.JSON(200, gin.H{
		"data": user,
	})
}

func Login(c *gin.Context) {
	var user models.User
	c.ShouldBindJSON(&user)

	var existingUser models.User
	config.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	match, err := argon2id.ComparePasswordAndHash(user.Password, existingUser.Password)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to login",
		})
		return
	}

	if !match {
		c.JSON(400, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	existingUser.Password = ""

	token, err := auth.CreateToken(existingUser.ID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to login",
		})
		return
	}

	c.JSON(200, gin.H{
		"data":  existingUser,
		"token": token,
	})
}
