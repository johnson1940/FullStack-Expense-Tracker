// Package utils holds small, reusable helper functions that several
// parts of the app share — here, the password hashing logic.
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword turns a plain-text password into a secure hash for storage.
//
// We NEVER store raw passwords. bcrypt is a one-way hash: you can turn a
// password into a hash, but you cannot turn the hash back into the password.
// It's also deliberately slow, which makes brute-force attacks expensive,
// and it automatically mixes in a random "salt" so two users with the same
// password still get different hashes.
func HashPassword(password string) (string, error) {
	// DefaultCost controls how slow/strong the hashing is. The default is a
	// good balance for most apps. The result is a []byte, so we convert it
	// to a string for storing in the database.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies a login attempt against the stored hash.
//
// Notice we never "decrypt" the stored hash. Instead, bcrypt re-hashes the
// attempt and compares it internally. It returns nil on a match and an
// error otherwise, so we translate that into a simple bool for our callers.
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
