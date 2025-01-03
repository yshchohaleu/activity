package migrations

import (
	"log"
	"time"

	"gorm.io/gorm"

	"activity/internal/models"
)

// EnsureDefaultUser creates a default user if it doesn't exist
func EnsureDefaultUser(db *gorm.DB) {
	var user models.User
	result := db.Where("email = ?", "y.shchohaleu@gmail.com").First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		newUser := models.User{
			ID:        "iKezbgQvzHc27zyl9AYCUT1C89g2",
			Email:     "y.shchohaleu@gmail.com",
			CreatedAt: time.Now(),
		}
		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		} else {
			log.Printf("Created user with email: %s", newUser.Email)
		}
	} else if result.Error != nil {
		log.Printf("Error checking for existing user: %v", result.Error)
	} else {
		log.Printf("User with email %s already exists", user.Email)
	}
}
