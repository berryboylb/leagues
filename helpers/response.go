package helpers

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

func CreateResponse(ctx *gin.Context, response Response) {
	ctx.JSON(response.StatusCode, gin.H{
		"message":    response.Message,
		"data":       response.Data,
		"statusCode": response.StatusCode,
	})
	ctx.Abort()
}
