// Package config handles application startup configuration — right now that's
// just loading environment variables from the .env file.
package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv reads key=value pairs from a .env file in the working directory and
// puts them into the process environment, so os.Getenv("DB_HOST") etc. work.
//
// If there's no .env file (for example in production, where you set real
// environment variables directly), that's fine — we just log it and carry on
// using whatever the system already provides.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found — using system environment variables")
	}
}
