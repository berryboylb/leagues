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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" validate:"required" json:"name"`
	Type      string             `bson:"type" validate:"required" json:"type"` // e.g., "League", "World Cup", "Champions League", etc.
	Country   string             `bson:"country" validate:"required" json:"country"`
	StartDate time.Time          `bson:"start_date" validate:"required" json:"start_date"`
	EndDate   time.Time          `bson:"end_date" validate:"required" json:"end_date"`
}

type Fixture struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	CompetitionID   primitive.ObjectID `bson:"competition_id" validate:"required"  json:"competition_id"`
	HomeTeamID      primitive.ObjectID `bson:"home_team_id" validate:"required" json:"home_team_id"`
	AwayTeamID      primitive.ObjectID `bson:"away_team_id" validate:"required" json:"away_team_id"`
	Date            time.Time          `bson:"date" validate:"required" json:"date"`
	Result          string             `bson:"result" validate:"required" json:"result"`
	Status          Status             `bson:"status"  json:"status"` // Completed, Pending, etc.
	UniqueLink      string             `bson:"unique_link" validate:"required" json:"unique_link"`
	Location        string             `bson:"location" validate:"required" json:"location"`
	Stadium         string             `bson:"stadium" json:"stadium"`
	Referee         string             `bson:"referee" json:"referee"`
	Attendance      int                `bson:"attendance" json:"attendance"`
	TicketsSold     int                `bson:"tickets_sold" json:"tickets_sold"`
	HomeFormation   string             `bson:"home_formation" validate:"required" json:"home_formation"`
	AwayFormation   string             `bson:"away_formation" validate:"required" json:"away_formation"`
	HomeLineup      []Player           `bson:"home_lineup" json:"home_lineup"`
	AwayLineup      []Player           `bson:"away_lineup" json:"away_lineup"`
	HomeSubstitutes []Player           `bson:"home_substitutes" json:"home_substitutes"`
	AwaySubstitutes []Player           `bson:"away_substitutes" json:"away_substitutes"`
	Stats           Stats              `bson:"stats" json:"stats"`
	HomeGoalScorers []string           `bson:"home_goal_scorers" json:"home_goal_scorers"`
	AwayGoalScorers []string           `bson:"away_goal_scorers" json:"away_goal_scorers"`
}

type Player struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" `
	Name   string             `bson:"name" json:"name"`
	Image  string             `bson:"img" json:"img"`
	TeamID primitive.ObjectID `bson:"team_id" json:"team_id"`
	Status PlayerStatus       `bson:"status" json:"status"`
}

type Stats struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	HomeTeamStats StatDetails        `bson:"home_team_stats" json:"home_team_stats"`
	AwayTeamStats StatDetails        `bson:"away_team_stats" json:"away_team_stats"`
}

type StatDetails struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" `
	Goals          int                `bson:"goals" json:"goals"`
	Shots          int                `bson:"shots" json:"shots"`
	ShotsOnTarget  int                `bson:"shots_on_target" json:"shots_on_target"`
	Possession     float64            `bson:"possession" json:"possession"`
	Passes         int                `bson:"passes" json:"passes"`
	PassesAccuracy int                `bson:"passes_accuracy" json:"passes_accuracy"`
	Fouls          int                `bson:"fouls" json:"fouls"`
	YellowCards    int                `bson:"yellow_cards" json:"yellow_cards"`
	RedCards       int                `bson:"red_cards" json:"red_cards"`
	OffSides       int                `bson:"off_sides" json:"off_sides"`
	Corners        int                `bson:"corners" json:"corners"`
}
