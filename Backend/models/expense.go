package models

import "time"

// Expense represents one row in the "expenses" table.
type Expense struct {
	ID     uint      `json:"id" gorm:"primaryKey"`
	Title  string    `json:"title" gorm:"not null"`
	Amount float64   `json:"amount" gorm:"not null"`
	Date   time.Time `json:"date" gorm:"not null"`

	// Foreign key for the User who owns this expense
	UserID uint `json:"user_id" gorm:"not null"`

	// Foreign key for the Category (can be NULL)
	CategoryID *uint `json:"category_id"`
	// The actual Category data, populated by GORM's Preload
	Category *Category `json:"category"`

	CreatedAt time.Time `json:"created_at"`
}
