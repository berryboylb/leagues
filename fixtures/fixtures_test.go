package fixtures

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"league/models"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBURI    = ""
	testDatabase = "real"
)

func TestGenerateRandomString(t *testing.T) {
	rand.Seed(42) // Seed for deterministic testing

	tests := []struct {
		length int
	}{
		{10},
		{20},
		{30},
	}

	for _, test := range tests {
		randomString, err := generateRandomString(test.length)
		if err != nil {
			t.Errorf("Error generating random string: %v", err)
		}

		expectedLength := test.length
		if len(randomString) != expectedLength {
			t.Errorf("Generated string length mismatch. Expected: %d, Got: %d", expectedLength, len(randomString))
		}

	}
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
	fixtureCollection = client.Database(testDatabase).Collection("fixtures")
}

func cleanupTestEnvironment(t *testing.T) {
	// Drop the test collection after tests are completed
	// err := fixtureCollection.Drop(context.Background())
	// if err != nil {
	// 	t.Fatalf("Error dropping test collection: %v", err)
	// }

	// Disconnect from MongoDB
	err := fixtureCollection.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

func setupTestCompEnvironment(t *testing.T) {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(testDBURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Set up the database and collection
	competitionCollection = client.Database(testDatabase).Collection("competitions")
}

func cleanupTestCompEnvironment(t *testing.T) {
	// Drop the test collection after tests are completed
	// err := fixtureCollection.Drop(context.Background())
	// if err != nil {
	// 	t.Fatalf("Error dropping test collection: %v", err)
	// }

	// Disconnect from MongoDB
	err := competitionCollection.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

func TestCreate_Competition(t *testing.T) {
	setupTestCompEnvironment(t)
	defer cleanupTestCompEnvironment(t)
	comp := models.Competition{
		Name: "uefa",
		Type: "domestic",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := competitionCollection.InsertOne(context.Background(), comp)
	assert.NoError(t, err)
}

func TestCreateFixture(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	competitionID, err := primitive.ObjectIDFromHex("6606af2f8ea9f277021e23ea") //gotten from db
	assert.NoError(t, err)

	team1ID, err := primitive.ObjectIDFromHex("660595c06c25f01f95f72670") //gotten from db
	assert.NoError(t, err)

	team2ID, err := primitive.ObjectIDFromHex("6605964e0c3b6abc49e55641") //gotten from db
	assert.NoError(t, err)

	randomString, err := generateRandomString(50)
	assert.NoError(t, err)

	subs := []string{"john", "doe"}
	mains := []string{"john", "doe"}
	forms := "4-3-3"

	var fixture = models.Fixture{
		HomeTeamID:    team1ID,
		AwayTeamID:    team2ID,
		CompetitionID: competitionID,
		Status:        models.Pending,
		Date:          time.Now(),
		Stadium:       "emirates",
		Referee:       "john snow",
		UniqueLink:    randomString,
		Away: models.Details{
			Substitutes: subs,
			Lineup:      mains,
			Formation:   forms,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Home: models.Details{
			Substitutes: subs,
			Lineup:      mains,
			Formation:   forms,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Call the createFixture function
	dbfixture, err := createFixture(fixture)
	assert.NoError(t, err)
	assert.NotNil(t, dbfixture)

	assert.Equal(t, randomString, dbfixture.UniqueLink)
	assert.Equal(t, competitionID, dbfixture.CompetitionID)
	assert.Equal(t, team1ID, dbfixture.HomeTeamID)
	assert.Equal(t, team2ID, dbfixture.AwayTeamID)
	assert.Equal(t, competitionID, dbfixture.CompetitionID)
	assert.Equal(t, models.Pending, dbfixture.Status)
	assert.Equal(t, "emirates", dbfixture.Stadium)
	assert.Equal(t, "john snow", dbfixture.Referee)

	assert.Equal(t, subs, dbfixture.Home.Substitutes)
	assert.Equal(t, mains, dbfixture.Home.Lineup)
	assert.Equal(t, forms, dbfixture.Home.Formation)
	assert.Equal(t, 0, dbfixture.Home.Goals)
	assert.Equal(t, 0, dbfixture.Home.Fouls)
	assert.Equal(t, 0, dbfixture.Home.Corners)
	assert.Equal(t, 0.0, dbfixture.Home.Possession)
	assert.Equal(t, 0, dbfixture.Home.PassesAccuracy)
	assert.Equal(t, 0, dbfixture.Home.YellowCards)
	assert.Equal(t, 0, dbfixture.Home.RedCards)
	assert.Equal(t, 0, dbfixture.Home.Shots)
	assert.Equal(t, 0, dbfixture.Home.ShotsOnTarget)

	assert.Equal(t, subs, dbfixture.Away.Substitutes)
	assert.Equal(t, mains, dbfixture.Away.Lineup)
	assert.Equal(t, forms, dbfixture.Away.Formation)
	assert.Equal(t, 0, dbfixture.Away.Goals)
	assert.Equal(t, 0, dbfixture.Away.Fouls)
	assert.Equal(t, 0, dbfixture.Away.Corners)
	assert.Equal(t, 0.0, dbfixture.Away.Possession)
	assert.Equal(t, 0, dbfixture.Away.PassesAccuracy)
	assert.Equal(t, 0, dbfixture.Away.YellowCards)
	assert.Equal(t, 0, dbfixture.Away.RedCards)
	assert.Equal(t, 0, dbfixture.Away.Shots)
	assert.Equal(t, 0, dbfixture.Away.ShotsOnTarget)
}

func TestGetSingleFixture_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	// Call the getSingleTeam function
	resp, err := getSingleFixture("6606b1acda826498e1205a47") // gotten from the db

	// Assert that the function returns no error and the teamWithCreator is not nil
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetSingleFixture_Failure(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	id := primitive.NewObjectID().Hex()
	// Call the getSingleTeam function
	resp, err := getSingleFixture(id)

	// Assert that the function returns no error and the teamWithCreator is not nil
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "fixture with the id "+id+" is not found")
}

func TestDeleteFixture_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the deleteTeam function
	err := deleteFixture("6606b1acda826498e1205a47")

	// Assert that the function returns no error
	assert.NoError(t, err)
}

func TestUpdateFixture_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	unique, err := generateRandomString(10)
	assert.NoError(t, err)

	competitionID, err := primitive.ObjectIDFromHex("6606af2f8ea9f277021e23ea")
	assert.NoError(t, err)

	team1ID, err := primitive.ObjectIDFromHex("660595c06c25f01f95f72670")
	assert.NoError(t, err)

	team2ID, err := primitive.ObjectIDFromHex("6605964e0c3b6abc49e55641")
	assert.NoError(t, err)

	// Update fields
	update := UpdateFixture{
		HomeTeamID:    team1ID,
		AwayTeamID:    team2ID,
		CompetitionID: competitionID,
		Stadium:       "emirates",
		Status:        models.Completed,
		UniqueLink:    unique,
		Referee:       "jon snow",
	}

	// Call the updatefixture function
	resp, err := updateFixture("6606b1acda826498e1205a47", update)

	// Assert that the function returns no error and the updatedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Assert that the fields of the updated team match the updated values
	assert.Equal(t, update.AwayTeamID, resp.AwayTeamID)
	assert.Equal(t, update.HomeTeamID, resp.HomeTeamID)
	assert.Equal(t, update.CompetitionID, resp.CompetitionID)
	assert.Equal(t, update.Stadium, resp.Stadium)
	assert.Equal(t, update.UniqueLink, resp.UniqueLink)
	assert.Equal(t, update.Referee, resp.Referee)
}

func TestUpdateFixtureStats_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Update fields
	update := UpdateFixtureStats{
		Home: Stats{
			Goals:           1,
			GoalScorers:     []string{"Neymar (47)"},
			Shots:           15,
			ShotsOnTarget:   5,
			Possession:      53.5,
			Passes:          873,
			Fouls:           3,
			YellowCards:     1,
			RedCards:        0,
			Corners:         15,
		},
		Away: Stats{
			Goals:           3,
			GoalScorers:     []string{"c ronaldo (47)", "bale (58)", "benzema (87)"},
			Shots:           10,
			ShotsOnTarget:   7,
			Possession:      46.5,
			Passes:          973,
			Fouls:           5,
			YellowCards:     3,
			RedCards:        0,
			Corners:         10,
		},
	}

	// Call the updatefixture function
	resp, err := updateFixtureStats("6606b1acda826498e1205a47", update)

	// Assert that the function returns no error and the updatedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Assert that the fields of the updated stats
	assert.Equal(t, update.Home.Goals, resp.Home.Goals)
	assert.Equal(t, update.Home.GoalScorers, resp.Home.GoalScorers)
	assert.Equal(t, update.Home.Shots, resp.Home.Shots)
	assert.Equal(t, update.Home.ShotsOnTarget, resp.Home.ShotsOnTarget)
	assert.Equal(t, update.Home.Possession, resp.Home.Possession)
	assert.Equal(t, update.Home.Passes, resp.Home.Passes)
	assert.Equal(t, update.Home.Fouls, resp.Home.Fouls)
	assert.Equal(t, update.Home.YellowCards, resp.Home.YellowCards)
	assert.Equal(t, update.Home.RedCards, resp.Home.RedCards)
	assert.Equal(t, update.Home.Corners, resp.Home.Corners)

	assert.Equal(t, update.Away.Goals, resp.Away.Goals)
	assert.Equal(t, update.Away.GoalScorers, resp.Away.GoalScorers)
	assert.Equal(t, update.Away.Shots, resp.Away.Shots)
	assert.Equal(t, update.Away.ShotsOnTarget, resp.Away.ShotsOnTarget)
	assert.Equal(t, update.Away.Possession, resp.Away.Possession)
	assert.Equal(t, update.Away.Passes, resp.Away.Passes)
	assert.Equal(t, update.Away.Fouls, resp.Away.Fouls)
	assert.Equal(t, update.Away.YellowCards, resp.Away.YellowCards)
	assert.Equal(t, update.Away.RedCards, resp.Away.RedCards)
	assert.Equal(t, update.Away.Corners, resp.Away.Corners)
}


func TestGetFixtures_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Define search features request
	query := SearchFeaturesRequest{
		Query: "completed", // made sure my fixture status is completed
	}

	// Call getFixtures function
	fixtures, total, page, perPage, err := getFixtures(query, "1", "15")

	// Assert no error
	assert.NoError(t, err)

	// Assert page, perPage and total
	assert.Equal(t, int64(1), page)
	assert.Equal(t, int64(15), perPage)
	assert.Equal(t, int64(1), total)

	// Assert returned fixtures
	assert.NotNil(t, fixtures)
	assert.Greater(t, len(fixtures), 0)
}
