package services

import (
	"errors"
	"expense-tracker/database"
	"expense-tracker/models"
	"expense-tracker/utils"
	"os"

	"gorm.io/gorm"
)

// AuthService handles the business logic for authentication.
type AuthService struct{}

// Signup contains the core logic for creating a new user.
func (s *AuthService) Signup(input *models.AuthInput) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash the password
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	// Create the user record
	user := models.User{Email: input.Email, Password: hashed}
	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Login contains the core logic for verifying credentials and issuing a token.
func (s *AuthService) Login(input *models.AuthInput) (string, *models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Use a generic error for security, whether user not found or other DB error.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid email or password")
		}
		return "", nil, err
	}

	// Compare the submitted password against the stored hash.
	if !utils.CheckPassword(user.Password, input.Password) {
		return "", nil, errors.New("invalid email or password")
	}

	// Credentials are valid, issue a token.
	token, err := utils.GenerateToken(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", nil, errors.New("could not create token")
	}

	return token, &user, nil
}
