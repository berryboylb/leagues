package fixtures

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"league/db"
	"league/models"

	"crypto/rand"
	"encoding/base64"

	"context"
	"fmt"
	"strconv"
	"time"
)

var competitionCollection *mongo.Collection = db.GetCollection(db.MongoClient, "competitions")
var fixtureCollection *mongo.Collection = db.GetCollection(db.MongoClient, "fixtures")
var duration time.Duration = 10 * time.Second

func init() {
	//check for unique link index
	exists, err := db.IsIndexExists(context.Background(), fixtureCollection, "unique_link")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !exists {
		err = db.IndexField(*fixtureCollection, "unique_link", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	statusexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "status")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !statusexists {
		err = db.IndexSparse(*fixtureCollection, "status", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	homeexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "home_team_id")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !homeexists {
		err = db.IndexNormalField(*fixtureCollection, "home_team_id", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	awayexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "away_team_id")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !awayexists {
		err = db.IndexNormalField(*fixtureCollection, "away_team_id", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	compexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "competition_id")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !compexists {
		err = db.IndexNormalField(*fixtureCollection, "competition_id", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	staexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "stadium")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !staexists {
		err = db.IndexNormalField(*fixtureCollection, "stadium", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	refexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "referee")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !refexists {
		err = db.IndexNormalField(*fixtureCollection, "referee", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}

	dateexists, err := db.IsIndexExists(context.Background(), fixtureCollection, "date")
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return
	}
	if !dateexists {
		err = db.IndexNormalField(*fixtureCollection, "date", 1)
		if err != nil {
			fmt.Println("Failed to index:", err)
			return
		}
	}
}

func generateRandomString(length int) (string, error) {
	numBytes := (length * 6) / 8

	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	randomString = randomString[:length]

	return randomString, nil
}

func createFixture(fixture models.Fixture) (*models.Fixture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	result, err := fixtureCollection.InsertOne(ctx, fixture)
	if err != nil {
		//check for duplicates
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, e := range mongoErr.WriteErrors {
				if e.Code == 11000 { // 11000 is the code for duplicate key error
					return nil, fmt.Errorf("link already exists: %s", fixture.UniqueLink)
				}
			}
		}
		return nil, err
	}
	var inserted models.Fixture
	err = fixtureCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&inserted)
	if err != nil {
		// Handle error
		return nil, fmt.Errorf("failed to fetch inserted fixture: %v", err)
	}

	return &inserted, nil
}

func getFixturesByStatus(status string, pageNumber string, pageSize string) ([]models.Fixture, int64, int64, int64, error) {
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

	filter := bson.M{"status": status}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	fOpt := options.FindOptions{Limit: &perPage, Skip: &offset, Sort: bson.D{{"created_at", -1}}}
	cOpt := options.CountOptions{Limit: &perPage, Skip: &offset}

	total, err := fixtureCollection.CountDocuments(ctx, filter, &cOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to count fixtures : %v", err)
	}

	cursor, err := fixtureCollection.Find(ctx, filter, &fOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to find fixtures: %v", err)
	}
	defer cursor.Close(ctx)

	var fixture []models.Fixture
	if err := cursor.All(ctx, &fixture); err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to decode fixtures: %v", err)
	}

	return fixture, total, page, perPage, nil
}


func getFixtures(query SearchFeaturesRequest, pageNumber string, pageSize string) ([]Fixture, int64, int64, int64, error) {
	perPage := int64(15)
	page := int64(1)

	if pageSize != "" {
		if perPageNum, err := strconv.Atoi(pageSize); err == nil {
			perPage = int64(perPageNum)
		}
		if perPage > 100 {
			perPage = 100
		}
	}

	if pageNumber != "" {
		if num, err := strconv.Atoi(pageNumber); err == nil {
			page = int64(num)
		}

	}

	filter := bson.M{}

	if query.Query != "" {
		regexQuery := bson.M{"$regex": query.Query, "$options": "i"}
		filter["$or"] = bson.A{
			bson.M{"unique_link": regexQuery},
			bson.M{"status": regexQuery},
			bson.M{"stadium": regexQuery},
			bson.M{"referee": regexQuery},
		}
	}

	// Check if query.From is not zero, then add it to the filter
	if !query.From.IsZero() {
		filter["created_at"] = bson.M{"$gte": query.From}
	}

	if !query.To.IsZero() {
		if _, exists := filter["created_at"]; exists {
			filter["created_at"].(bson.M)["$lte"] = query.To
		} else {
			filter["created_at"] = bson.M{"$lte": query.To}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$sort", bson.D{{"created_at", -1}}}},
		{{"$skip", (page - 1) * perPage}},
		{{"$limit", perPage}},
		{{"$lookup", bson.D{
			{"from", "competitions"},
			{"localField", "competition_id"},
			{"foreignField", "_id"},
			{"as", "competition_id"},
		}}},
		{{"$unwind", "$competition_id"}},
		{{"$lookup", bson.D{
			{"from", "teams"},
			{"localField", "home_team_id"},
			{"foreignField", "_id"},
			{"as", "home_team_id"},
		}}},
		{{"$unwind", "$home_team_id"}},
		{{"$lookup", bson.D{
			{"from", "teams"},
			{"localField", "away_team_id"},
			{"foreignField", "_id"},
			{"as", "away_team_id"},
		}}},
		{{"$unwind", "$away_team_id"}},
	}

	cursor, err := fixtureCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to aggregate fixtures: %v", err)
	}
	defer cursor.Close(ctx)

	var fixture []Fixture
	if err := cursor.All(ctx, &fixture); err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to decode fixtures: %v", err)
	}

	// Count total documents
	total, err := fixtureCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to count fixtures: %v", err)
	}

	return fixture, total, page, perPage, nil
}

func getSingleFixture(ID string) (*Fixture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"_id", objID}}}},
		{{"$lookup", bson.D{
			{"from", "competitions"},
			{"localField", "competition_id"},
			{"foreignField", "_id"},
			{"as", "competition_id"},
		}}},
		{{"$lookup", bson.D{
			{"from", "teams"},
			{"localField", "home_team_id"},
			{"foreignField", "_id"},
			{"as", "home_team_id"},
		}}},
		{{"$lookup", bson.D{
			{"from", "teams"},
			{"localField", "away_team_id"},
			{"foreignField", "_id"},
			{"as", "away_team_id"},
		}}},
		{{"$unwind", "$competition_id"}},
		{{"$unwind", "$home_team_id"}},
		{{"$unwind", "$away_team_id"}},
	}

	cursor, err := fixtureCollection.Aggregate(ctx, pipeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team with the id %v is not found", ID)
		}
		return nil, fmt.Errorf("failed to execute aggregation pipeline: %v", err)
	}

	var fixture Fixture
	if cursor.Next(ctx) {
		if err := cursor.Decode(&fixture); err != nil {
			return nil, fmt.Errorf("failed to decode team: %v", err)
		}
	} else {
		// If no documents are returned by the pipeline
		return nil, fmt.Errorf("fixture with the id %v is not found", ID)
	}

	return &fixture, nil
}

