package graph

import (
	"firebase.google.com/go/v4"
	"gorm.io/gorm"
)

type Resolver struct {
	DB  *gorm.DB
	App *firebase.App
}

func NewResolver(db *gorm.DB, app *firebase.App) *Resolver {
	return &Resolver{
		DB:  db,
		App: app,
	}
} 