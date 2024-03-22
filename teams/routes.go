package teams

import (
	"github.com/gin-gonic/gin"

	"league/jwt"
	// "fixtures/middleware"
	// "fixtures/models"
)

// var everybody []models.RoleAllowed = []models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func TeamRoutes(superRoute *gin.RouterGroup) {
	//only admins have access to tems
	teamRouter := superRoute.Group("/teams")
	{
		teamRouter.Use(jwt.Middleware())
		// teamRouter.Use(roles middleware)
		teamRouter.GET("/", getHandler)
		teamRouter.GET("/:id", getSingleHandler)
		teamRouter.POST("/", createHandler)
		teamRouter.PATCH("/user", updateHandler)
		teamRouter.DELETE("/user", deleteHandler)
		//    router.HandleFunc("/teams", ViewTeams).Methods("GET")
	}
}
