package auth

import (
	"fmt"
	// "reflect"
	"league/helpers"
	"github.com/joho/godotenv"
	"log"
	"os"
	"league/models"
	"testing"
	"time"
)

//load env
func TestMain(m *testing.M) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run tests
	exitCode := m.Run()

	// Exit with the same exit code as the test
	os.Exit(exitCode)
}

// test the create user
func TestCreateUser(t *testing.T) {
	dbHost := os.Getenv("MONGO_URI")

	fmt.Printf("DB_HOST: %s\n", dbHost)

	hash, _ := helpers.HashPassword("password", 8)
	newUser := models.User{
		FirstName: "john",
		LastName:  "doe",
		Email:     "johndoerand@gmail.com",
		RoleName:  models.UserRole,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err := createUser(newUser)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected a response, but got nil")
	}

	//do the rest of your unit test
}
