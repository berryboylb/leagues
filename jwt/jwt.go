package jwt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"league/db"
	"league/helpers"
	"league/models"
	cisredis "league/redis"

	"github.com/go-redis/redis/v8"
)

var SecretKey []byte
var userCollection *mongo.Collection = db.GetCollection(db.MongoClient, "users")

func init() {
	// database = db.GetDB()
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Error loading jwt secret")
	}
	SecretKey = []byte(secretKey)
}

func GenerateJWT(providerID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = providerID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetSingleUser(ID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	err = userCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(ID string) (models.User, error) {
	userStr, err := cisredis.Retrieve(ID)
	if err != nil {
		if err == redis.Nil {
			// Fetch user data from database
			user, err := GetSingleUser(ID)
			if err != nil {
				return models.User{}, fmt.Errorf("failed to fetch user data from database: %w", err)
			}
			// Store user data in Redis with an expiration time
			userByte, err := cisredis.StoreStruct(user) 
			if err != nil {
				return models.User{}, fmt.Errorf("failed to store user data in Redis: %w", err)
			}
			expiration := 10 * time.Minute
			err = cisredis.Store(ID, userByte, expiration)
			if err != nil {
				return models.User{}, fmt.Errorf("failed to store user data in Redis with expiration: %w", err)
			}
			return *user, nil
		}
		return models.User{}, fmt.Errorf("failed to retrieve user data from Redis: %w", err)
	}

	var user models.User
	if err := cisredis.UnmarshalStruct([]byte(userStr), &user); err != nil {
		return models.User{}, fmt.Errorf("failed to unmarshal user data from Redis: %w", err)
	}
	return user, nil
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.CreateResponse(c, helpers.Response{
				Message:    "invalid jwt",
				StatusCode: http.StatusUnauthorized,
				Data:       nil,
			})
			return
		}

		token, err := jwt.Parse(authHeader[len("Bearer "):], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key
			return SecretKey, nil
		})

		if err != nil {
			helpers.CreateResponse(c, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			helpers.CreateResponse(c, helpers.Response{
				Message:    "invalid jwt",
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}

		providerID := claims["id"].(string)
		user, err := GetUser(providerID)
		if err != nil {
			helpers.CreateResponse(c, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}
		c.Set("claims", claims)
		c.Set("user", user)
		c.Next()
	}
}
