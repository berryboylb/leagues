package jwt

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateJWT(t *testing.T) {
	// Mock providerID
	providerID := primitive.NewObjectID()

	// Call GenerateJWT function
	tokenString, err := GenerateJWT(providerID)

	// Assert no error
	assert.NoError(t, err)

	// Assert that the token string is not empty
	assert.NotEmpty(t, tokenString)
}
