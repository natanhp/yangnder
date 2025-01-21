package controllers_test

import (
	"fmt"
	"log"
	"time"

	"github.com/natanhp/yangnder/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Sqlite driver based on CGO

// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

var DBTest *gorm.DB

func ConnectTest() {
	var err error
	DBTest, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}
}

func MigrateTest() {
	var err error

	err = DBTest.SetupJoinTable(&models.User{}, "RSwipes", &models.RSwipe{})

	if err != nil {
		log.Fatal("Failed to create join table")
	}

	err = DBTest.SetupJoinTable(&models.User{}, "LSwipes", &models.LSwipe{})

	if err != nil {
		log.Fatal("Failed to create join table")
	}

	err = DBTest.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Failed to migrate user")
	}
}

func PopulateUsersTest() {
	for i := 0; i <= 10; i++ {
		log.Println("Creating user", i)
		user := models.User{
			Name:           fmt.Sprintf("User %d", i),
			Email:          fmt.Sprintf("user%d", i) + "@example.com",
			Password:       "password",
			DOB:            "1990-01-01",
			DESC:           "I am a user",
			SwipeNum:       10,
			NextSwipeReset: time.Now().AddDate(0, 0, 1),
		}

		DBTest.Create(&user)
	}
}
