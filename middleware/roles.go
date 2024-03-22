package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"league/helpers"
	"league/models"
)

func RolesMiddleware(roles []models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := models.GetUserFromContext(c)
		if err != nil {
			helpers.CreateResponse(c, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}
		
		// Check if user role matches any of the allowed roles
		for _, allowedRole := range roles {
			if allowedRole == user.RoleName {
				c.Next()
				return
			}
		}

		helpers.CreateResponse(c, helpers.Response{
			Message:    "You don't have the required role for this resource",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
	}
}
