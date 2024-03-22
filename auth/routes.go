package auth

import (
	"github.com/gin-gonic/gin"
	// "league/jwt"
	// "league/middleware"
	// "league/models"
)

// var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
// var everybody []models.RoleAllowed = []models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func AuthRoutes(superRoute *gin.RouterGroup) {
	authRouter := superRoute.Group("/auth")
	{
		authRouter.POST("/user/signup", signUpUserHandler)
		authRouter.POST("/admin/signup", signUpAdminHandler)
		authRouter.POST("/user/login", loginUserHandler)
		authRouter.POST("/admin/login", loginAdminHandler)
		authRouter.POST("/admin/login/confirm", confirmLoginAdminHandler)
		authRouter.POST("/forgot-password", forgotPasswordHandler)
		authRouter.POST("/reset-password", resetPasswordHandler)
		// authRouter.Use(jwt.Middleware())
	}
}
