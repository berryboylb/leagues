package auth

import (
	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"league/db"
	"league/emails"
	"league/helpers"
	"league/models"

	mrand "math/rand"

	

	"context"
	"fmt"
	"log"
	"os"
	"time"

	"strings"
)

var adminFirstName string
var adminLastName string
var adminEmail string
var adminPassword string

var userCollection *mongo.Collection = db.GetCollection(db.MongoClient, "users")
var duration time.Duration = 10 * time.Second

func init() {
	//check for email index
	exists, err := db.IsIndexExists(context.Background(), userCollection, "email")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !exists {
		err = db.IndexField(*userCollection, "email", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	// err = godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }
	adminFirstName = os.Getenv("ADMIN_FIRST_NAME")
	adminLastName = os.Getenv("ADMIN_FIRST_NAME")
	adminEmail = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")

	if adminFirstName == "" || adminEmail == "" || adminPassword == "" || adminLastName == "" {
		log.Fatal("Error loading super admin details")
	}
	createAdmin()
}




func createAdmin() {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": strings.ToLower(adminEmail)}).Decode(&user)
	if err != nil {
		hash, _ := helpers.HashPassword(adminPassword, 8)
		if err == mongo.ErrNoDocuments {
			newUser := models.User{
				FirstName: adminFirstName,
				LastName:  adminLastName,
				Email:     adminEmail,
				RoleName:  models.SuperAdminRole,
				Password:  hash,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_, err := createUser(newUser)
			if err != nil {
				log.Fatal("Error creating super admin details")
				return
			}
			fmt.Println("Super admin created successfully")
		}
	}
	if user.Id != primitive.NilObjectID {
		fmt.Println("Super admin already exists")
	}
}

func createUser(user models.User) (*models.User, error) {
	user.SetEmail()
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		//check for duplicates
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, e := range mongoErr.WriteErrors {
				if e.Code == 11000 { // 11000 is the code for duplicate key error
					return nil, fmt.Errorf("email already exists: %s", user.Email)
				}
			}
		}
		return nil, err
	}
	var insertedUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedUser)
	if err != nil {
		// Handle error
		return nil, fmt.Errorf("failed to fetch inserted user: %v", err)
	}

	return &insertedUser, nil
}

func getUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": strings.ToLower(email)}).Decode(&user)
	if err != nil {
		// Handle error
		if err == mongo.ErrNoDocuments {
			// If no user found with the specified email
			return nil, fmt.Errorf("user with the email %v is not found", email)
		}
		// If other error occurred
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	return &user, nil
}

const charset = "0123456789"

var seededRand *mrand.Rand = mrand.New(
	mrand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	fmt.Println(string(b), "otp")
	return string(b)
}

func GenerateOtp(length int) string {
	return StringWithCharset(length, charset)
}

func sendOtp(user *models.User, title string) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	otp := GenerateOtp(4)
	expiryTime := time.Now().Add(24 * time.Hour) //  24 hours
	update := bson.M{"verification_token": otp, "expires_at": expiryTime, "updated_at": time.Now()}

	_, err := userCollection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": update})
	if err != nil {
		fmt.Printf("could not update user: %v \n", err)
		return
	}

	err = emails.SendOTPEmail(user.Email, otp, title)
	if err != nil {
		fmt.Printf("could not send email: %v \n", err)
		return
	}
	fmt.Printf("successfully sent otp %v \n", otp)
}

func getUserFromOtp(otp string, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	var user models.User
	defer cancel()
	filter := bson.M{
		"email": strings.ToLower(email),
		// "verification_token": otp,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// Handle error
		if err == mongo.ErrNoDocuments {
			// If no user found with the specified verification token
			return nil, fmt.Errorf("user with token %s not found", otp)
		}
		// If other error occurred
		return nil, fmt.Errorf("failed to fetch user by token: %v", err)
	}

	if user.VerificationToken != otp {
		return nil, fmt.Errorf("invalid verification token")
	}

	return &user, nil
}

func destroyToken(ID primitive.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	update := bson.M{"verification_token": "", "expires_at": time.Time{}, "updated_at": time.Now()}
	_, err := userCollection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": update})
	if err != nil {
		fmt.Printf("could not update user: %v \n", err)
		return
	}
	fmt.Printf("destroyed user token %v \n", ID)
}

func forgotPassword(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": strings.ToLower(email)}).Decode(&user)
	if err != nil {
		// Handle error
		if err == mongo.ErrNoDocuments {
			// If no user found with the specified email
			return fmt.Errorf("user with the email %v is not found", email)
		}
		// If other error occurred
		return fmt.Errorf("failed to fetch user: %v", err)
	}
	return nil
}

func getUser(otp string, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	var user models.User
	defer cancel()
	filter := bson.M{
		"email": strings.ToLower(email),
		// "verification_token": otp,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// Handle error
		if err == mongo.ErrNoDocuments {
			// If no user found with the specified verification token
			return nil, fmt.Errorf("user with token %s not found", otp)
		}
		// If other error occurred
		return nil, fmt.Errorf("failed to fetch user by token: %v", err)
	}

	if user.VerificationToken != otp {
		return nil, fmt.Errorf("invalid verification token")
	}

	return &user, nil
}

func changePassword(ID primitive.ObjectID, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	hash, err := helpers.HashPassword(password, 8)
	if err != nil {
		return err
	}
	update := bson.M{
		"password":   hash,
		"updated_at": time.Now(),
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": update})
	if err != nil {
		fmt.Printf("could not update user: %v \n", err)
		return err
	}
	return nil
}
