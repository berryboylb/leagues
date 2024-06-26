package teams

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"league/helpers"
	"league/models"
)

func createHandler(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	var req TeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	newTeam := models.Team{
		Name:        req.Name,
		Country:     req.Country,
		State:       req.State,
		FoundedYear: req.FoundedYear,
		Stadium:     req.Stadium,
		Sponsor:     req.Sponsor,
		CreatedBy:   user.Id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := createTeam(newTeam)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created team",
		StatusCode: http.StatusOK,
		Data:       result,
	})
}

func getHandler(ctx *gin.Context) {
	var (
		foundedYear int
		err         error
	)

	if year := ctx.Query("founded_year"); year != "" {
		if foundedYear, err = strconv.Atoi(year); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	query := TeamQueryRequest{
		Name:        ctx.Query("name"),
		Country:     ctx.Query("country"),
		State:       ctx.Query("state"),
		FoundedYear: foundedYear,
		Stadium:     ctx.Query("stadium"),
		Sponsor:     ctx.Query("sponsor"),
		Query:      strings.ToLower(ctx.Query("query")),
	}

	teams, total, page, perPage, err := getTeam(query, ctx.Query("page"), ctx.Query("per_page"))

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched teams",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     teams,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})

}

func getSingleHandler(ctx *gin.Context) {
	team, err := getSingleTeam(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched team",
		StatusCode: http.StatusOK,
		Data:       team,
	})
}

func updateHandler(ctx *gin.Context) {
	var req TeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	updatedTeam, err := updateUser(ctx.Param("id"), req)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated team",
		StatusCode: http.StatusOK,
		Data:       updatedTeam,
	})
}

func deleteHandler(ctx *gin.Context) {
	err := deleteTeam(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted team",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func getPlayersHandler(ctx *gin.Context) {
	fmt.Println("hello")
	var (
		teamID primitive.ObjectID
		err    error
	)

	if team := ctx.Query("team_id"); team != "" {
		if teamID, err = primitive.ObjectIDFromHex(team); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	query := PlayerRequest{
		Name:        ctx.Query("name"),
		Position:  ctx.Query("position"),
		TeamID: teamID,
	}

	players, total, page, perPage, err := getPlayers(query, ctx.Query("page"), ctx.Query("per_page"))

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched players",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     players,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getPlayerHandler(ctx *gin.Context) {
	player, err := getSinglePlayer(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched player",
		StatusCode: http.StatusOK,
		Data:       player,
	})
}
