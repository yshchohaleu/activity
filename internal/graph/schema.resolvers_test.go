package graph

import (
	"activity/internal/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityResolver_Date(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	testDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	activity := &models.Activity{
		Date: testDate,
	}

	r := resolver.Activity()
	result, err := r.Date(context.TODO(), activity)

	require.NoError(t, err)
	assert.Equal(t, "2025-01-01", result)
}

func TestActivityResolver_CreatedAt(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	activity := &models.Activity{
		CreatedAt: testTime,
	}

	r := resolver.Activity()
	result, err := r.CreatedAt(context.TODO(), activity)

	require.NoError(t, err)
	assert.Equal(t, "2025-01-01T12:00:00Z", result)
}

func TestMutationResolver_CreateActivity(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test user and tag
	user := createTestUser(t, resolver.DB)
	tag := createTestTag(t, resolver.DB, user.ID)

	tests := []struct {
		name    string
		ctx     func() context.Context
		input   CreateActivityInput
		wantErr string
	}{
		{
			name: "creates valid activity",
			ctx:  func() context.Context { return contextWithUserID(user.ID) },
			input: CreateActivityInput{
				TagID:    tag.ID,
				Date:     "2025-01-01",
				Duration: 30,
			},
			wantErr: "",
		},
		{
			name: "fails with invalid date format",
			ctx:  func() context.Context { return contextWithUserID(user.ID) },
			input: CreateActivityInput{
				TagID:    tag.ID,
				Date:     "invalid-date",
				Duration: 30,
			},
			wantErr: "parsing time",
		},
		{
			name: "fails with non-existent tag",
			ctx:  func() context.Context { return contextWithUserID(user.ID) },
			input: CreateActivityInput{
				TagID:    "non-existent-tag",
				Date:     "2025-01-01",
				Duration: 30,
			},
			wantErr: "tag not found",
		},
		{
			name: "fails without user context",
			ctx:  func() context.Context { return context.Background() },
			input: CreateActivityInput{
				TagID:    tag.ID,
				Date:     "2025-01-01",
				Duration: 30,
			},
			wantErr: "user not authenticated",
		},
	}

	r := resolver.Mutation()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.CreateActivity(tt.ctx(), tt.input)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.input.TagID, result.TagID)
				assert.Equal(t, tt.input.Duration, result.Duration)
				assert.Equal(t, user.ID, result.UserID)
			}
		})
	}
}

func TestMutationResolver_CreateTag(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test user
	user := createTestUser(t, resolver.DB)

	tests := []struct {
		name    string
		ctx     func() context.Context
		value   string
		wantErr string
	}{
		{
			name:    "creates valid tag",
			ctx:     func() context.Context { return contextWithUserID(user.ID) },
			value:   "RUNNING",
			wantErr: "",
		},
		{
			name:    "fails without user context",
			ctx:     func() context.Context { return context.Background() },
			value:   "SWIMMING",
			wantErr: "user not authenticated",
		},
	}

	r := resolver.Mutation()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.CreateTag(tt.ctx(), tt.value)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.value, result.Value)
				assert.Equal(t, user.ID, result.CreatorID)
				assert.NotZero(t, result.CreatedAt)
			}
		})
	}
}

func TestQueryResolver_Me(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test user with activities
	user := createTestUser(t, resolver.DB)
	_ = createTestActivity(t, resolver.DB, user.ID)

	tests := []struct {
		name    string
		ctx     func() context.Context
		wantErr string
	}{
		{
			name:    "returns user with activities",
			ctx:     func() context.Context { return contextWithUserID(user.ID) },
			wantErr: "",
		},
		{
			name:    "fails without user context",
			ctx:     func() context.Context { return context.Background() },
			wantErr: "user not authenticated",
		},
		{
			name:    "fails with non-existent user",
			ctx:     func() context.Context { return contextWithUserID("non-existent") },
			wantErr: "record not found",
		},
	}

	r := resolver.Query()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.Me(tt.ctx())
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, user.ID, result.ID)
				assert.Equal(t, user.Email, result.Email)
				assert.NotEmpty(t, result.Activities)
			}
		})
	}
}

func TestQueryResolver_Activities(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test user with activities
	user := createTestUser(t, resolver.DB)
	activity := createTestActivity(t, resolver.DB, user.ID)

	tests := []struct {
		name    string
		ctx     func() context.Context
		wantErr string
	}{
		{
			name:    "returns user activities",
			ctx:     func() context.Context { return contextWithUserID(user.ID) },
			wantErr: "",
		},
		{
			name:    "fails without user context",
			ctx:     func() context.Context { return context.Background() },
			wantErr: "user not authenticated",
		},
	}

	r := resolver.Query()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.Activities(tt.ctx())
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, activity.ID, result[0].ID)
				assert.Equal(t, activity.TagID, result[0].TagID)
				assert.Equal(t, activity.Duration, result[0].Duration)
			}
		})
	}
}

func TestQueryResolver_Tags(t *testing.T) {
	resolver, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test user and tags
	user := createTestUser(t, resolver.DB)
	tag1 := createTestTag(t, resolver.DB, user.ID)

	// Create another tag with different value
	tag2 := &models.Tag{
		ID:        "test-tag-id-2",
		Value:     "SWIMMING",
		CreatorID: user.ID,
		CreatedAt: time.Now(),
	}
	require.NoError(t, resolver.DB.Create(tag2).Error)

	r := resolver.Query()
	result, err := r.Tags(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.ElementsMatch(t, []string{tag1.Value, tag2.Value}, []string{result[0].Value, result[1].Value})
}
