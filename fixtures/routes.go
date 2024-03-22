package fixtures

import (
	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func FixtureRoutes(superRoute *gin.RouterGroup) {
	fixtureRouter := superRoute.Group("/fixtures")
	{
		fixtureRouter.GET("/latest", getLatestFixturesHandler)
		fixtureRouter.GET("/completed", viewCompletedFixturesHandler)
		fixtureRouter.GET("/pending", viewPendingFixturesHandler)
		fixtureRouter.GET("/search", searchHandler)
		fixtureRouter.GET("/:id", singleFixtureHandler)
		// the remaining routes for admin
	}
}
