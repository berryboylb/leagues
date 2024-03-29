package users

import (
	"context"

	"testing"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBURI    = "mongodb+srv://admin:test1234@test.ct3433r.mongodb.net/?retryWrites=true&w=majority"
	testDatabase = "real"
)

func setupTestEnvironment(t *testing.T) {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(testDBURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Set up the database and collection
	userCollection = client.Database(testDatabase).Collection("users")
}

func cleanupTestEnvironment(t *testing.T) {
	// Drop the test collection after tests are completed
	// err := userCollection.Drop(context.Background())
	// if err != nil {
	// 	t.Fatalf("Error dropping test collection: %v", err)
	// }

	// Disconnect from MongoDB
	err := userCollection.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

// TestUpdateUser_Success tests the updateUser function.
func TestUpdateUser_Success(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Define update fields
	update := UserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
	}

	// Call updateUser function
	updatedUser, err := updateUser("6606bbc3d5330e500a58558a", update) //gotten from db
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)

	// Verify that user fields are updated correctly
	assert.Equal(t, update.FirstName, updatedUser.FirstName)
	assert.Equal(t, update.LastName, updatedUser.LastName)
	assert.Equal(t, update.Email, updatedUser.Email)
}

// TestDeleteUser_Success tests the deleteUser function.
func TestDeleteUser_Success(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	ID, err := primitive.ObjectIDFromHex("6606bbc3d5330e500a58558a") //gotten from db
	assert.NoError(t, err)

	// Call deleteUser function
	err = deleteUser(ID)
	assert.NoError(t, err)
}

// TestGetUsers_Success tests the getUsers function.
func TestGetUsers_Success(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call getUsers function
	filters := UserRequest{}
	pageNumber := "1"
	pageSize := "10"
	users, total, _, _, err := getUsers(filters, pageNumber, pageSize)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Greater(t, total, int64(0))
}
