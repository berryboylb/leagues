package teams

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"league/db"
	"league/models"

	"context"
	"fmt"
	"strconv"
	"time"
)

var teamCollection *mongo.Collection = db.GetCollection(db.MongoClient, "teams")
var trophyCollection *mongo.Collection = db.GetCollection(db.MongoClient, "trophies")
var playerCollection *mongo.Collection = db.GetCollection(db.MongoClient, "players")
var duration time.Duration = 10 * time.Second

func init() {
	indexExists, err := db.IsIndexExists(context.Background(), trophyCollection, "name")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !indexExists {
		err = db.IndexField(*trophyCollection, "name", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}
	//check for name index
	nameIndexExists, err := db.IsIndexExists(context.Background(), teamCollection, "name")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !nameIndexExists {
		err = db.IndexField(*teamCollection, "name", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	//check for stadium index
	stadiumIndexExists, err := db.IsIndexExists(context.Background(), teamCollection, "stadium")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !stadiumIndexExists {
		err = db.IndexField(*teamCollection, "stadium", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}
}

func createTeam(team models.Team) (*models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	result, err := teamCollection.InsertOne(ctx, team)
	if err != nil {
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, e := range mongoErr.WriteErrors {
				if e.Code == 11000 {
					return nil, fmt.Errorf("name  %s or stadium %s already exists", team.Name, team.Stadium)
				}
			}
		}
		return nil, err
	}

	var insertedTeam models.Team
	err = teamCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedTeam)
	if err != nil {
		// Handle error
		return nil, fmt.Errorf("failed to fetch inserted user: %v", err)
	}

	return &insertedTeam, nil
}

func getSingleTeam(id string) (*TeamWithCreator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"_id", objID}}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "created_by"},
			{"foreignField", "_id"},
			{"as", "created_by"},
		}}},
		{{"$unwind", "$created_by"}},
	}

	cursor, err := teamCollection.Aggregate(ctx, pipeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team with the id %v is not found", id)
		}
		return nil, fmt.Errorf("failed to execute aggregation pipeline: %v", err)
	}

	var team TeamWithCreator
	if cursor.Next(ctx) {
		if err := cursor.Decode(&team); err != nil {
			return nil, fmt.Errorf("failed to decode team: %v", err)
		}
	} else {
		// If no documents are returned by the pipeline
		return nil, fmt.Errorf("team with the id %v is not found", id)
	}

	return &team, nil
}

func deleteTeam(ID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %v", err)
	}

	result, err := teamCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return fmt.Errorf("failed to delete team: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no team found with ID %s", ID)
	}

	return nil
}

func updateUser(ID string, update TeamRequest) (*models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Parse ObjectID from string
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Create update fields
	updates := bson.M{
		"name":         update.Name,
		"country":      update.Country,
		"state":        update.State,
		"founded_year": update.FoundedYear,
		"stadium":      update.Stadium,
		"sponsor":      update.Sponsor,
	}

	// Perform the update operation
	_, err = teamCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		// Handle specific error types
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, e := range mongoErr.WriteErrors {
				if e.Code == 11000 {
					return nil, fmt.Errorf("name  %v or stadium %v already exists", true, true)
				}
			}
		}
		return nil, fmt.Errorf("could not update user: %v", err)
	}

	// Fetch the updated user from the database
	var team models.Team
	err = teamCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&team)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %v", err)
	}
	return &team, nil
}

func getTeam(filters TeamRequest, pageNumber string, pageSize string) ([]models.Team, int64, int64, int64, error) {
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
	if filters.Name != "" {
		filter["name"] = filters.Name
	}
	if filters.Country != "" {
		filter["country"] = filters.Country
	}
	if filters.State != "" {
		filter["state"] = filters.State
	}
	if filters.FoundedYear != 0 {
		filter["founded_year"] = filters.FoundedYear
	}
	if filters.Stadium != "" {
		filter["stadium"] = filters.Stadium
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	fOpt := options.FindOptions{Limit: &perPage, Skip: &offset, Sort: bson.D{{"created_at", -1}}}
	cOpt := options.CountOptions{Limit: &perPage, Skip: &offset}

	total, err := teamCollection.CountDocuments(ctx, filter, &cOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to count users: %v", err)
	}

	cursor, err := teamCollection.Find(ctx, filter, &fOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to find users: %v", err)
	}
	defer cursor.Close(ctx)

	var teams []models.Team
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to decode users: %v", err)
	}

	return teams, total, page, perPage, nil
}
