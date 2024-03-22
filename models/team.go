package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" `
	Name        string             `bson:"name" validate:"required" json:"name"`
	Country     string             `bson:"country" validate:"required" json:"country"`
	FoundedYear int                `bson:"founded_year" validate:"required" json:"founded_year"`
	Stadium     string             `bson:"stadium" validate:"required" json:"stadium"`
	Sponsor     string             `bson:"sponsor" validate:"required" json:"sponsor"`
	Trophies    []Trophy           `bson:"trophies" json:"trophies"`
	Players     []Player           `bson:"players" json:"players"`
}

type Trophy struct {
	Name         string `bson:"name" validate:"required" json:"name"`
	YearObtained int    `bson:"year_obtained" validate:"required" json:"year_obtained"`
}
