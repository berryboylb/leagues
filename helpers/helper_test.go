package helpers

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	// Define a password and cost for hashing
	password := "testPassword"
	cost := bcrypt.DefaultCost

	// Call HashPassword function
	hashedPassword, err := HashPassword(password, cost)

	// Assert no error
	assert.NoError(t, err)

	// Assert that the hashed password is not empty
	assert.NotEmpty(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	// Define a password and cost for hashing
	password := "testPassword"
	cost := bcrypt.DefaultCost

	// Call HashPassword function
	hashedPassword, _ := HashPassword(password, cost)

	// Call CheckPasswordHash function
	match := CheckPasswordHash(password, hashedPassword)

	// Assert that the password and hash match
	assert.True(t, match)
}
