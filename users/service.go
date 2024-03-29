package users

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"league/db"
	"league/models"
	"league/redis"

	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var userCollection *mongo.Collection = db.GetCollection(db.MongoClient, "users")
// var userCollection *mongo.Collection //for tests
var duration time.Duration = 10 * time.Second

func updateUser(ID string, update UserRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Parse ObjectID from string
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Create update fields
	updates := bson.M{
		"first_name": update.FirstName,
		"last_name":  update.LastName,
		"email":      strings.ToLower(update.Email),
	}

	// Perform the update operation
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		// Handle specific error types
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, e := range mongoErr.WriteErrors {
				if e.Code == 11000 {
					return nil, fmt.Errorf("email already exists: %s", update.Email)
				}
			}
		}
		return nil, fmt.Errorf("could not update user: %v", err)
	}

	// Fetch the updated user from the database
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %v", err)
	}

	// Update user data in Redis
	userByte, err := redis.StoreStruct(user)
	if err != nil {
		return nil, fmt.Errorf("failed to store user data in Redis: %v", err)
	}
	expiration := 10 * time.Minute
	err = redis.Store(objID.Hex(), userByte, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to store user data in Redis with expiration: %v", err)
	}

	return &user, nil
}

func deleteUser(ID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// objId, err := primitive.ObjectIDFromHex(ID)
	// if err != nil {
	// 	return fmt.Errorf("invalid ObjectID: %v", err)
	// }

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no user found with ID %s", ID)
	}

	// Delete user data from Redis
	err = redis.Delete(ID.Hex())
	if err != nil {
		return err
	}

	return nil
}

func getUsers(filters UserRequest, pageNumber string, pageSize string) ([]models.User, int64, int64, int64, error) {
	perPage := int64(15)
	page := int64(1)

	if pageSize != "" {
		if perPageNum, err := strconv.Atoi(pageSize); err == nil {
			perPage = int64(perPageNum)
		}
	}

	if pageNumber != "" {
		if num, err := strconv.Atoi(pageNumber); err == nil {
			page = int64(num)
		}
	}

	offset := (page - 1) * perPage

	filter := bson.M{}

	if filters.Email != "" {
		filter["email"] = filters.Email
	}
	if filters.FirstName != "" {
		filter["first_name"] = filters.FirstName
	}
	if filters.LastName != "" {
		filter["last_name"] = filters.LastName
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	fOpt := options.FindOptions{Limit: &perPage, Skip: &offset, Sort: bson.D{{"created_at", -1}}}
	cOpt := options.CountOptions{Limit: &perPage, Skip: &offset}

	total, err := userCollection.CountDocuments(ctx, filter, &cOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to count users: %v", err)
	}

	cursor, err := userCollection.Find(ctx, filter, &fOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to find users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to decode users: %v", err)
	}

	return users, total, page, perPage, nil
}
