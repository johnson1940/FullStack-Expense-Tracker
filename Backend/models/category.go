package models

import "time"

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`

	// UserID is a pointer (*uint) so it can be NULL in the database.
	// If NULL, it's a default category for everyone. If set, it's a custom user category.
	UserID *uint `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
}
