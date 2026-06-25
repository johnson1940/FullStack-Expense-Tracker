// Package handlers contains the functions that actually handle incoming
// HTTP requests. Each handler reads the request, does some work (talk to
// the service), and writes back a JSON response.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"expense-tracker/models"
	"expense-tracker/services"
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
// This is now defined in the models package to be shared.

// We create an instance of our service. In a larger app, this would be
// "injected" for better testability, but this is fine for now.
var authService = services.AuthService{}

// Signup creates a new user account.
// In Gin, every handler takes a *gin.Context ("c"), which bundles together
// the incoming request and the tools to write a response.
func Signup(c *gin.Context) {
	var input models.AuthInput

	// ShouldBindJSON parses the request body into our struct AND runs the
	// validation rules from the binding tags. If anything is wrong, we
	// reply 400 Bad Request and stop (the early `return` is important —
	// without it the function would keep running).
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to perform the business logic
	user, err := authService.Signup(&input)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var input models.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to perform the business logic
	token, user, err := authService.Login(&input)
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "data": gin.H{"token": token, "user": user}})
}
