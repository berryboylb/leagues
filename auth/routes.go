package auth

import (
	"github.com/gin-gonic/gin"
)

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
	}
}
