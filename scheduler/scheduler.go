package scheduler

import (
	"time"

	"github.com/natanhp/yangnder/config"
	"github.com/natanhp/yangnder/models"
)

func resetSwipeNumber() {
	config.DB.Model(&models.User{}).Where("next_swipe_reset <=", time.Now()).Updates(map[string]interface{}{
		"swipe_num":        69,
		"next_swipe_reset": time.Now().AddDate(0, 0, 1),
	})
}

// func deleteLSwipes() {
// 	var lSwipes []models.LSwipe
// 	config.DB.Find(&lSwipes).Delete(&lSwipes)
// }

func Start() {
	resetSwipeNumber()
	// deleteLSwipes()
	time.Sleep(5 * time.Minute)
}
