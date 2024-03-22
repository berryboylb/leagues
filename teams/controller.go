package teams

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	query := TeamRequest{
		Name:        ctx.Query("name"),
		Country:     ctx.Query("country"),
		State:       ctx.Query("state"),
		FoundedYear: foundedYear,
		Stadium:     ctx.Query("stadium"),
		Sponsor:     ctx.Query("sponsor"),
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
