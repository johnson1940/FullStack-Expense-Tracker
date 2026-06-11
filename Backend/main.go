// Package main is the entry point of the backend. Running `go run .` starts
// here. Its job is to wire the pieces together in the right order and then
// start listening for HTTP requests.
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"expense-tracker/config"
	"expense-tracker/database"
	"expense-tracker/routes"
)

func main() {
	// 1. Load configuration from .env into the environment. We do this FIRST
	//    because the next steps (DB connection, JWT secret) read from it.
	config.LoadEnv()

	// 2. Open the database connection and run migrations. After this returns,
	//    database.DB is ready for the handlers to use.
	database.Connect()

	// 3. Create the Gin engine — the HTTP router/server. gin.Default() comes
	//    with logging and crash-recovery middleware already attached.
	r := gin.Default()

	// 4. Register all our routes (GET /health, POST /auth/signup, /auth/login).
	routes.Setup(r)

	// 5. Start the server. We read the port from the environment so it can be
	//    changed without touching code, defaulting to 8080 for local dev.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
