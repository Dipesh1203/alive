package services

import (
	"context"
	"errors"
	"log"

	db "github.com/Dipesh1203/alive/backend/db"

	"golang.org/x/crypto/bcrypt"
)

// UserSignup creates a new user with email and password
func UserSignup(ctx context.Context, database *db.PrismaClient, email string, password string) (*db.UserModel, error) {
	log.Printf("[SERVICE] UserSignup: Checking if user %s already exists", email)
	// Check if user already exists
	existingUser, err := database.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err == nil && existingUser != nil {
		log.Printf("[SERVICE] UserSignup: ERROR - User %s already exists", email)
		return nil, errors.New("user already exists")
	}

	log.Printf("[SERVICE] UserSignup: Hashing password for user %s", email)
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[SERVICE] UserSignup: ERROR - Failed to hash password for %s - %v", email, err)
		return nil, errors.New("failed to hash password")
	}

	log.Printf("[SERVICE] UserSignup: Creating user %s in database", email)
	// Create new user
	user, err := database.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Password.Set(string(hashedPassword)),
	).Exec(ctx)

	if err != nil {
		log.Printf("[SERVICE] UserSignup: ERROR - Failed to create user %s - %v", email, err)
		return nil, errors.New("failed to create user")
	}

	log.Printf("[SERVICE] UserSignup: User %s created successfully with ID: %s", email, user.ID)

	// Create a default empty profile for the new user so frontend can rely on /api/profile
	go func() {
		// run in background to avoid slowing signup response; use a background context
		_, err := database.UserProfile.CreateOne(
			db.UserProfile.User.Link(
				db.User.ID.Equals(user.ID),
			),
			db.UserProfile.FirstName.Set(""),
			db.UserProfile.LastName.Set(""),
			db.UserProfile.Phone.Set(""),
			db.UserProfile.Bio.Set(""),
			db.UserProfile.Avatar.Set(""),
			db.UserProfile.Preferences.Set([]byte("{}")),
		).Exec(context.Background())
		if err != nil {
			log.Printf("[SERVICE] UserSignup: WARNING - failed to create default profile for user %s: %v", user.ID, err)
		} else {
			log.Printf("[SERVICE] UserSignup: Default profile created for user %s", user.ID)
		}
	}()

	return user, nil
}

// UserLogin authenticates a user and returns the user if credentials are valid
func UserLogin(ctx context.Context, database *db.PrismaClient, email string, password string) (*db.UserModel, error) {
	log.Printf("[SERVICE] UserLogin: Attempting to find user by email %s", email)
	// Find user by email
	user, err := database.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] UserLogin: ERROR - User %s not found - %v", email, err)
		return nil, errors.New("user not found")
	}

	if user == nil {
		log.Printf("[SERVICE] UserLogin: ERROR - User %s not found (nil result)", email)
		return nil, errors.New("user not found")
	}

	log.Printf("[SERVICE] UserLogin: User %s found, verifying password", email)
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("[SERVICE] UserLogin: ERROR - Invalid password for user %s", email)
		return nil, errors.New("invalid password")
	}

	log.Printf("[SERVICE] UserLogin: Password verified successfully for user %s", email)
	return user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(ctx context.Context, database *db.PrismaClient, userID string) (*db.UserModel, error) {
	log.Printf("[SERVICE] GetUserByID: Fetching user %s", userID)
	user, err := database.User.FindUnique(db.User.ID.Equals(userID)).Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] GetUserByID: ERROR - User %s not found - %v", userID, err)
		return nil, errors.New("user not found")
	}

	log.Printf("[SERVICE] GetUserByID: User %s retrieved successfully", userID)
	return user, nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(ctx context.Context, database *db.PrismaClient, email string) (*db.UserModel, error) {
	log.Printf("[SERVICE] GetUserByEmail: Fetching user by email %s", email)
	user, err := database.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] GetUserByEmail: ERROR - User %s not found - %v", email, err)
		return nil, errors.New("user not found")
	}

	log.Printf("[SERVICE] GetUserByEmail: User %s retrieved successfully", email)
	return user, nil
}

// ListUsers retrieves all users
func ListUsers(ctx context.Context, database *db.PrismaClient) ([]db.UserModel, error) {
	log.Printf("[SERVICE] ListUsers: Fetching all users")
	users, err := database.User.FindMany().Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] ListUsers: ERROR - Failed to fetch users - %v", err)
		return nil, errors.New("failed to fetch users")
	}

	log.Printf("[SERVICE] ListUsers: Retrieved %d users", len(users))
	return users, nil
}

// DeleteUser removes a user from the database
func DeleteUser(ctx context.Context, database *db.PrismaClient, userID string) (*db.UserModel, error) {
	log.Printf("[SERVICE] DeleteUser: Attempting to delete user %s", userID)
	user, err := database.User.FindUnique(db.User.ID.Equals(userID)).Delete().Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] DeleteUser: ERROR - Failed to delete user %s - %v", userID, err)
		return nil, errors.New("failed to delete user")
	}

	log.Printf("[SERVICE] DeleteUser: User %s deleted successfully", userID)
	return user, nil
}
