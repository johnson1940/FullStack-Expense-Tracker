// Package database owns everything about the database connection.
// Keeping this in one place means there's exactly one spot to look at
// (or change) when something is wrong with the DB, or when you want to
// swap Postgres for another database later.
package database

import (
	"expense-tracker/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// Replace with your own module path (see note in main.go).
)

// DB is a package-level variable holding the live database connection.
// It's exported (capital D) so other packages — like our handlers —
// can use it as database.DB. Using a shared connection like this is the
// normal pattern: GORM manages a pool of connections under the hood, so
// you create it once and reuse it everywhere.
var DB *gorm.DB

// Connect opens the database connection and runs migrations.
// We call this once, from main(), at startup.
func Connect() {
	// The DSN (Data Source Name) is a single string that tells GORM
	// where Postgres lives and how to authenticate. We build it from
	// environment variables so secrets like the password never get
	// hardcoded into the source file.
	//
	//   sslmode=disable        -> fine for local dev; enable it in production
	//   TimeZone=Asia/Kolkata  -> makes timestamps use your local time
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// gorm.Open establishes the connection. The postgres.Open(dsn) part
	// selects the Postgres "driver" — if you ever switch databases, this
	// is essentially the only line that changes.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Fatal prints the message and immediately exits the program.
		// There's no point continuing if we can't reach the database.
		log.Fatal("failed to connect to database:", err)
	}

	// AutoMigrate inspects the User struct and creates/updates the matching
	// table to fit it. IMPORTANT: it creates TABLES inside the database, but
	// it does NOT create the database itself — that has to already exist
	// (you create it once with `CREATE DATABASE expense_tracker;`).
	//
	// As you add more models later (Expense, Category), you list them all
	// here: db.AutoMigrate(&models.User{}, &models.Expense{}, ...)
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Expense{}, // This will now use the updated struct
	); err != nil {
		log.Fatal("failed to run database migrations:", err)
	}

	// Seed the database with default categories
	seedCategories(db)

	// Store the connection in our package-level variable so the rest of
	// the app can use it.
	DB = db
	log.Println("database connected & migrated")
}

// seedCategories populates the database with 18 default categories if the table is empty.
func seedCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)

	if count == 0 {
		defaultCategories := []models.Category{
			{Name: "Housing"},
			{Name: "Utilities"},
			{Name: "Groceries"},
			{Name: "Dining Out"},
			{Name: "Transportation"},
			{Name: "Insurance"},
			{Name: "Healthcare"},
			{Name: "Debt/Loans"},
			{Name: "Entertainment"},
			{Name: "Personal Care"},
			{Name: "Education"},
			{Name: "Clothing/Shopping"},
			{Name: "Savings/Investments"},
			{Name: "Travel"},
			{Name: "Gifts/Donations"},
			{Name: "Subscriptions"},
			{Name: "Pets"},
			{Name: "Miscellaneous"},
		}
		if err := db.Create(&defaultCategories).Error; err != nil {
			log.Println("failed to seed default categories:", err)
		} else {
			log.Println("successfully seeded 18 default categories")
		}
	}
}
