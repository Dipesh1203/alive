package services

import (
	"backend/db"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserSignup creates a new user with email and password
func UserSignup(ctx context.Context, database *db.PrismaClient, email string, password string) (*db.UserModel, error) {
	// Check if user already exists
	existingUser, err := database.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new user
	user, err := database.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Password.Set(string(hashedPassword)),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// UserLogin authenticates a user and returns the user if credentials are valid
func UserLogin(ctx context.Context, database *db.PrismaClient, email string, password string) (*db.UserModel, error) {
	// Find user by email
	user, err := database.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(ctx context.Context, database *db.PrismaClient, userID string) (*db.UserModel, error) {
	user, err := database.User.FindUnique(db.User.ID.Equals(userID)).Exec(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(ctx context.Context, database *db.PrismaClient, email string) (*db.UserModel, error) {
	user, err := database.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ListUsers retrieves all users
func ListUsers(ctx context.Context, database *db.PrismaClient) ([]db.UserModel, error) {
	users, err := database.User.FindMany().Exec(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	return users, nil
}

// DeleteUser removes a user from the database
func DeleteUser(ctx context.Context, database *db.PrismaClient, userID string) (*db.UserModel, error) {
	user, err := database.User.FindUnique(db.User.ID.Equals(userID)).Delete().Exec(ctx)
	if err != nil {
		return nil, errors.New("failed to delete user")
	}

	return user, nil
}
