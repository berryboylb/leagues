package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id" `
	Name        string             `bson:"name" validate:"required" json:"name"`
	State       string             `bson:"state" validate:"required" json:"state"`
	Country     string             `bson:"country" validate:"required" json:"country"`
	FoundedYear int                `bson:"founded_year" validate:"required" json:"founded_year"`
	Stadium     string             `bson:"stadium" validate:"required" json:"stadium"`
	Sponsor     string             `bson:"sponsor" validate:"required" json:"sponsor"`
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// type Trophy struct {
// 	Name      string    `bson:"name" validate:"required" json:"name"`
// 	Image     string    `bson:"image" validate:"required" json:"image"`
// 	CreatedAt time.Time `bson:"created_at" json:"created_at"`
// 	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
// }

// type TeamTrophy struct {
// 	TrophyID     primitive.ObjectID `bson:"trophy_id" json:"trophy_id"`
// 	TeamID       primitive.ObjectID `bson:"team_id" json:"team_id"`
// 	YearObtained int                `bson:"year_obtained" validate:"required" json:"year_obtained"`
// 	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
// 	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
// }
