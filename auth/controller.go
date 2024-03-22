package auth

import (
	// "errors"
	// "fmt"
	// "log"
	"net/http"
	// "net/url"
	// "os"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"

	"league/helpers"
	"league/jwt"
	"league/models"
	// "league/models"
	// "league/notifications"
)

func signUpUserHandler(ctx *gin.Context) {
	var req SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	hash, err := helpers.HashPassword(req.Password, 8)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	newUser := models.User{
		FirstName: req.FirstName,
		LastName:  req.FirstName,
		Email:     req.Email,
		RoleName:  models.UserRole,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := createUser(newUser)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created user",
		StatusCode: http.StatusOK,
		Data:       result,
	})
}

func loginUserHandler(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	user, err := getUserByEmail(req.Email)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	if user.RoleName != models.UserRole {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "you don't have permission to access this resource",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}

	isMatch := helpers.CheckPasswordHash(req.Password, user.Password)
	if !isMatch {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Invalid credentials",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	token, err := jwt.GenerateJWT(user.Id.String())
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully signed in",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
}

func signUpAdminHandler(ctx *gin.Context) {
	var req SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	hash, err := helpers.HashPassword(req.Password, 8)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	newUser := models.User{
		FirstName: req.FirstName,
		LastName:  req.FirstName,
		Email:     req.Email,
		RoleName:  models.AdminRole,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := createUser(newUser)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created admin",
		StatusCode: http.StatusOK,
		Data:       result,
	})
}

func loginAdminHandler(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	user, err := getUserByEmail(req.Email)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	if user.RoleName != models.AdminRole && user.RoleName != models.SuperAdminRole {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "you don't have permission to access this resource",
			StatusCode: http.StatusUnauthorized,
			Data:       nil,
		})
		return
	}

	isMatch := helpers.CheckPasswordHash(req.Password, user.Password)
	if !isMatch {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "Invalid password",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	// send the otp in a seperate thread
	go sendOtp(user, "Please confirm login")

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully sent otp",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func confirmLoginAdminHandler(ctx *gin.Context) {
	var req OtpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	user, err := getUserFromOtp(req.Otp, req.Email)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	token, err := jwt.GenerateJWT(user.Id.String())
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	user.VerificationToken = ""
	user.ExpiresAt = time.Time{}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully logged In admin",
		StatusCode: http.StatusOK,
		Data:       token,
	})

	go destroyToken(user.Id)
}

func forgotPasswordHandler(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	user, err := getUserByEmail(req.Email)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	// send the otp in a seperate thread
	go sendOtp(user, "Reset Password")

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully sent mail",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func resetPasswordHandler(ctx *gin.Context) {
	var req Otp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	user, err := getUser(req.Otp)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = changePassword(user.Id, req.Password)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	go destroyToken(user.Id)

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully changed password",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}
