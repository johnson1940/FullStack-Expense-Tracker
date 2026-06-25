// Package routes maps URL paths to the handler functions that serve them.
// Keeping all routes in one file lets you see your entire API surface at a
// single glance — a useful "table of contents" for the backend.
package routes

import (
	"github.com/gin-gonic/gin"

	// Replace with your own module path (see note in main.go).
	"expense-tracker/handlers"
	"expense-tracker/middleware"
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

	// 1. Public Auth Routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
	}

	// 2. Protected Expense Routes
	// middleware.AuthRequired() ensures only logged-in users with a valid JWT can enter
	expenses := r.Group("/expenses")
	expenses.Use(middleware.AuthRequired())
	{
		expenses.POST("/", handlers.CreateExpense) // POST /expenses/
		expenses.GET("/", handlers.ListExpenses)
		expenses.PUT("/:id", handlers.UpdateExpense)
		expenses.DELETE("/:id", handlers.DeleteExpense)
	}

	// 3. Protected Category Routes
	categories := r.Group("/categories")
	categories.Use(middleware.AuthRequired())
	{
		categories.GET("/", handlers.ListCategories)
		categories.POST("/", handlers.CreateCategory)
	}
}
