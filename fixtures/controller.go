package fixtures

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"league/helpers"
	"league/models"

	"net/http"
	"time"
)

func createFixtureHandler(ctx *gin.Context) {
	var req CreateFixture
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	homeID, err := req.ParseHex(req.HomeTeamID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	awayID, err := req.ParseHex(req.AwayTeamID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	compID, err := req.ParseHex(req.CompetitionID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	status, err := parseStatus(req.Status)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	date, err := req.ParseDate()
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	newFixture := models.Fixture{
		HomeTeamID:    homeID,
		AwayTeamID:    awayID,
		CompetitionID: compID,
		Status:        status,
		Date:          date,
		Stadium:       req.Stadium,
		Referee:       req.Referee,
		UniqueLink:    "to be done",
		Away: models.Details{
			Substitutes: req.Away.Substitutes,
			Lineup:      req.Away.Lineup,
			Formation:   req.Away.Formation,
		},
		Home: models.Details{
			Substitutes: req.Home.Substitutes,
			Lineup:      req.Home.Lineup,
			Formation:   req.Home.Formation,
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
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	homeID, err := req.ParseHex(req.HomeTeamID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	awayID, err := req.ParseHex(req.AwayTeamID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	compID, err := req.ParseHex(req.CompetitionID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	status, err := parseStatus(req.Status)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	date, err := req.ParseDate()
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	update := UpdateFixture{
		HomeTeamID:    homeID,
		AwayTeamID:    awayID,
		CompetitionID: compID,
		Status:        status,
		Date:          date,
		Stadium:       req.Stadium,
		Referee:       req.Referee,
		UniqueLink:    "to be done",
	}

	resp, err := updateFixture(ctx.Param("id"), update)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
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

func updateFixtureStatus(ctx *gin.Context) {
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

}

func getFixtureByHash(ctx *gin.Context) {
	resp, err := getSingleFixtureByHash(ctx.Param("link"))
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

func searchHandler(ctx *gin.Context) {
	var (
		compId primitive.ObjectID
		homeid primitive.ObjectID
		awayid primitive.ObjectID
		status models.Status
		err    error
	)

	if id := ctx.Query("competition_id"); id != "" {
		if compId, err = primitive.ObjectIDFromHex(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("home_team_id"); id != "" {
		if homeid, err = primitive.ObjectIDFromHex(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("away_team_id"); id != "" {
		if awayid, err = primitive.ObjectIDFromHex(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("status"); id != "" {
		if status, err = parseStatus(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	query := SearchFeaturesRequest{
		Competition: compId,
		HomeTeam:    homeid,
		AwayTeam:    awayid,
		Status:      status,
		UniqueLink:  ctx.Query("unique_link"),
		Referee:     ctx.Query("referee"),
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
