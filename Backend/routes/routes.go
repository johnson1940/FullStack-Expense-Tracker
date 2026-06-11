// Package routes maps URL paths to the handler functions that serve them.
// Keeping all routes in one file lets you see your entire API surface at a
// single glance — a useful "table of contents" for the backend.
package routes

import (
	"github.com/gin-gonic/gin"

	// Replace with your own module path (see note in main.go).
	"expense-tracker/handlers"
)

// Setup attaches all our routes to the router. main() passes the router
// (r) in, and we hang our endpoints off it.
func Setup(r *gin.Engine) {
	// A simple health-check endpoint. Hitting GET /health is the quickest
	// way to confirm the server is alive without touching the database.
	// The inline func here is a tiny handler defined right where it's used.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// A route group lets related endpoints share a common URL prefix.
	// Everything inside this group is automatically prefixed with "/auth",
	// so the full paths become POST /auth/signup and POST /auth/login.
	//
	// Groups are also where you'd later attach middleware to a whole set of
	// routes at once — for example, wrapping your expense routes in a JWT
	// auth check so only logged-in users can reach them.
	auth := r.Group("/auth")
	{
		// We pass the handler function itself (no parentheses) — Gin calls
		// it for us whenever a matching request arrives.
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
	}
}
