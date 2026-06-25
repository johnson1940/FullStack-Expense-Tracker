package models

import "time"

// CreateExpenseInput defines the expected JSON structure for adding/updating an expense.
type CreateExpenseInput struct {
	Title      string    `json:"title" binding:"required"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Date       time.Time `json:"date" binding:"required"`
	CategoryID *uint     `json:"category_id"`
}
