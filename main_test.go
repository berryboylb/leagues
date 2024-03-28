// main_test.go

package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"
	"io"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "league/fixtures"
	// "league/models"
	// "log"
	"net/http"
	"testing"
	"time"
)

var jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MTE1OTI2NDIsImlkIjoiNjVmZDY3YjQ4MDczMTQ2NzY2MGY4NGQ3In0.IPUBBPSzBmTRjCicnW5m4fur8U8lgA41M47jF2Vgyvw"

type LoginUser struct {
	Data struct {
		Token string `json:"token"`
		User  struct {
			ID                string    `json:"_id"`
			FirstName         string    `json:"first_name"`
			LastName          string    `json:"last_name"`
			Email             string    `json:"email"`
			Role              string    `json:"role"`
			VerificationToken string    `json:"verification_token"`
			ExpiresAt         time.Time `json:"expires_at"`
		} `json:"user"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestLogin(t *testing.T) {
	// login creds
	payload := map[string]string{
		"email":    "johndoe4586@gmail.com",
		"password": "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/user/login", "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response LoginUser
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully signed in", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, payload["email"], response.Data.User.Email)

}

type UserResponse struct {
	Data       interface{} ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

type SignUpPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignUpResponse struct {
	Data struct {
		ID                string    `json:"_id"`
		FirstName         string    `json:"first_name"`
		LastName          string    `json:"last_name"`
		Email             string    `json:"email"`
		Role              string    `json:"role"`
		VerificationToken string    `json:"verification_token"`
		ExpiresAt         time.Time `json:"expires_at"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestSignup(t *testing.T) {
	// signup creds
	payload := SignUpPayload{
		FirstName: "john",
		LastName:  "doe",
		Email:     "johndoe4586@gmail.com",
		Password:  "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/user/signup", "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response SignUpResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully created user", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, payload.FirstName, response.Data.FirstName)
	assert.Equal(t, payload.LastName, response.Data.LastName)
	assert.Equal(t, payload.Email, response.Data.Email)
	assert.Equal(t, "user", response.Data.Role)
}

func TestSignupAdmin(t *testing.T) {
	// signup creds
	payload := SignUpPayload{
		FirstName: "john",
		LastName:  "doe",
		Email:     "admin0178@gmail.com",
		Password:  "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/admin/signup", "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	defer resp.Body.Close()
	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response SignUpResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully created admin", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, payload.FirstName, response.Data.FirstName)
	assert.Equal(t, payload.LastName, response.Data.LastName)
	assert.Equal(t, payload.Email, response.Data.Email)
	assert.Equal(t, "admin", response.Data.Role)
}

type LoginAdmin struct {
	Data       interface{} ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

func TestLoginAdmin(t *testing.T) {
	// signup creds
	payload := map[string]string{
		"email":    "admin0178@gmail.com",
		"password": "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/admin/login", "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response LoginAdmin
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully sent otp", response.Message)

}

type ConfirmLoginAdmin struct {
	Data       string ` json:"data"`
	Message    string ` json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestConfirmLoginAdmin(t *testing.T) {
	payload := map[string]string{
		"email": "admin0178@gmail.com",
		"otp":   "9328",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/admin/login/confirm", "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response ConfirmLoginAdmin
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully logged In admin", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type Team struct {
	Data struct {
		ID          string    `json:"_id"`
		Name        string    `json:"name"`
		State       string    `json:"state"`
		Country     string    `json:"country"`
		FoundedYear int       `json:"founded_year"`
		Stadium     string    `json:"stadium"`
		Sponsor     string    `json:"sponsor"`
		CreatedBy   string    `json:"created_by"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type GetTeam struct {
	Data struct {
		ID          string      `json:"_id"`
		Name        string      `json:"name"`
		State       string      `json:"state"`
		Country     string      `json:"country"`
		FoundedYear int         `json:"founded_year"`
		Stadium     string      `json:"stadium"`
		Sponsor     string      `json:"sponsor"`
		CreatedBy   interface{} `json:"created_by"`
		CreatedAt   time.Time   `json:"created_at"`
		UpdatedAt   time.Time   `json:"updated_at"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type TeamPayload struct {
	Name        string `json:"name"`
	Country     string `json:"country"`
	State       string `json:"state"`
	FoundedYear int    `json:"founded_year"`
	Stadium     string `json:"stadium"`
	Sponsor     string `json:"sponsor"`
}

func TestAddTeam(t *testing.T) {
	payload := TeamPayload{
		Name:        "testing",
		Country:     "england",
		State:       "london",
		FoundedYear: 2001,
		Stadium:     "stamford34",
		Sponsor:     "adidas",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)
	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://league-ybry.onrender.com/api/v1/teams", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response Team
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully created team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	//parse the id returned from the response
	_, err = primitive.ObjectIDFromHex(response.Data.ID)
	assert.NoError(t, err)

	assert.Equal(t, payload.Name, response.Data.Name)
	assert.Equal(t, payload.Country, response.Data.Country)
	assert.Equal(t, payload.State, response.Data.State)
	assert.Equal(t, payload.FoundedYear, response.Data.FoundedYear)
	assert.Equal(t, payload.Stadium, response.Data.Stadium)
	assert.Equal(t, payload.Sponsor, response.Data.Sponsor)
	fmt.Print(response.Data)
}

func TestEditTeam(t *testing.T) {

	// Create a new request with the necessary payload
	payload := TeamPayload{
		Name:        "teste2e",
		Country:     "england",
		State:       "london",
		FoundedYear: 2001,
		Stadium:     "stamforde2e",
		Sponsor:     "adidas",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)
	// Create a new HTTP request
	req, err := http.NewRequest("PATCH", "https://league-ybry.onrender.com/api/v1/teams/6604cdfc6912f7e29fee1ec2", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response Team
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully updated team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, payload.Name, response.Data.Name)
	assert.Equal(t, payload.Country, response.Data.Country)
	assert.Equal(t, payload.State, response.Data.State)
	assert.Equal(t, payload.FoundedYear, response.Data.FoundedYear)
	assert.Equal(t, payload.Stadium, response.Data.Stadium)
	assert.Equal(t, payload.Sponsor, response.Data.Sponsor)
	fmt.Print(response.Data)
}

func TestViewTeam(t *testing.T) {
	req, err := http.NewRequest("GET", "https://league-ybry.onrender.com/api/v1/teams/6604cdfc6912f7e29fee1ec2", nil)
	assert.NoError(t, err)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response GetTeam
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	_, err = primitive.ObjectIDFromHex(response.Data.ID)
	assert.NoError(t, err)
}

func TestRemoveTeam(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "https://league-ybry.onrender.com/api/v1/teams/6604cdfc6912f7e29fee1ec2", nil)
	assert.NoError(t, err)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response Team
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully deleted team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type SearchTeam struct {
	Data struct {
		Data []struct {
			ID          string    `json:"_id"`
			Name        string    `json:"name"`
			State       string    `json:"state"`
			Country     string    `json:"country"`
			FoundedYear int       `json:"founded_year"`
			Stadium     string    `json:"stadium"`
			Sponsor     string    `json:"sponsor"`
			CreatedBy   string    `json:"created_by"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"data"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Total   int `json:"total"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestSearchTeam(t *testing.T) {
	resp, err := http.Get("https://league-ybry.onrender.com/api/v1/teams?page=1&per_page=15")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response SearchTeam
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully fetched teams", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Greater(t, len(response.Data.Data), 0, "Array length should be greater than zero")
}

type Fixture struct {
	Data       interface{} ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

type AggFixture struct {
	Data struct {
		ID            string `json:"ID"`
		CompetitionID struct {
			ID        string    `json:"_id"`
			Name      string    `json:"name"`
			Type      string    `json:"type"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"competition_id"`
		HomeTeamID struct {
			ID          string    `json:"_id"`
			Name        string    `json:"name"`
			State       string    `json:"state"`
			Country     string    `json:"country"`
			FoundedYear int       `json:"founded_year"`
			Stadium     string    `json:"stadium"`
			Sponsor     string    `json:"sponsor"`
			CreatedBy   string    `json:"created_by"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"home_team_id"`
		AwayTeamID struct {
			ID          string    `json:"_id"`
			Name        string    `json:"name"`
			State       string    `json:"state"`
			Country     string    `json:"country"`
			FoundedYear int       `json:"founded_year"`
			Stadium     string    `json:"stadium"`
			Sponsor     string    `json:"sponsor"`
			CreatedBy   string    `json:"created_by"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"away_team_id"`
		Home struct {
			Goals          int         `json:"goals"`
			GoalScorers    interface{} `json:"goal_scorers"`
			Substitutes    []string    `json:"substitutes"`
			Lineup         []string    `json:"lineup"`
			Formation      string      `json:"formation"`
			Shots          int         `json:"shots"`
			ShotsOnTarget  int         `json:"shots_on_target"`
			Possession     float64     `json:"possession"`
			Passes         int         `json:"passes"`
			PassesAccuracy int         `json:"passes_accuracy"`
			Fouls          int         `json:"fouls"`
			YellowCards    int         `json:"yellow_cards"`
			RedCards       int         `json:"red_cards"`
			OffSides       int         `json:"off_sides"`
			Corners        int         `json:"corners"`
			CreatedAt      time.Time   `json:"created_at"`
			UpdatedAt      time.Time   `json:"updated_at"`
		} `json:"home"`
		Away struct {
			Goals          int         `json:"goals"`
			GoalScorers    interface{} `json:"goal_scorers"`
			Substitutes    []string    `json:"substitutes"`
			Lineup         []string    `json:"lineup"`
			Formation      string      `json:"formation"`
			Shots          int         `json:"shots"`
			ShotsOnTarget  int         `json:"shots_on_target"`
			Possession     float64     `json:"possession"`
			Passes         int         `json:"passes"`
			PassesAccuracy int         `json:"passes_accuracy"`
			Fouls          int         `json:"fouls"`
			YellowCards    int         `json:"yellow_cards"`
			RedCards       int         `json:"red_cards"`
			OffSides       int         `json:"off_sides"`
			Corners        int         `json:"corners"`
			CreatedAt      time.Time   `json:"created_at"`
			UpdatedAt      time.Time   `json:"updated_at"`
		} `json:"away"`
		Date       time.Time `json:"date"`
		Status     string    `json:"status"`
		UniqueLink string    `json:"unique_link"`
		Stadium    string    `json:"stadium"`
		Referee    string    `json:"referee"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestAddFixture(t *testing.T) {
	// Create a new request with the necessary payload
	payload := map[string]interface{}{
		"competition_id": "65ffc4d497df3fa190013424",
		"home_team_id":   "65fdd0c9193055e3c7a77b08",
		"away_team_id":   "65fdd0c9193055e3c7a77af7",
		"date":           "2022-03-28T15:04:05Z",
		"status":         "ongoing",
		"stadium":        "stadium",
		"referee":        "jill jackman",
		"home": map[string]interface{}{
			"substitutes": []string{"John", "Emma", "Michael", "Sophia", "William"},
			"lineup":      []string{"sophie", "William", "Olivia", "James", "Amelia", "Benjamin", "Isabella", "Tim", "Joe", "Fred", "Joyboy"},
			"formation":   "4-4-2",
		},
		"away": map[string]interface{}{
			"substitutes": []string{"John", "Emma", "Michael", "Sophia", "William"},
			"lineup":      []string{"sophie", "William", "Olivia", "James", "Amelia", "Benjamin", "Isabella", "Tim", "Joe", "Fred", "Joyboy"},
			"formation":   "4-3-3",
		},
	}
	payloadBytes, err := json.Marshal(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://league-ybry.onrender.com/api/v1/fixtures", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response Fixture
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully created fixture", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	fmt.Println(response.Data)
}

func TestEditFixture(t *testing.T) {

	// Create a new request with the necessary payload
	payload := map[string]string{
		"competition_id": "65fdfb3556a912cf828b39ec",
		"home_team_id":   "65fdd0c9193055e3c7a77b08",
		"away_team_id":   "65fdd0c9193055e3c7a77af6",
		"date":           "2022-03-28T15:04:05Z",
		"status":         "ongoing",
		"stadium":        "stadium",
		"referee":        "jill jackman",
		"unique_link":    "xlpu0tMy2A",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)
	// Create a new HTTP request
	req, err := http.NewRequest("PATCH", "https://league-ybry.onrender.com/api/v1/fixtures/6604c1ce6912f7e29fee1ec0", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response Fixture
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully updated fixture", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestViewFixture(t *testing.T) {
	req, err := http.NewRequest("GET", "https://league-ybry.onrender.com/api/v1/fixtures/B5VrnxYhBQ", nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response AggFixture
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully fetched fixture", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type DeleteFixture struct {
	Data       interface{} ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

func TestRemoveFixture(t *testing.T) {
	req, err := http.NewRequest("DELETE", "https://league-ybry.onrender.com/api/v1/fixtures/6604c1ce6912f7e29fee1ec0", nil)
	assert.NoError(t, err)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response DeleteFixture
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, "successfully deleted fixture", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestGetPendingFixtures(t *testing.T) {
	req, err := http.NewRequest("GET", "https://league-ybry.onrender.com/api/v1/fixtures/status/pending", nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response StatusFixtures
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully fetched response", response.Message)
	assert.Greater(t, len(response.Data.Data), 0, "Array length should be greater than zero")
}

type StatusFixtures struct {
	Data struct {
		Data []struct {
			ID            string `json:"_id"`
			CompetitionID string `json:"competition_id"`
			HomeTeamID    string `json:"home_team_id"`
			AwayTeamID    string `json:"away_team_id"`
			Home          struct {
				Goals          int         `json:"goals"`
				GoalScorers    interface{} `json:"goal_scorers"`
				Substitutes    []string    `json:"substitutes"`
				Lineup         []string    `json:"lineup"`
				Formation      string      `json:"formation"`
				Shots          int         `json:"shots"`
				ShotsOnTarget  int         `json:"shots_on_target"`
				Possession     float64     `json:"possession"`
				Passes         int         `json:"passes"`
				PassesAccuracy int         `json:"passes_accuracy"`
				Fouls          int         `json:"fouls"`
				YellowCards    int         `json:"yellow_cards"`
				RedCards       int         `json:"red_cards"`
				OffSides       int         `json:"off_sides"`
				Corners        int         `json:"corners"`
				CreatedAt      time.Time   `json:"created_at"`
				UpdatedAt      time.Time   `json:"updated_at"`
			} `json:"home"`
			Away struct {
				Goals          int         `json:"goals"`
				GoalScorers    interface{} `json:"goal_scorers"`
				Substitutes    []string    `json:"substitutes"`
				Lineup         []string    `json:"lineup"`
				Formation      string      `json:"formation"`
				Shots          int         `json:"shots"`
				ShotsOnTarget  int         `json:"shots_on_target"`
				Possession     float64     `json:"possession"`
				Passes         int         `json:"passes"`
				PassesAccuracy int         `json:"passes_accuracy"`
				Fouls          int         `json:"fouls"`
				YellowCards    int         `json:"yellow_cards"`
				RedCards       int         `json:"red_cards"`
				OffSides       int         `json:"off_sides"`
				Corners        int         `json:"corners"`
				CreatedAt      time.Time   `json:"created_at"`
				UpdatedAt      time.Time   `json:"updated_at"`
			} `json:"away"`
			Date       time.Time `json:"date"`
			Status     string    `json:"status"`
			UniqueLink string    `json:"unique_link"`
			Stadium    string    `json:"stadium"`
			Referee    string    `json:"referee"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
		} `json:"data"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Total   int `json:"total"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestGetCompletedFixtures(t *testing.T) {
	req, err := http.NewRequest("GET", "https://league-ybry.onrender.com/api/v1/fixtures/status/completed", nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response StatusFixtures
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully fetched response", response.Message)
	assert.Greater(t, len(response.Data.Data), 0, "Array length should be greater than zero")

}

type Fixtures struct {
	Data struct {
		Data []struct {
			ID            string `json:"ID"`
			CompetitionID struct {
				ID        string    `json:"_id"`
				Name      string    `json:"name"`
				Type      string    `json:"type"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"competition_id"`
			HomeTeamID struct {
				ID          string    `json:"_id"`
				Name        string    `json:"name"`
				State       string    `json:"state"`
				Country     string    `json:"country"`
				FoundedYear int       `json:"founded_year"`
				Stadium     string    `json:"stadium"`
				Sponsor     string    `json:"sponsor"`
				CreatedBy   string    `json:"created_by"`
				CreatedAt   time.Time `json:"created_at"`
				UpdatedAt   time.Time `json:"updated_at"`
			} `json:"home_team_id"`
			AwayTeamID struct {
				ID          string    `json:"_id"`
				Name        string    `json:"name"`
				State       string    `json:"state"`
				Country     string    `json:"country"`
				FoundedYear int       `json:"founded_year"`
				Stadium     string    `json:"stadium"`
				Sponsor     string    `json:"sponsor"`
				CreatedBy   string    `json:"created_by"`
				CreatedAt   time.Time `json:"created_at"`
				UpdatedAt   time.Time `json:"updated_at"`
			} `json:"away_team_id"`
			Home struct {
				Goals          int           `json:"goals"`
				GoalScorers    []interface{} `json:"goal_scorers"`
				Substitutes    []string      `json:"substitutes"`
				Lineup         []string      `json:"lineup"`
				Formation      string        `json:"formation"`
				Shots          int           `json:"shots"`
				ShotsOnTarget  int           `json:"shots_on_target"`
				Possession     float64       `json:"possession"`
				Passes         int           `json:"passes"`
				PassesAccuracy int           `json:"passes_accuracy"`
				Fouls          int           `json:"fouls"`
				YellowCards    int           `json:"yellow_cards"`
				RedCards       int           `json:"red_cards"`
				OffSides       int           `json:"off_sides"`
				Corners        int           `json:"corners"`
				CreatedAt      time.Time     `json:"created_at"`
				UpdatedAt      time.Time     `json:"updated_at"`
			} `json:"home"`
			Away struct {
				Goals          int           `json:"goals"`
				GoalScorers    []interface{} `json:"goal_scorers"`
				Substitutes    []string      `json:"substitutes"`
				Lineup         []string      `json:"lineup"`
				Formation      string        `json:"formation"`
				Shots          int           `json:"shots"`
				ShotsOnTarget  int           `json:"shots_on_target"`
				Possession     float64       `json:"possession"`
				Passes         int           `json:"passes"`
				PassesAccuracy int           `json:"passes_accuracy"`
				Fouls          int           `json:"fouls"`
				YellowCards    int           `json:"yellow_cards"`
				RedCards       int           `json:"red_cards"`
				OffSides       int           `json:"off_sides"`
				Corners        int           `json:"corners"`
				CreatedAt      time.Time     `json:"created_at"`
				UpdatedAt      time.Time     `json:"updated_at"`
			} `json:"away"`
			Date       time.Time `json:"date"`
			Status     string    `json:"status"`
			UniqueLink string    `json:"unique_link"`
			Stadium    string    `json:"stadium"`
			Referee    string    `json:"referee"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
		} `json:"data"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Total   int `json:"total"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestGetAllFixtures(t *testing.T) {
	resp, err := http.Get("https://league-ybry.onrender.com/api/v1/fixtures")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response Fixtures
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "successfully fetched response", response.Message)
	assert.Greater(t, len(response.Data.Data), 0, "Array length should be greater than zero")
}
