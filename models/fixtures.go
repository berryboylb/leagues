package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status string

const (
	Ongoing   Status = "ongoing"
	Pending   Status = "pending"
	Completed Status = "completed"
)

type PlayerStatus string

const (
	Active  PlayerStatus = "active"
	Injured PlayerStatus = "injured"
)

type Competition struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `bson:"name" validate:"required" json:"name"`
	Type      string             `bson:"type" validate:"required" json:"type"` // e.g., "League", "World Cup", "Champions League", etc.
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Fixture struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CompetitionID primitive.ObjectID `bson:"competition_id" validate:"required" json:"competition_id"`
	HomeTeamID    primitive.ObjectID `bson:"home_team_id" validate:"required" json:"home_team_id"`
	AwayTeamID    primitive.ObjectID `bson:"away_team_id" validate:"required" json:"away_team_id"`
	Home          Details            `bson:"home" json:"home"`
	Away          Details            `bson:"away" json:"away"`
	Date          time.Time          `bson:"date" validate:"required" json:"date"`
	Status        Status             `bson:"status"  json:"status"` // Completed, Pending, etc.
	UniqueLink    string             `bson:"unique_link" validate:"required" json:"unique_link"`
	Stadium       string             `bson:"stadium" json:"stadium"`
	Referee       string             `bson:"referee" json:"referee"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type Player struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `bson:"name" json:"name"`
	Image     string             `bson:"img" json:"img"`
	Position  string             `bson:"position" json:"position"`
	TeamID    primitive.ObjectID `bson:"team_id" json:"team_id"`
	Status    PlayerStatus       `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Details struct {
	Goals          int       `bson:"goals" json:"goals"`
	GoalScorers    []string  `bson:"goal_scorers" json:"goal_scorers"`
	Substitutes    []string  `bson:"substitutes" json:"substitutes"`
	Lineup         []string  `bson:"lineup" json:"lineup"`
	Formation      string    `bson:"formation" validate:"required" json:"formation"`
	Shots          int       `bson:"shots" json:"shots"`
	ShotsOnTarget  int       `bson:"shots_on_target" json:"shots_on_target"`
	Possession     float64   `bson:"possession" json:"possession"`
	Passes         int       `bson:"passes" json:"passes"`
	PassesAccuracy int       `bson:"passes_accuracy" json:"passes_accuracy"`
	Fouls          int       `bson:"fouls" json:"fouls"`
	YellowCards    int       `bson:"yellow_cards" json:"yellow_cards"`
	RedCards       int       `bson:"red_cards" json:"red_cards"`
	OffSides       int       `bson:"off_sides" json:"off_sides"`
	Corners        int       `bson:"corners" json:"corners"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}
