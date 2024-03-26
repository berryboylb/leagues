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
	fixtureCollection = client.Database(testDatabase).Collection("fixtures")
}

func cleanupTestEnvironment(t *testing.T) {
	// Drop the test collection after tests are completed
	err := fixtureCollection.Drop(context.Background())
	if err != nil {
		t.Fatalf("Error dropping test collection: %v", err)
	}

	// Disconnect from MongoDB
	err = fixtureCollection.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

var fixture = models.Fixture{
	HomeTeamID:    primitive.NewObjectID(),
	AwayTeamID:    primitive.NewObjectID(),
	CompetitionID: primitive.NewObjectID(),
	Status:        models.Pending,
	Date:          time.Now(),
	Stadium:       "emirates",
	Referee:       "john snow",
	UniqueLink:    "gggg",
	Away: models.Details{
		Substitutes: []string{"john", "doe"},
		Lineup:      []string{"john", "doe"},
		Formation:   "4-3-3",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	Home: models.Details{
		Substitutes: []string{"john", "doe"},
		Lineup:      []string{"john", "doe"},
		Formation:   "4-3-3",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestCreateFixture(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the createFixture function
	insertedTeam, err := createFixture(fixture)

	// Assert that the function returns no error and the insertedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, insertedTeam)
}

func TestGetSingleFixture_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the createFixture function
	insertedTeam, err := createFixture(fixture)

	// Assert that the function returns no error and the insertedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, insertedTeam)

	// Call the getSingleTeam function
	resp, err := getSingleFixture(insertedTeam.ID.Hex())

	// Assert that the function returns no error and the teamWithCreator is not nil
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetSingleFixture_Failure(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the createFixture function
	insertedTeam, err := createFixture(fixture)

	// Assert that the function returns no error and the insertedFixture is not nil
	assert.NoError(t, err)
	assert.NotNil(t, insertedTeam)

	// Call the getSingleTeam function
	resp, err := getSingleFixture(primitive.NewObjectID().Hex())

	// Assert that the function returns no error and the teamWithCreator is not nil
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "fixture with the id "+primitive.NewObjectID().Hex()+" is not found")
}

func TestDeleteFixture_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the createFixture function
	insertedTeam, err := createFixture(fixture)

	// Assert that the function returns no error and the insertedFixture is not nil
	assert.NoError(t, err)
	assert.NotNil(t, insertedTeam)

	// Call the deleteTeam function
	err = deleteFixture(insertedTeam.ID.Hex())

	// Assert that the function returns no error
	assert.NoError(t, err)
}

func TestUpdateFixture_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the createFixture function
	insertedTeam, err := createFixture(fixture)

	// Assert that the function returns no error and the insertedFixture is not nil
	assert.NoError(t, err)
	assert.NotNil(t, insertedTeam)

	unique, err := generateRandomString(10)
	assert.NoError(t, err)

	// Update fields
	update := UpdateFixture{
		HomeTeamID:    primitive.NewObjectID(),
		AwayTeamID:    primitive.NewObjectID(),
		CompetitionID: primitive.NewObjectID(),
		Stadium:       "emirates",
		Status:        models.Completed,
		UniqueLink:    unique,
		Referee:       "jon snow",
	}

	// Call the updatefixture function
	resp, err := updateFixture(insertedTeam.ID.Hex(), update)

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

	// Call the createFixture function
	insertedTeam, err := createFixture(fixture)

	// Assert that the function returns no error and the insertedFixture is not nil
	assert.NoError(t, err)
	assert.NotNil(t, insertedTeam)

	// Update fields
	update := UpdateFixtureStats{
		Home: Stats{
			Goals: 1, 
			GoalScorers: []string{"Neymar (47)"}, 
			Shots: 15,
			ShotsOnTarget: 5,
			Possession: 53.5,
			Passes: 873,
			Fouls: 3,
			YellowCards: 1,
			RedCards: 0,
			Corners: 15,
		},
		Away: Stats{
			Goals: 3, 
			GoalScorers: []string{"c ronaldo (47)", "bale (58)", "benzema (87)"}, 
			Shots: 10,
			ShotsOnTarget: 7,
			Possession: 46.5,
			Passes: 973,
			Fouls: 5,
			YellowCards: 3,
			RedCards: 0,
			Corners: 10,
		},
	}

	// Call the updatefixture function
	resp, err := updateFixtureStats(insertedTeam.ID.Hex(), update)

	// Assert that the function returns no error and the updatedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Assert that the fields of the updated stats
	assert.Equal(t, update.Home, resp.Home)
	assert.Equal(t, update.Away, resp.Away)
}


func TestGetFixtures_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)



	// Define search features request
	query := SearchFeaturesRequest{
		Query: "completed",
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
