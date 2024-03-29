package jwt

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
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

func TestGetUser(t *testing.T) {
	id, err := primitive.ObjectIDFromHex("6606bbc3d5330e500a58558a") //gotten from db
	assert.NoError(t, err)
	user, err := GetUser(id)
	assert.NoError(t, err)

	assert.NotNil(t,user)
}
