package teams

import (
	"time"

	"league/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Country     string `json:"country" binding:"required,min=3"`
	State       string `json:"state" binding:"required,min=3"`
	FoundedYear int    `json:"founded_year" binding:"required"`
	Stadium     string `json:"stadium" binding:"required,min=3"`
	Sponsor     string `json:"sponsor" binding:"required,min=3"`
}

type TeamQueryRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Country     string `json:"country" binding:"required,min=3"`
	State       string `json:"state" binding:"required,min=3"`
	FoundedYear int    `json:"founded_year" binding:"required"`
	Stadium     string `json:"stadium" binding:"required,min=3"`
	Sponsor     string `json:"sponsor" binding:"required,min=3"`
	Query string `json:"query" binding:"required,min=3"`
}

type TeamFilterRequest struct{}

type TeamWithCreator struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id" `
	Name        string             `bson:"name" validate:"required" json:"name"`
	State       string             `bson:"state" validate:"required" json:"state"`
	Country     string             `bson:"country" validate:"required" json:"country"`
	FoundedYear int                `bson:"founded_year" validate:"required" json:"founded_year"`
	Stadium     string             `bson:"stadium" validate:"required" json:"stadium"`
	Sponsor     string             `bson:"sponsor" validate:"required" json:"sponsor"`
	CreatedBy   models.User        `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type PlayerRequest struct {
	Name     string             `json:"name" binding:"required,min=3"`
	Position string             `json:"position" binding:"required,min=3"`
	TeamID   primitive.ObjectID `json:"team_id" binding:"required,min=3"`
}
