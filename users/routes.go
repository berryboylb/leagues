package users

import (
	"github.com/gin-gonic/gin"

	"league/jwt"
	"league/middleware"
	"league/models"
)


func UserRoutes(superRoute *gin.RouterGroup) {
	userRouter := superRoute.Group("/users")
	{
		userRouter.Use(jwt.Middleware())
		userRouter.GET("/", middleware.RolesMiddleware([]models.Role{models.AdminRole, models.SuperAdminRole}), getUsersHandler)
		userRouter.GET("/:id", getUserHandler)
		userRouter.PATCH("/:id", updateUserHandler)
		userRouter.DELETE("/:id", deleteUserHandler)
	}
}
