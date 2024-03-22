package auth

import (
	"github.com/gin-gonic/gin"

	"league/jwt"
	// "league/middleware"
	// "league/models"
)

// var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
// var everybody []models.RoleAllowed = []models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func AuthRoutes(superRoute *gin.RouterGroup) {
	authRouter := superRoute.Group("/auth")
	{
		authRouter.Use(jwt.Middleware())
	}
}