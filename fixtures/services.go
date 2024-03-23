package fixtures

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

func getFixtures(query SearchFeaturesRequest, pageNumber string, pageSize string) ([]models.Fixture, int64, int64, int64, error) {
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

	if query.Competition != primitive.NilObjectID {
		filter["competition_id"] = query.Competition
	}

	if query.HomeTeam != primitive.NilObjectID {
		filter["home_team_id"] = query.Competition
	}

	if query.AwayTeam != primitive.NilObjectID {
		filter["away_team_id"] = query.AwayTeam
	}

	if query.UniqueLink != "" {
		filter["unique_link"] = query.UniqueLink
	}

	if query.Referee != "" {
		filter["referee"] = query.Referee
	}

	if query.Status != "" {
		filter["status"] = query.Status
	}

	if !query.From.IsZero() {
		filter["created_at"] = bson.M{"$gte": query.From}
	}

	if !query.To.IsZero() && !query.From.IsZero() {
		filter["created_at"] = bson.M{"$gte": query.From, "$lte": query.To}
	}

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
		return nil, fmt.Errorf("team with the id %v is not found", ID)
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
			return nil, fmt.Errorf("fixture with the hash %v is not found", hash)
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
		return nil, fmt.Errorf("fixture with the hash %v is not found", hash)
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
	updates := bson.M{
		"home_team_id": update.HomeTeamID,
		"away_team_id": update.AwayTeamID,
		"date":         update.Date,
		"status":       update.Stadium,
		"unique_link":  update.UniqueLink,
		"referee":      update.Referee,
		"updated_at":   time.Now(),
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

func updateFixtureStats(ID string, update UpdateFixtureStats) (*models.Fixture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Parse ObjectID from string
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	// Create update fields
	updates := bson.M{
		"home": bson.M{
			"goals":           update.Home.Goals,
			"goal_scorers":    update.Home.GoalScorers,
			"substitutes":     update.Home.Substitutes,
			"lineup":          update.Home.Lineup,
			"formation":       update.Home.Formation,
			"shots_on_target": update.Home.ShotsOnTarget,
			"possession":      update.Home.Possession,
			"passes":          update.Home.Passes,
			"passes_accuracy": update.Home.PassesAccuracy,
			"fouls":           update.Home.Fouls,
			"yellow_cards":    update.Home.YellowCards,
			"red_cards":       update.Home.RedCards,
			"off_sides":       update.Home.OffSides,
			"corners":         update.Home.Corners,
		},
		"away": bson.M{
			"goals":           update.Away.Goals,
			"goal_scorers":    update.Away.GoalScorers,
			"substitutes":     update.Away.Substitutes,
			"lineup":          update.Away.Lineup,
			"formation":       update.Away.Formation,
			"shots_on_target": update.Away.ShotsOnTarget,
			"possession":      update.Away.Possession,
			"passes":          update.Away.Passes,
			"passes_accuracy": update.Away.PassesAccuracy,
			"fouls":           update.Away.Fouls,
			"yellow_cards":    update.Away.YellowCards,
			"red_cards":       update.Away.RedCards,
			"off_sides":       update.Away.OffSides,
			"corners":         update.Away.Corners,
		},
		"updated_at": time.Now(),
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
