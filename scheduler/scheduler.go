package scheduler

import (
	"log"
	"time"

	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func resetSwipeNumber() {
	config.DB.Model(&models.User{}).Where("next_swipe_reset <= ?", time.Now()).Updates(map[string]interface{}{
		"swipe_num":        10,
		"next_swipe_reset": time.Now().AddDate(0, 0, 1),
	})
}

func deleteLSwipes() {
	config.DB.Where("delete_on <= ?", time.Now()).Delete(&models.LSwipe{})
}

func Start() {
	log.Default().Println("Scheduler started")
	resetSwipeNumber()
	deleteLSwipes()
	time.Sleep(5 * time.Minute)
}
