// Package models holds the data structures that represent rows in our
// database tables. In GORM, a Go struct and a database table are two
// views of the same thing: the struct defines the table's columns.
package models

import "time"

// User represents one row in the "users" table.
//
// The backtick strings after each field are called "struct tags". They
// are metadata that libraries read using reflection. Here we use two:
//
//	json:"..."  -> controls how this field is named (or hidden) when the
//	               struct is converted to/from JSON in API responses.
//	gorm:"..."  -> tells GORM how to build the database column.
//
// GORM also automatically pluralizes the struct name for the table name,
// so "User" -> "users".
type User struct {
	// primaryKey marks this as the table's primary key. GORM auto-increments
	// uint primary keys for you, so you never set the ID manually.
	ID uint `json:"id" gorm:"primaryKey"`

	// unique     -> the database rejects a second user with the same email,
	//               which is how we prevent duplicate signups.
	// not null   -> the column can't be empty.
	Email string `json:"email" gorm:"unique;not null"`

	// json:"-" is the important one here: the dash means "never include this
	// field in JSON output". So even if we accidentally return a whole User
	// in a response, the (hashed) password stays hidden. Security by default.
	Password string `json:"-" gorm:"not null"`

	// CreatedAt is a "magic" field name that GORM recognizes: it fills this
	// in automatically with the current time when the row is first created.
	// (GORM also recognizes UpdatedAt and DeletedAt if you add them.)
	CreatedAt time.Time `json:"created_at"`
}
