package auth

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"league/models"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	os.Setenv("MONGODB_URI", "mongodb+srv://admin:test1234@test.ct3433r.mongodb.net/league?retryWrites=true&w=majority")
}

var (
	testDBURI    = "mongodb://localhost:27017/testdb"
	testDatabase = "testdb"
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
	err := userCollection.Drop(context.Background())
	if err != nil {
		t.Fatalf("Error dropping test collection: %v", err)
	}

	// Disconnect from MongoDB
	err = userCollection.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

func TestGenerateOtp(t *testing.T) {
	// Remove redundant code, such as setting MONGODB_URI here, it's set in init() already
	fmt.Println("Testing GenerateOtp function")

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

type MockUserCollection struct{}



func TestCreateUser_Success(t *testing.T) {

	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		RoleName:  "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedUser, err := createUser(user)

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
	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
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

	email := "existing@example.com"

	user, err := getUserByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	email := "nonexisting@example.com"

	user, err := getUserByEmail(email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, fmt.Sprintf("user with the email %v is not found", email))
}

func TestSendOtp_Success(t *testing.T) {
	user := models.User{
		Email:     "johndoe@example.com",
	}

	// Call the sendOtp function
	sendOtp(&user, "Verification")
}

func TestGetUserFromOtp_ValidToken(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)

	// Insert a user with a valid verification token
	// Insert a user with an expired verification token
	expiredTime := time.Now().Add(time.Hour) // Set token expiration time in the past
	user := models.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john.doe@example.com",
		RoleName:          "user",
		VerificationToken: "4444",
		ExpiresAt:         expiredTime,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Insert the user into the test database
	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		t.Fatalf("Error inserting user with valid token: %v", err)
	}

	// Call getUserFromOtp with the valid token and email
	foundUser, err := getUserFromOtp(user.VerificationToken, user.Email)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	// Clean up test environment
	cleanupTestEnvironment(t)
}


func TestGetUserFromOtp_ExpiredToken(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Insert a user with an expired verification token
	expiredTime := time.Now().Add(-time.Hour) // Set token expiration time in the past
	user := models.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john.doe@example.com",
		RoleName:          "user",
		VerificationToken: "4444",
		ExpiresAt:         expiredTime,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Insert the user into the test database
	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		t.Fatalf("Error inserting user with expired token: %v", err)
	}

	// Call getUserFromOtp with the expired token and email
	foundUser, err := getUserFromOtp(user.VerificationToken, user.Email)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "not found")
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
	user, err := getUser("nonexistingOTP", "nonexisting@example.com")

	// Assert that the function returns an error indicating user not found
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user with token nonexistingOTP not found")
}



func TestChangePassword_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	//existing user

	existinguser := models.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john.doe@example.com",
		RoleName:          "user",
	}

	err := changePassword(existinguser.Id, "newpassword")

	// Assert that the function returns no error
	assert.NoError(t, err)

}
