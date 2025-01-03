package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Activities []Activity
}

type Activity struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userId" gorm:"index"`
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
} 