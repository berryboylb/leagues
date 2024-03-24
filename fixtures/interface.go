package fixtures

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"

	"league/models"
)

type CompetitionRequest struct {
	Name string `json:"name" binding:"required,min=3"`
	Type string `json:"type" binding:"required,min=3"`
}

type SearchFeaturesRequest struct {
	Competition primitive.ObjectID `json:"competition"`
	UniqueLink  string             `json:"link"`
	HomeTeam    primitive.ObjectID `json:"home_team"`
	AwayTeam    primitive.ObjectID `json:"away_team"`
	Query       string             `json:"query"`
	Referee     string             `json:"referee"`
	Status      models.Status      `json:"status"`
	From        time.Time          `json:"from"`
	To          time.Time          `json:"to"`
}

type CreateStats struct {
	Substitutes []string ` json:"substitutes" binding:"required,len=5"`
	Lineup      []string ` json:"lineup" binding:"required,len=11"`
	Formation   string   `  json:"formation" binding:"required"`
}

type CreateTestFixture struct {
	CompetitionID primitive.ObjectID `json:"competition_id" binding:"required"`
	HomeTeamID    primitive.ObjectID `json:"home_team_id" binding:"required"`
	AwayTeamID    primitive.ObjectID `json:"away_team_id" binding:"required"`
	Date          time.Time          `json:"date" binding:"required" time_format:"2006-01-02"`
	Status        models.Status      `json:"status" binding:"required,oneof=completed ongoing pending"`
	Stadium       string             `json:"stadium" binding:"required"`
	Referee       string             `json:"referee" binding:"required"`
	Home          CreateStats        `json:"home" binding:"required"`
	Away          CreateStats        `json:"away" binding:"required"`
}

type CreateFixture struct {
	CompetitionID string      `json:"competition_id" binding:"required"`
	HomeTeamID    string      `json:"home_team_id" binding:"required"`
	AwayTeamID    string      `json:"away_team_id" binding:"required"`
	Date          string      `json:"date" binding:"required"`
	Status        string      `json:"status" binding:"required,parseStatus"`
	UniqueLink    string      `json:"unique_link" binding:"required"`
	Stadium       string      `json:"stadium" binding:"required"`
	Referee       string      `json:"referee" binding:"required"`
	Home          CreateStats `json:"home" binding:"required"`
	Away          CreateStats `json:"away" binding:"required"`
}

func (c *CreateFixture) ParseDate() (time.Time, error) {
	return time.Parse("2006-01-02", c.Date)
}

func (c *CreateFixture) ParseHex(ID string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(ID)
}

func parseStatus(status string) (models.Status, error) {
	switch status {
	case "completed":
		return models.Completed, nil
	case "ongoing":
		return models.Ongoing, nil
	case "pending":
		return models.Pending, nil
	default:
		return "", fmt.Errorf("%v is not a valid status type", status)
	}
}

// this is to update a fixture
type UpdateRequest struct {
	CompetitionID string `json:"competition_id" binding:"required"`
	HomeTeamID    string `json:"home_team_id" binding:"required"`
	AwayTeamID    string `json:"away_team_id" binding:"required"`
	Date          string `json:"date" binding:"required"`
	Status        string `json:"status" binding:"required"`
	UniqueLink    string `json:"unique_link" binding:"required"`
	Stadium       string `json:"stadium" binding:"required"`
	Referee       string `json:"referee" binding:"required"`
}

func (u *UpdateRequest) ParseDate() (time.Time, error) {
	return time.Parse("2006-01-02", u.Date)
}

func (u *UpdateRequest) ParseHex(ID string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(ID)
}

type UpdateFixture struct {
	CompetitionID primitive.ObjectID `json:"competition_id" binding:"omitempty"`
	HomeTeamID    primitive.ObjectID `json:"home_team_id" binding:"omitempty"`
	AwayTeamID    primitive.ObjectID `json:"away_team_id" binding:"omitempty"`
	Date          time.Time          `json:"date" binding:"omitempty" time_format:"2006-01-02"`
	Status        models.Status      `json:"status" binding:"omitempty,oneof=completed ongoing pending"`
	UniqueLink    string             `json:"unique_link" binding:"omitempty"`
	Stadium       string             `json:"stadium" binding:"omitempty"`
	Referee       string             `json:"referee" binding:"omitempty"`
}

type UpdateFixtureStats struct {
	Home Stats ` json:"home" binding:"required"`
	Away Stats ` json:"away" binding:"required"`
}

type Stats struct {
	Goals          int      ` json:"goals" binding:"omitempty"`
	GoalScorers    []string `json:"goal_scorers" binding:"omitempty"`
	Substitutes    []string ` json:"substitutes" binding:"omitempty,len=5"`
	Lineup         []string ` json:"lineup" binding:"omitempty,len=11"`
	Formation      string   `  json:"formation" binding:"omitempty"`
	Shots          int      ` json:"shots" binding:"omitempty"`
	ShotsOnTarget  int      `json:"shots_on_target" binding:"omitempty"`
	Possession     float64  ` json:"possession" binding:"omitempty"`
	Passes         int      ` json:"passes" binding:"omitempty"`
	PassesAccuracy int      ` json:"passes_accuracy" binding:"omitempty"`
	Fouls          int      ` json:"fouls" binding:"omitempty"`
	YellowCards    int      ` json:"yellow_cards" binding:"omitempty"`
	RedCards       int      ` json:"red_cards" binding:"omitempty"`
	OffSides       int      ` json:"off_sides" binding:"omitempty"`
	Corners        int      ` json:"corners" binding:"omitempty"`
}

// this is for aggregated fixture
type Fixture struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	CompetitionID models.Competition `bson:"competition_id" json:"competition_id"`
	HomeTeamID    models.Team        `bson:"home_team_id" validate:"required" json:"home_team_id"`
	AwayTeamID    models.Team        `bson:"away_team_id" validate:"required" json:"away_team_id"`
	Home          models.Details     `bson:"home" json:"home"`
	Away          models.Details     `bson:"away" json:"away"`
	Date          time.Time          `bson:"date" validate:"required" json:"date"`
	Status        models.Status      `bson:"status"  json:"status"` // Completed, Pending, etc.
	UniqueLink    string             `bson:"unique_link" validate:"required" json:"unique_link"`
	Stadium       string             `bson:"stadium" json:"stadium"`
	Referee       string             `bson:"referee" json:"referee"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
