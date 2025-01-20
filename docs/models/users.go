package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Email          string         `json:"email" gorm:"unique"`
	Password       string         `json:"password"`
	Name           string         `json:"name"`
	DOB            datatypes.Date `json:"dob"`
	DESC           string         `json:"desc"`
	Photo          string         `json:"photo"`
	SwipeNum       int            `json:"swipe_num"`
	NextSwipeReset time.Time      `json:"next_swipe_reset"`
	RSwipes        []*RSwipe      `json:"r_swipes" gorm:"many2many:r_swipes"`
	LSwipes        []*LSwipe      `json:"l_swipes" gorm:"many2many:l_swipes"`
}

type RSwipe struct {
	UserID   uint `json:"user_id" goem:"primaryKey"`
	RSwipeID uint `json:"r_swipe_id" goem:"primaryKey"`
}

type LSwipe struct {
	UserID   uint `json:"user_id" goem:"primaryKey"`
	LSwipeID uint `json:"l_swipe_id" goem:"primaryKey"`
	DeleteOn time.Time
}
