package main

import (
	"github.com/gin-gonic/gin"

	"league/auth"
	"league/users"
	"league/fixtures"
)

func AddRoutes(superRoute *gin.RouterGroup) {
	//register routes
	auth.AuthRoutes(superRoute)
	users.UserRoutes(superRoute)
	fixtures.FixtureRoutes(superRoute)
}
