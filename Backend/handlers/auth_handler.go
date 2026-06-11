// Package handlers contains the functions that actually handle incoming
// HTTP requests. Each handler reads the request, does some work (talk to
// the database, hash a password, etc.), and writes back a JSON response.
package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	// Replace with your own module path (see note in main.go).
	"expense-tracker/database"
	"expense-tracker/models"
	"expense-tracker/utils"
)

// AuthInput defines the shape of the JSON we EXPECT to receive for both
// signup and login. Defining a dedicated input struct (instead of reading
// the User model directly) is good practice: it controls exactly which
// fields a client is allowed to send.
//
// The `binding` tags tell Gin to validate the data automatically:
//
//	required  -> the field must be present
//	email     -> must look like a valid email address
//	min=6     -> the password must be at least 6 characters
//
// If validation fails, ShouldBindJSON returns an error before our logic runs.
type AuthInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Signup creates a new user account.
// In Gin, every handler takes a *gin.Context ("c"), which bundles together
// the incoming request and the tools to write a response.
func Signup(c *gin.Context) {
	var input AuthInput

	// ShouldBindJSON parses the request body into our struct AND runs the
	// validation rules from the binding tags. If anything is wrong, we
	// reply 400 Bad Request and stop (the early `return` is important —
	// without it the function would keep running).
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before it ever touches the database.
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// Build the user record and insert it. database.DB.Create runs the
	// SQL INSERT. We check .Error to see if it failed — the most likely
	// failure is the unique-email constraint firing because the email is
	// already registered, so we return 409 Conflict.
	user := models.User{Email: input.Email, Password: hashed}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// 201 Created is the correct status for "a new resource was made".
	// gin.H is just a shortcut for map[string]interface{} — a quick way to
	// build a JSON object. Wrapping data in {"data": ...} gives every
	// response a consistent shape, which makes the Flutter side simpler.
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully signed up", "data": user})
}

// Login verifies credentials and returns a token on success.
func Login(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Look up the user by email. .First fills `user` with the first match
	// and returns an error if no row is found.
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// SECURITY NOTE: we deliberately return the SAME generic message
		// whether the email doesn't exist OR the password is wrong. If we
		// said "no such email" specifically, an attacker could use that to
		// discover which emails are registered.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Compare the submitted password against the stored hash.
	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Credentials are valid — issue a token. os.Getenv reads the secret
	// we loaded from .env. This token is what Flutter will save and send
	// on every future request to prove the user is logged in.
	token, err := utils.GenerateToken(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "data": gin.H{"token": token, "user": user}})
}
