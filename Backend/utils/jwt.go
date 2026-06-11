package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT (JSON Web Token) for a logged-in user.
//
// A JWT is a small, signed string the client stores after logging in and
// sends back on future requests to prove who they are. It has three parts
// separated by dots (header.payload.signature). The payload holds "claims"
// — facts about the user — and the signature is what makes it tamper-proof.
func GenerateToken(userID uint, secret string) (string, error) {
	// Claims are the data we embed in the token. These are standard claim
	// names with specific meanings:
	//   "sub" (subject) -> WHO the token is about — here, the user's ID.
	//   "exp" (expiry)  -> WHEN the token stops being valid (a Unix timestamp).
	//
	// We set expiry to 72 hours from now. After that the token is rejected
	// and the user must log in again — a safety net if a token is stolen.
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	// Build the token using HS256, a symmetric signing algorithm: the same
	// secret both signs the token and later verifies it.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SignedString produces the final token string, signed with our secret.
	// The secret must be kept private (it lives in JWT_SECRET in .env).
	// Anyone with the secret can forge valid tokens, so never expose it.
	return token.SignedString([]byte(secret))
}

// NOTE: This file only CREATES tokens. The matching logic to READ and
// VERIFY a token lives in the auth middleware, which you'll add when you
// start protecting routes (like the expense endpoints).
