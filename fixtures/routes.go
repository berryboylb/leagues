package fixtures

import (
	"github.com/gin-gonic/gin"

	"league/jwt" //remove during unit tests
	"league/middleware"
	"league/models"
)

var admins []models.Role = []models.Role{models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func FixtureRoutes(superRoute *gin.RouterGroup) {
	fixtureRouter := superRoute.Group("/fixtures")
	{
		//public
		fixtureRouter.GET("/", searchHandler)

		//protected
		fixtureRouter.Use(jwt.Middleware())
		fixtureRouter.POST("/", middleware.RolesMiddleware(admins), createFixtureHandler)
		fixtureRouter.POST("/hash", middleware.RolesMiddleware(admins), generateUniqueHash)
		fixtureRouter.GET("/status/:status", viewFixturesByTypeHandler)
		fixtureRouter.GET("/:link", getFixtureByHash)
		fixtureRouter.GET("/fixture/:id", singleFixtureHandler)
		fixtureRouter.PATCH("/:id", middleware.RolesMiddleware(admins), updateFixtureHandler)
		fixtureRouter.PATCH("/stats/:id", middleware.RolesMiddleware(admins), updateFixtureStatsHandler)
		fixtureRouter.DELETE("/:id", middleware.RolesMiddleware(admins), deleteFixtureHandler)
		fixtureRouter.GET("/competitions", getCompetitionsHandler)
		fixtureRouter.GET("/competitions/:id", getSingleCompetitionsHandler)
	}
}
