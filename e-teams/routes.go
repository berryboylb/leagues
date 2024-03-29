package teams

import (
	"github.com/gin-gonic/gin"

	"league/jwt" //remove during test
	"league/middleware"
	"league/models"
)

var admins []models.Role = []models.Role{models.AdminRole, models.SuperAdminRole}

func TeamRoutes(superRoute *gin.RouterGroup) {
	teamRouter := superRoute.Group("/teams")
	{
		teamRouter.GET("/", getHandler)

		teamRouter.Use(jwt.Middleware())
		teamRouter.POST("/", middleware.RolesMiddleware(admins), createHandler)
		teamRouter.GET("/:id", getSingleHandler)
		teamRouter.PATCH("/:id", middleware.RolesMiddleware(admins), updateHandler)
		teamRouter.DELETE("/:id", middleware.RolesMiddleware(admins), deleteHandler)
		teamRouter.GET("/players", getPlayersHandler)
		teamRouter.GET("/players/:id", getPlayerHandler)
	}
}
