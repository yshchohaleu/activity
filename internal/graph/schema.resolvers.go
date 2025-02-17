package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.62

import (
	"activity/internal/auth"
	"activity/internal/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Date is the resolver for the date field.
func (r *activityResolver) Date(ctx context.Context, obj *models.Activity) (string, error) {
	return obj.Date.Format("2006-01-02"), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *activityResolver) CreatedAt(ctx context.Context, obj *models.Activity) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// CreateActivity is the resolver for the createActivity field.
func (r *mutationResolver) CreateActivity(ctx context.Context, input CreateActivityInput) (*models.Activity, error) {
	userIDValue := ctx.Value(auth.UserIDKey)
	if userIDValue == nil {
		return nil, errors.New("user not authenticated")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("invalid user context")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, err
	}

	// Verify tag exists and user has access to it
	var tag models.Tag
	if err := r.DB.First(&tag, "id = ?", input.TagID).Error; err != nil {
		return nil, errors.New("tag not found")
	}

	activity := &models.Activity{
		ID:        uuid.New().String(),
		UserID:    userID,
		TagID:     input.TagID,
		Date:      date,
		Duration:  input.Duration,
		CreatedAt: time.Now(),
	}

	if err := r.DB.Create(activity).Error; err != nil {
		return nil, err
	}

	return activity, nil
}

// CreateTag is the resolver for the createTag field.
func (r *mutationResolver) CreateTag(ctx context.Context, value string) (*models.Tag, error) {
	userIDValue := ctx.Value(auth.UserIDKey)
	if userIDValue == nil {
		return nil, errors.New("user not authenticated")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("invalid user context")
	}

	tag := &models.Tag{
		ID:        uuid.New().String(),
		Value:     value,
		CreatorID: userID,
		CreatedAt: time.Now(),
	}

	if err := r.DB.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	userIDValue := ctx.Value(auth.UserIDKey)
	if userIDValue == nil {
		return nil, errors.New("user not authenticated")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("invalid user context")
	}

	var user models.User
	if err := r.DB.Preload("Activities").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Activities is the resolver for the activities field.
func (r *queryResolver) Activities(ctx context.Context) ([]*models.Activity, error) {
	userIDValue := ctx.Value(auth.UserIDKey)
	if userIDValue == nil {
		return nil, errors.New("user not authenticated")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("invalid user context")
	}

	var activities []*models.Activity
	if err := r.DB.Where("user_id = ?", userID).Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context) ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := r.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *tagResolver) CreatedAt(ctx context.Context, obj *models.Tag) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// Activity returns ActivityResolver implementation.
func (r *Resolver) Activity() ActivityResolver { return &activityResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Tag returns TagResolver implementation.
func (r *Resolver) Tag() TagResolver { return &tagResolver{r} }

type activityResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
