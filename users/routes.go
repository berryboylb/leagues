package users

import (
	"github.com/gin-gonic/gin"

	"league/jwt"
	// "league/middleware"
	// "league/models"
)

// var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
// var everybody []models.RoleAllowed = []models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func UserRoutes(superRoute *gin.RouterGroup) {
	userRouter := superRoute.Group("/users")
	{
		userRouter.Use(jwt.Middleware())
		// userRouter.GET("/", middleware.RolesMiddleware([]models.RoleAllowed{models.AdminRole, models.SuperAdminRole}), GetAllUsers)
		// userRouter.POST("/", middleware.RolesMiddleware([]models.RoleAllowed{models.AdminRole, models.SuperAdminRole}), CreateAdmin)
		// userRouter.GET("/user", User)
		// userRouter.PATCH("/user", UpdateUser)
		// userRouter.PATCH("/user/:id", middleware.RolesMiddleware([]models.RoleAllowed{models.AdminRole, models.SuperAdminRole}), ReinStateAccount)
		// userRouter.DELETE("/user", DeleteUser)
	}
}