package fixtures

import (
	"github.com/gin-gonic/gin"

	"league/jwt"
	"league/middleware"
	"league/models"
)

var admins []models.Role = []models.Role{models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func FixtureRoutes(superRoute *gin.RouterGroup) *gin.RouterGroup {
	fixtureRouter := superRoute.Group("/fixtures")
	{
		//public
		fixtureRouter.GET("/", searchHandler)

		//protected
		fixtureRouter.Use(jwt.Middleware())
		fixtureRouter.POST("/", middleware.RolesMiddleware(admins), createFixtureHandler)
		fixtureRouter.GET("/:status", viewFixturesByTypeHandler)
		fixtureRouter.GET("/:link", getFixtureByHash)
		fixtureRouter.GET("/:id", singleFixtureHandler)
		fixtureRouter.PATCH("/:id", middleware.RolesMiddleware(admins), updateFixtureHandler)
		fixtureRouter.DELETE("/:id", middleware.RolesMiddleware(admins), deleteFixtureHandler)
		fixtureRouter.GET("/competitions", getCompetitionsHandler)
		fixtureRouter.GET("/competitions/:id", getSingleCompetitionsHandler)
	}
	return fixtureRouter
}
