package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Role string

const (
	SuperAdminRole Role = "super-admin"
	AdminRole      Role = "admin"
	UserRole       Role = "user"
)

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	FirstName         string             `bson:"first_name,omitempty" validate:"required" json:"first_name"`
	LastName          string             `bson:"last_name,omitempty" validate:"required" json:"last_name"`
	Email             string             `bson:"email" validate:"required" json:"email"`
	RoleName          Role               `bson:"role" validate:"required" json:"role"`
	VerificationToken string             `bson:"verification_token" json:"verification_token"`
	ExpiresAt         time.Time          `bson:"expires_at" json:"expires_at"`
	Password          string             `bson:"password" json:"-"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

func (u *User) SetEmail() {
	u.Email = strings.ToLower(u.Email)
}

func GetUserFromContext(ctx *gin.Context) (*User, error) {
	value, exists := ctx.Get("user")
	if !exists {
		return nil, errors.New("user not found in session")
	}
	user, ok := value.(User)
	if !ok {
		return nil, errors.New("mismatching types")
	}
	return &user, nil
}
