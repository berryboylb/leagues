package fixtures

import (
	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	"league/helpers"
	"league/models"

	"net/http"
	"strconv"
	"time"
	"fmt"
)

func createFixtureHandler(ctx *gin.Context) {
	var req CreateTestFixture
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	randomString, err := generateRandomString(10)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	if req.HomeTeamID == req.AwayTeamID {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "home and away teams cannot be same",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	newFixture := models.Fixture{
		HomeTeamID:    req.HomeTeamID,
		AwayTeamID:    req.AwayTeamID,
		CompetitionID: req.CompetitionID,
		Status:        models.Status(req.Status),
		Date:          req.Date,
		Stadium:       req.Stadium,
		Referee:       req.Referee,
		UniqueLink:    randomString,
		Away: models.Details{
			Substitutes: req.Away.Substitutes,
			Lineup:      req.Away.Lineup,
			Formation:   req.Away.Formation,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Home: models.Details{
			Substitutes: req.Home.Substitutes,
			Lineup:      req.Home.Lineup,
			Formation:   req.Home.Formation,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := createFixture(newFixture)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created fixture",
		StatusCode: http.StatusOK,
		Data:       result,
	})
}

func viewFixturesByTypeHandler(ctx *gin.Context) {
	resp, total, page, perPage, err := getFixturesByStatus(ctx.Param("status"), ctx.Query("page"), ctx.Query("per_page"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched response",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     resp,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func updateFixtureHandler(ctx *gin.Context) {
	var req UpdateFixture
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	resp, err := updateFixture(ctx.Param("id"), req)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	if req.HomeTeamID == req.AwayTeamID {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    "home and away teams cannot be same",
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated fixture",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func updateFixtureStatsHandler(ctx *gin.Context) {
	var req UpdateFixtureStats
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	resp, err := updateFixtureStats(ctx.Param("id"), req)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated fixture stats",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func deleteFixtureHandler(ctx *gin.Context) {
	err := deleteFixture(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted fixture",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func generateUniqueHash(ctx *gin.Context) {
	length := 10
	if ctx.Query("length") != "" {
		if len, err := strconv.Atoi(ctx.Query("length")); err == nil {
			length = len
		}
	}
	randomString, err := generateRandomString(length)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully generated hash",
		StatusCode: http.StatusOK,
		Data:       randomString,
	})
}

func getFixtureByHash(ctx *gin.Context) {
	hash := ctx.Param("link")
	fmt.Println(hash)
	resp, err := getSingleFixtureByHash(hash)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched fixture",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

var (
	from      time.Time
	to        time.Time
	searchErr error
)

func searchHandler(ctx *gin.Context) {
	if date := ctx.Query("from"); date != "" {
		if from, searchErr = time.Parse("2006-01-02T15:04:05Z", date); searchErr != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    searchErr.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if date := ctx.Query("to"); date != "" {
		if to, searchErr = time.Parse("2006-01-02T15:04:05Z", date); searchErr != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    searchErr.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}
	query := SearchFeaturesRequest{
		Query: ctx.Query("query"),
		From:  from,
		To:    to,
	}

	resp, total, page, perPage, err := getFixtures(query, ctx.Query("page"), ctx.Query("per_page"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched response",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     resp,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func singleFixtureHandler(ctx *gin.Context) {
	resp, err := getSingleFixture(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched fixture",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func getCompetitionsHandler(ctx *gin.Context) {
	query := CompetitionRequest{
		Name: ctx.Query("name"),
		Type: ctx.Query("type"),
	}

	players, total, page, perPage, err := getCompetitions(query, ctx.Query("page"), ctx.Query("per_page"))

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched competitions",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     players,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleCompetitionsHandler(ctx *gin.Context) {
	competition, err := getSingleCompetition(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched  competition",
		StatusCode: http.StatusOK,
		Data:       competition,
	})
}
