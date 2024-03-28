package teams

import (
	"context"
	"fmt"
	"testing"
	"time"

	"league/models"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBURI    = ""
	testDatabase = "real"
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
	teamCollection = client.Database(testDatabase).Collection("teams")
}

func cleanupTestEnvironment(t *testing.T) {
	// Drop the test collection after tests are completed
	// err := teamCollection.Drop(context.Background())
	// if err != nil {
	// 	t.Fatalf("Error dropping test collection: %v", err)
	// }

	// Disconnect from MongoDB
	err := teamCollection.Database().Client().Disconnect(context.Background())
	if err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

func TestCreateTeam_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Create a team object
	team := models.Team{
		Name:        "manchester",
		State:       "hudders",
		Country:     "englan",
		FoundedYear: 200,
		Stadium:     "hudders",
		Sponsor:     "three",
		CreatedBy:   primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Call the createTeam function
	resp, err := createTeam(team)

	// Assert that the function returns no error and the insertedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, team.Name, resp.Name)
	assert.Equal(t, team.State, resp.State)
	assert.Equal(t, team.Country, resp.Country)
	assert.Equal(t, team.FoundedYear, resp.FoundedYear)
	assert.Equal(t, team.Stadium, resp.Stadium)
	assert.Equal(t, team.Sponsor, resp.Sponsor)
	assert.Equal(t, team.CreatedBy, resp.CreatedBy)
}

func TestCreateTeam_DuplicateError(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	team := models.Team{
		Name:        "chelsea",
		State:       "cobham",
		Country:     "englang",
		FoundedYear: 200,
		Stadium:     "stamford",
		Sponsor:     "three",
		CreatedBy:   primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Attempt to create a new team with the same name and stadium

	insertedTeam, err := createTeam(team)

	// Assert that the function returns an error indicating duplicate key
	assert.Error(t, err)
	assert.Nil(t, insertedTeam)
	assert.EqualError(t, err, fmt.Sprintf("name  %s or stadium %s already exists", team.Name, team.Stadium))
}

func TestGetSingleTeam_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Call the getSingleTeam function
	teamWithCreator, err := getSingleTeam("66058baf3166ffd82cd5ff46")

	assert.NoError(t, err)
	assert.NotNil(t, teamWithCreator)
}

func TestGetSingleTeam_NotFound(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	id := primitive.NewObjectID().Hex()
	// Call the getSingleTeam function with a non-existing team ID
	teamWithCreator, err := getSingleTeam(id)

	// Assert that the function returns an error indicating team not found
	assert.Error(t, err)
	assert.Nil(t, teamWithCreator)
	assert.EqualError(t, err, "team with the id "+id+" is not found")
}

func TestDeleteTeam_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)


	// Call the deleteTeam function
	err := deleteTeam("66058baf3166ffd82cd5ff46")

	// Assert that the function returns no error
	assert.NoError(t, err)
}

func TestUpdateTeam_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Update fields
	update := TeamRequest{
		Name:        "al nassar",
		State:       "emirates",
		Country:     "suadi arabia",
		FoundedYear: 2020,
		Stadium:     "al-awwal park",
		Sponsor:     "Lebara",
	}

	// Call the updateUser function
	updatedTeam, err := updateUser("66058baf3166ffd82cd5ff46", update)

	// Assert that the function returns no error and the updatedTeam is not nil
	assert.NoError(t, err)
	assert.NotNil(t, updatedTeam)

	// Assert that the fields of the updated team match the updated values
	assert.Equal(t, update.Name, updatedTeam.Name)
	assert.Equal(t, update.State, updatedTeam.State)
	assert.Equal(t, update.Country, updatedTeam.Country)
	assert.Equal(t, update.FoundedYear, updatedTeam.FoundedYear)
	assert.Equal(t, update.Stadium, updatedTeam.Stadium)
	assert.Equal(t, update.Sponsor, updatedTeam.Sponsor)
}

func TestGetTeam_Success(t *testing.T) {
	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	
	// Define filters
	filters := TeamQueryRequest{
		Query: "suadi arabia",
	}

	// Call the getTeam function
	teams, _, _, _, err := getTeam(filters, "1", "10")

	// Assert that the function returns no error and the teams slice is not empty
	assert.NoError(t, err)
	assert.NotEmpty(t, teams)
	assert.Greater(t, len(teams), 0)
}
