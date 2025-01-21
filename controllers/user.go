package controllers

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/natanhp/yangnder/auth"
	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func UserRoutes(route *gin.Engine) {
	user := route.Group("/users")
	user.GET("", findAll)
	user.POST("/register", create)
	user.GET("/detail", findOne)
	user.POST("/upload-photo", uploadPhoto)
	user.POST("/login", login)
}

func findAll(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)
	c.JSON(200, gin.H{
		"data": users,
	})
}

func findOne(c *gin.Context) {
	var user models.User
	config.DB.First(&user, c.Param("id"))

	c.JSON(200, gin.H{
		"data": user,
	})
}

func create(c *gin.Context) {
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

	c.JSON(200, gin.H{
		"data": user,
	})
}

func uploadPhoto(c *gin.Context) {
	var user models.User
	config.DB.First(&user, c.Param("id"))

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Photo is required",
		})
		return
	}

	err = c.SaveUploadedFile(file, "photos/"+file.Filename)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save photo",
		})
		return
	}

	user.Photo = file.Filename
	config.DB.Save(&user)

	c.JSON(200, gin.H{
		"data": user,
	})
}

func login(c *gin.Context) {
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
