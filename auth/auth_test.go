package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"league/models"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBURI    = ""
	testDatabase = "real"
)

func createEmailIndex(client *mongo.Client, databaseName, collectionName string) error {
	// Specify the index model
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},            // Create an ascending index on the email field
		Options: options.Index().SetUnique(true), // Set the unique constraint on the email field
	}

	// Create the index
	_, err := client.Database(databaseName).Collection(collectionName).Indexes().CreateOne(
		context.Background(),
		indexModel,
	)
	if err != nil {
		return err
	}
	return nil
}

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

func TestGenerateOtp(t *testing.T) {
	testCases := []struct {
		length      int
		expectedLen int
	}{
		{5, 5},
		{8, 8},
		{10, 10},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Length %d", tc.length), func(t *testing.T) {
			otp := GenerateOtp(tc.length)
			if len(otp) != tc.expectedLen {
				t.Errorf("Expected OTP length %d, but got %d", tc.expectedLen, len(otp))
			}
		})
	}
}

func TestCreateUser_Success(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john11@example.com",
		RoleName:  "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedUser, err := createUser(user)
	fmt.Println(insertedUser)

	assert.NoError(t, err)
	assert.NotNil(t, insertedUser)
	assert.NotNil(t, insertedUser.Id)
	assert.Equal(t, user.FirstName, insertedUser.FirstName)
	assert.Equal(t, user.LastName, insertedUser.LastName)
	assert.Equal(t, user.Email, insertedUser.Email)
	assert.Equal(t, user.RoleName, insertedUser.RoleName)
	assert.NotEmpty(t, insertedUser.CreatedAt)
	assert.NotEmpty(t, insertedUser.UpdatedAt)
}

func TestCreateUser_DuplicateKeyError(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	//we are using the same email to get a duplicate key error on purpose
	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john11@example.com",
		RoleName:  "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedUser, err := createUser(user)
	assert.Error(t, err)
	assert.Nil(t, insertedUser)
	assert.EqualError(t, err, fmt.Sprintf("email already exists: %s", user.Email))
}

func TestGetUserByEmail_Success(t *testing.T) {

	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	user, err := getUserByEmail("john11@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john11@example.com", user.Email)
}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	user, err := getUserByEmail("nonexisting@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, fmt.Sprintf("user with the email %v is not found", "nonexisting@example.com"))
}

func TestSendOtp_Success(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Fetch a user from the database
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": "john11@example.com"}).Decode(&user)
	assert.NoError(t, err)

	assert.NotNil(t, user)

	sendOtp(&user, "Verification")
}

func TestGetUserFromOtp_ValidToken(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	token := GenerateOtp(4)
	email := "john11@example.com"
	update := bson.M{"verification_token": token, "expires_at": time.Now().Add(24 * time.Hour), "updated_at": time.Now()}

	_, err := userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": update})
	assert.NoError(t, err)

	// Call getUserFromOtp with the valid token and email
	foundUser, err := getUserFromOtp(token, email)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, email, foundUser.Email)
	assert.Equal(t, token, foundUser.VerificationToken)
}

func TestGetUserFromOtp_ExpiredToken(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// we purpoely gonna backdate it
	token := GenerateOtp(4)
	email := "john11@example.com"
	update := bson.M{"verification_token": token, "expires_at": time.Now().Add(-48 * time.Hour), "updated_at": time.Now()}

	_, err := userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": update})
	assert.NoError(t, err)

	// Call getUserFromOtp with the valid token and email
	foundUser, err := getUserFromOtp(token, email)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestForgotPassword_UserNotFound(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the function with a non-existing email
	err := forgotPassword("nonexisting@example.com")

	// Assert that the function returns an error indicating user not found
	assert.Error(t, err)
	assert.EqualError(t, err, "user with the email nonexisting@example.com is not found")
}

func TestGetUser_UserNotFound(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the function with a non-existing OTP and email combination
	user, err := getUser("3758", "nonexisting@example.com")

	// Assert that the function returns an error indicating user not found
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user with token 3758 not found")
}



func TestChangePassword_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	//existing user

	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": "john11@example.com"}).Decode(&user)
	assert.NoError(t, err)

	err = changePassword(user.Id, "123456789")
		assert.NoError(t, err)

	// Assert that the function returns no error
	assert.NoError(t, err)

}
