package users

import (
	"github.com/gin-gonic/gin"

	"league/helpers"
	"league/models"
	
	"net/http"
)

func getUserHandler(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched user",
		StatusCode: http.StatusOK,
		Data:       user,
	})
}

func getUsersHandler(ctx *gin.Context) {
	query := UserRequest{
		Email:     ctx.Query("email"),
		FirstName: ctx.Query("first_name"),
		LastName:  ctx.Query("last_name"),
	}

	users, total, page, perPage, err := getUsers(query, ctx.Query("page"), ctx.Query("per_page"))

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched users",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     users,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func updateUserHandler(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	var req UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	id := user.Id.String()
	updatedUser, err := updateUser(id, req)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated user",
		StatusCode: http.StatusOK,
		Data:       updatedUser,
	})
}

func deleteUserHandler(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	err = deleteUser(user.Id.String())
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted user",
		StatusCode: http.StatusOK,
		Data:       nil,
	})

}
