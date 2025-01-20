package main

import (
	"log"

	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func init() {
	config.Connect()
}

func main() {
	var err error

	err = config.DB.SetupJoinTable(&models.User{}, "RSwipes", &models.RSwipe{})

	if err != nil {
		log.Fatal("Failed to create join table")
	}

	err = config.DB.SetupJoinTable(&models.User{}, "LSwipes", &models.LSwipe{})

	if err != nil {
		log.Fatal("Failed to create join table")
	}

	err = config.DB.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Failed to migrate user")
	}
}
