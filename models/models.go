package models

import (
	"time"
)

type User struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Email          string    `json:"email" gorm:"unique"`
	Password       string    `json:"password"`
	Name           string    `json:"name"`
	DOB            string    `json:"dob" gorm:"type:date"`
	DESC           string    `json:"desc"`
	Photo          string    `json:"photo"`
	SwipeNum       int       `json:"swipe_num"`
	NextSwipeReset time.Time `json:"-"`
	RSwipes        []*User   `json:"r_swipes" gorm:"many2many:r_swipes"`
	LSwipes        []*User   `json:"l_swipes" gorm:"many2many:l_swipes"`
}

type RSwipe struct {
	UserID   uint `json:"user_id" gorm:"primaryKey"`
	RSwipeID uint `json:"r_swipe_id" gorm:"primaryKey"`
}

type LSwipe struct {
	UserID   uint      `json:"user_id" gorm:"primaryKey"`
	LSwipeID uint      `json:"l_swipe_id" gorm:"primaryKey"`
	DeleteOn time.Time `json:"-"`
}