func getSingleFixtureByHash(hash string) (*Fixture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"unique_link", hash}}}},
		{{"$lookup", bson.D{
			{"from", "competitions"},
			{"localField", "competition_id"},
			{"foreignField", "_id"},
			{"as", "competition_id"},
		}}},
		{{"$unwind", "$competition_id"}},
		{{"$lookup", bson.D{
			{"from", "teams"},
			{"localField", "home_team_id"},
			{"foreignField", "_id"},
			{"as", "home_team_id"},
		}}},
		{{"$unwind", "$home_team_id"}},
		{{"$lookup", bson.D{
			{"from", "teams"},
			{"localField", "away_team_id"},
			{"foreignField", "_id"},
			{"as", "away_team_id"},
		}}},
		{{"$unwind", "$away_team_id"}},
	}

	cursor, err := fixtureCollection.Aggregate(ctx, pipeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("fixture with the hash  %v is not found", hash)
		}
		return nil, fmt.Errorf("failed to execute aggregation pipeline: %v", err)
	}

	var fixture Fixture
	if cursor.Next(ctx) {
		if err := cursor.Decode(&fixture); err != nil {
			return nil, fmt.Errorf("failed to decode team: %v", err)
		}
	} else {
		// If no documents are returned by the pipeline
		return nil, fmt.Errorf("fixture with the hash  %v is not found", hash)
	}

	return &fixture, nil
}

func updateFixture(ID string, update UpdateFixture) (*models.Fixture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Parse ObjectID from string
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Create update fields
	updates := bson.M{}
	if update.CompetitionID != primitive.NilObjectID {
		updates["competition_id"] = update.CompetitionID
	}
	if update.HomeTeamID != primitive.NilObjectID {
		updates["home_team_id"] = update.HomeTeamID
	}
	if update.AwayTeamID != primitive.NilObjectID {
		updates["away_team_id"] = update.AwayTeamID
	}
	if !update.Date.IsZero() {
		updates["date"] = update.Date
	}
	if update.Status != "" {
		updates["status"] = update.Status
	}
	if update.UniqueLink != "" {
		updates["unique_link"] = update.UniqueLink
	}
	if update.Referee != "" {
		updates["referee"] = update.Referee
	}

	// Add fields that are always updated
	updates["updated_at"] = time.Now()

	// Perform the update operation
	_, err = fixtureCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		return nil, fmt.Errorf("could not update link: %v", err)
	}

	var fixture models.Fixture
	err = fixtureCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&fixture)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated fixture: %v", err)
	}
	return &fixture, nil
}

