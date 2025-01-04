package graph

import (
	"activity/internal/auth"
	"activity/internal/models"
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// contextWithUserID creates a new context with a user ID
func contextWithUserID(userID string) context.Context {
	return context.WithValue(context.Background(), auth.UserIDKey, userID)
}

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) (*Resolver, func()) {
	// Using an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate the test database
	err = db.AutoMigrate(&models.User{}, &models.Activity{}, &models.Tag{})
	require.NoError(t, err)

	// Create a new resolver with the test database
	// Note: We don't need Firebase for tests since we mock auth context
	resolver := &Resolver{
		DB: db,
	}

	// Return cleanup function
	cleanup := func() {
		// Get the underlying SQL database
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close()
	}

	return resolver, cleanup
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	user := &models.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}
	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

// createTestTag creates a test tag in the database
func createTestTag(t *testing.T, db *gorm.DB, userID string) *models.Tag {
	tag := &models.Tag{
		ID:        "test-tag-id",
		Value:     "RUNNING",
		CreatorID: userID,
		CreatedAt: time.Now(),
	}
	err := db.Create(tag).Error
	require.NoError(t, err)
	return tag
}

// createTestActivity creates a test activity in the database
func createTestActivity(t *testing.T, db *gorm.DB, userID string) *models.Activity {
	// Create a test tag first
	tag := createTestTag(t, db, userID)

	activity := &models.Activity{
		ID:        "test-activity-id",
		UserID:    userID,
		TagID:     tag.ID,
		Date:      time.Now(),
		Duration:  30,
		CreatedAt: time.Now(),
	}
	err := db.Create(activity).Error
	require.NoError(t, err)
	return activity
}
