package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Dipesh1203/alive/apps/backend/db"
)

// CreateUserProfile creates a user profile
func CreateUserProfile(ctx context.Context, database *db.PrismaClient, userID, firstName, lastName, phone, bio, avatar string) (*db.UserProfileModel, error) {
	// Check if profile already exists
	existing, _ := database.UserProfile.FindUnique(db.UserProfile.UserID.Equals(userID)).Exec(ctx)
	if existing != nil {
		return nil, errors.New("user profile already exists")
	}

	profile, err := database.UserProfile.CreateOne(
		db.UserProfile.User.Link(db.User.ID.Equals(userID)),
		db.UserProfile.FirstName.Set(firstName),
		db.UserProfile.LastName.Set(lastName),
		db.UserProfile.Phone.Set(phone),
		db.UserProfile.Bio.Set(bio),
		db.UserProfile.Avatar.Set(avatar),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to create user profile")
	}

	return profile, nil
}

// GetUserProfile retrieves a user's profile
func GetUserProfile(ctx context.Context, database *db.PrismaClient, userID string) (*db.UserProfileModel, error) {
	profile, err := database.UserProfile.FindUnique(db.UserProfile.UserID.Equals(userID)).Exec(ctx)
	if err != nil {
		return nil, errors.New("user profile not found")
	}

	return profile, nil
}

// UpdateUserProfile updates only the fields provided for a user's profile.
func UpdateUserProfile(ctx context.Context, database *db.PrismaClient, userID string, firstName, lastName, phone, bio, avatar *string) (*db.UserProfileModel, error) {
	params := make([]db.UserProfileSetParam, 0)
	if firstName != nil {
		params = append(params, db.UserProfile.FirstName.Set(*firstName))
	}
	if lastName != nil {
		params = append(params, db.UserProfile.LastName.Set(*lastName))
	}
	if phone != nil {
		params = append(params, db.UserProfile.Phone.Set(*phone))
	}
	if bio != nil {
		params = append(params, db.UserProfile.Bio.Set(*bio))
	}
	if avatar != nil {
		params = append(params, db.UserProfile.Avatar.Set(*avatar))
	}
	if len(params) == 0 {
		return nil, errors.New("no fields to update")
	}
	log.Printf("Updating user profile for userID %s with params: %+v", userID, params)

	profile, err := database.UserProfile.FindUnique(db.UserProfile.UserID.Equals(userID)).Update(params...).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("DB ERROR: %v", err)
	}

	return profile, nil
}

// UpdateUserPreferences updates user preferences (settings)
func UpdateUserPreferences(ctx context.Context, database *db.PrismaClient, userID string, preferences map[string]interface{}) (*db.UserProfileModel, error) {
	preferencesJSON, err := json.Marshal(preferences)
	if err != nil {
		return nil, errors.New("failed to marshal preferences")
	}

	profile, err := database.UserProfile.FindUnique(db.UserProfile.UserID.Equals(userID)).Update(
		db.UserProfile.Preferences.Set(preferencesJSON),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to update preferences")
	}

	return profile, nil
}

// GetUserPreferences retrieves user preferences
func GetUserPreferences(ctx context.Context, database *db.PrismaClient, userID string) (map[string]interface{}, error) {
	profile, err := GetUserProfile(ctx, database, userID)
	if err != nil {
		return nil, err
	}

	var preferences map[string]interface{}
	preferencesValue, ok := profile.Preferences()
	if !ok || len(preferencesValue) == 0 {
		return make(map[string]interface{}), nil
	}
	err = json.Unmarshal(preferencesValue, &preferences)
	if err != nil {
		return nil, errors.New("failed to unmarshal preferences")
	}

	return preferences, nil
}

// DeleteUserProfile deletes a user's profile
func DeleteUserProfile(ctx context.Context, database *db.PrismaClient, userID string) (*db.UserProfileModel, error) {
	profile, err := database.UserProfile.FindUnique(db.UserProfile.UserID.Equals(userID)).Delete().Exec(ctx)
	if err != nil {
		return nil, errors.New("failed to delete user profile")
	}

	return profile, nil
}
