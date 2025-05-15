package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gankun2024/gin-demo-project/internal/db/models"
)

// User represents a user in the system
type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUser creates a new user
func CreateUser(email, password, firstName, lastName string) (*User, error) {
	// Check if user exists
	exists, err := models.UserExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate user ID
	userID := uuid.New().String()

	// Create user in database
	now := time.Now()
	user := &models.User{
		ID:             userID,
		Email:          email,
		HashedPassword: string(hashedPassword),
		FirstName:      firstName,
		LastName:       lastName,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := models.CreateUser(user); err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// AuthenticateUser authenticates a user with email and password
func AuthenticateUser(email, password string) (*User, error) {
	// Get user from database
	user, err := models.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*User, error) {
	// Get user from database
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUser updates a user's profile
func UpdateUser(id, firstName, lastName string) (*User, error) {
	// Get user from database
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Update user fields
	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	// Save user
	if err := models.UpdateUser(user); err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// SendPasswordResetEmail sends a password reset email to the user
func SendPasswordResetEmail(email string) error {
	// Check if user exists
	user, err := models.GetUserByEmail(email)
	if err != nil {
		// Don't reveal that the email doesn't exist
		return nil
	}

	// Generate reset token
	token := uuid.New().String()
	expires := time.Now().Add(time.Hour * 24) // 24 hours

	// Save reset token
	if err := models.SavePasswordResetToken(user.ID, token, expires); err != nil {
		return err
	}

	// In a real application, you would send an email with the token
	// For now, we'll just log it (implementation would depend on your email provider)
	// log.Printf("Password reset token for %s: %s", email, token)

	return nil
}

// ResetPassword resets a user's password using a token
func ResetPassword(token, newPassword string) error {
	// Verify token
	userID, err := models.VerifyPasswordResetToken(token)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user's password
	if err := models.UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		return err
	}

	// Invalidate token
	if err := models.InvalidatePasswordResetToken(token); err != nil {
		return err
	}

	return nil
}
