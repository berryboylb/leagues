package main

import (
	"github.com/gin-gonic/gin"

	"league/auth"
	"league/users"
	"league/fixtures"
	"league/e-teams"
)

func AddRoutes(superRoute *gin.RouterGroup) {
	//register routes
	auth.AuthRoutes(superRoute)
	users.UserRoutes(superRoute)
	fixtures.FixtureRoutes(superRoute)
	teams.TeamRoutes(superRoute)
}