func updateFixtureStats(ID string, update UpdateFixtureStats) (*models.Fixture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Parse ObjectID from string
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	updates := bson.M{
		"updated_at": time.Now(),
	}

	if update.Home.Goals != 0 {
		updates["home.goals"] = update.Home.Goals
	}

	if update.Home.ShotsOnTarget != 0 {
		updates["home.shots_on_target"] = update.Home.ShotsOnTarget
	}

	if update.Home.Possession != 0 {
		updates["home.possession"] = update.Home.Possession
	}

	if update.Home.Passes != 0 {
		updates["home.passes"] = update.Home.Passes
	}

	if update.Home.PassesAccuracy != 0 {
		updates["home.passes_accuracy"] = update.Home.PassesAccuracy
	}

	if update.Home.Fouls != 0 {
		updates["home.fouls"] = update.Home.Fouls
	}

	if update.Home.YellowCards != 0 {
		updates["home.yellow_cards"] = update.Home.YellowCards
	}

	if update.Home.RedCards != 0 {
		updates["home.red_cards"] = update.Home.RedCards
	}

	if update.Home.OffSides != 0 {
		updates["home.off_sides"] = update.Home.OffSides
	}

	if update.Home.Corners != 0 {
		updates["home.corners"] = update.Home.Corners
	}

	if update.Home.Formation != "" {
		updates["home.formation"] = update.Home.Formation
	}

	if len(update.Home.GoalScorers) > 0 {
		updates["home.goal_scorers"] = update.Home.GoalScorers
	}

	if len(update.Home.Substitutes) > 0 {
		updates["home.substitutes"] = update.Home.Substitutes
	}

	if len(update.Home.Lineup) > 0 {
		updates["home.lineup"] = update.Home.Lineup
	}

	if update.Away.Goals != 0 {
		updates["away.goals"] = update.Away.Goals
	}

	if update.Away.ShotsOnTarget != 0 {
		updates["away.shots_on_target"] = update.Away.ShotsOnTarget
	}

	if update.Away.Possession != 0 {
		updates["away.possession"] = update.Away.Possession
	}

	if update.Away.Passes != 0 {
		updates["away.passes"] = update.Away.Passes
	}

	if update.Away.PassesAccuracy != 0 {
		updates["away.passes_accuracy"] = update.Away.PassesAccuracy
	}

	if update.Away.Fouls != 0 {
		updates["away.fouls"] = update.Away.Fouls
	}

	if update.Away.YellowCards != 0 {
		updates["away.yellow_cards"] = update.Away.YellowCards
	}

	if update.Away.RedCards != 0 {
		updates["away.red_cards"] = update.Away.RedCards
	}

	if update.Away.OffSides != 0 {
		updates["away.off_sides"] = update.Away.OffSides
	}

	if update.Away.Corners != 0 {
		updates["away.corners"] = update.Away.Corners
	}

	if update.Away.Formation != "" {
		updates["away.formation"] = update.Away.Formation
	}

	if len(update.Away.GoalScorers) > 0 {
		updates["away.goal_scorers"] = update.Away.GoalScorers
	}

	if len(update.Away.Substitutes) > 0 {
		updates["away.substitutes"] = update.Away.Substitutes
	}

	if len(update.Away.Lineup) > 0 {
		updates["away.lineup"] = update.Away.Lineup
	}

	// Perform the update operation
	_, err = fixtureCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		return nil, fmt.Errorf("could not update link: %v", err)
	}

	var fixture models.Fixture
	err = fixtureCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&fixture)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated fixture: %v", err)
	}
	return &fixture, nil
}

func deleteFixture(ID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %v", err)
	}

	result, err := fixtureCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return fmt.Errorf("failed to delete team: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no team found with ID %s", ID)
	}

	return nil
}

func getCompetitions(filters CompetitionRequest, pageNumber string, pageSize string) ([]models.Competition, int64, int64, int64, error) {
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
	if filters.Type != "" {
		filter["type"] = filters.Type
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	fOpt := options.FindOptions{Limit: &perPage, Skip: &offset, Sort: bson.D{{"created_at", -1}}}
	cOpt := options.CountOptions{Limit: &perPage, Skip: &offset}

	total, err := competitionCollection.CountDocuments(ctx, filter, &cOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to count competitions: %v", err)
	}

	cursor, err := competitionCollection.Find(ctx, filter, &fOpt)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to find competitions: %v", err)
	}
	defer cursor.Close(ctx)

	var competitions []models.Competition
	if err := cursor.All(ctx, &competitions); err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to decode competitions: %v", err)
	}

	return competitions, total, page, perPage, nil
}

func getSingleCompetition(ID string) (*models.Competition, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	var competition models.Competition
	err = competitionCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&competition)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch competition: %v", err)
	}

	return &competition, nil
}
