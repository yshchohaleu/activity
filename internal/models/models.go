package models

import (
	"time"
)

type User struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Email       string     `json:"email"`
	CreatedAt   time.Time  `json:"createdAt"`
	Activities  []Activity `json:"activities"`
	CreatedTags []Tag      `json:"createdTags" gorm:"foreignKey:CreatorID"`
}

type Tag struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Value     string    `json:"value"`
	CreatorID string    `json:"creatorId" gorm:"index"`
	CreatedAt time.Time `json:"createdAt"`
	Creator   User      `json:"creator" gorm:"foreignKey:CreatorID"`
}

type Activity struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userId" gorm:"index"`
	TagID     string    `json:"tagId" gorm:"index"`
	Date      time.Time `json:"date"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Tag       Tag       `json:"tag" gorm:"foreignKey:TagID"`
}
