// main_test.go

package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"league/fixtures"
	"league/models"
	"log"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	// login creds
	payload := map[string]string{
		"email":    "johndoe@gmail.com",
		"password": "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/user/login", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

type UserResponse struct {
	Data       models.User ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

func TestSignup(t *testing.T) {
	// signup creds
	payload := map[string]string{
		"first_name": "john",
		"last_name":  "doe",
		"email":      "johndoe45@gmail.com",
		"password":   "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/user/signup", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("signup request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response UserResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully created user", response.Message)
}

func TestSignupAdmin(t *testing.T) {
	// signup creds
	payload := map[string]string{
		"first_name": "john",
		"last_name":  "doe",
		"email":      "Admin017@gmail.com",
		"password":   "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/admin/signup", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("signup request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response UserResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully created user", response.Message)
}

type LoginAdmin struct {
	Data       interface{} ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

func TestLoginAdmin(t *testing.T) {
	// signup creds
	payload := map[string]string{
		"email":    "Admin01@gmail.com",
		"password": "123456",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Simulate login request
	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/admin/login", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response LoginAdmin
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully sent otp", response.Message)

}

type ConfirmLoginAdmin struct {
	Data       string ` json:"data"`
	Message    string ` json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestConfirmLoginAdmin(t *testing.T) {
	payload := map[string]string{
		"email": "Admin01@gmail.com",
		"otp":   "5704",
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post("https://league-ybry.onrender.com/api/v1/auth/admin/login/confirm", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response ConfirmLoginAdmin
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "Successfully logged In admin", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type Team struct {
	Data       models.Team ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

func TestAddTeam(t *testing.T) {
	jwtToken := ""

	// Create a new request with the necessary payload
	payload := map[string]interface{}{
		"name":         "test",
		"country":      "england",
		"state":        "london",
		"founded_year": 2001,
		"stadium":      "stamford",
		"sponsor":      "adidas",
	}
	payloadBytes, _ := json.Marshal(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://league-ybry.onrender.com/api/v1/teams", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Team
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully created team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)

}

func TestEditTeam(t *testing.T) {
	jwtToken := ""

	// Create a new request with the necessary payload
	payload := map[string]interface{}{
		"name":         "test",
		"country":      "england",
		"state":        "london",
		"founded_year": 2001,
		"stadium":      "stamford",
		"sponsor":      "adidas",
	}
	payloadBytes, _ := json.Marshal(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("PATCH", "https://league-ybry.onrender.com/api/v1/teams/65fdab162a01779985f840f9", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Team
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully updated team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestViewTeam(t *testing.T) {
	jwtToken := ""
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "https://league-ybry.onrender.com/api/v1/teams/65fdab162a01779985f840f9", nil)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Team
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestRemoveTeam(t *testing.T) {
	jwtToken := ""

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "https://league-ybry.onrender.com/api/v1/teams/65fdab162a01779985f840f9", nil)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Team
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully deleted team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type SearchTeam struct {
	Data struct {
		Data    []models.Team `json: "data"`
		Page    int           `json:"page"`
		Perpage int           `json:"per_page"`
		Total   int           `json:"total"`
	} ` json:"data"`
	Message    string ` json:"message"`
	StatusCode int    `json:"statusCode"`
}

func TestSearchTeam(t *testing.T) {
	resp, err := http.Get("https://league-ybry.onrender.com/api/v1/teams?page=1&per_page=15")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Team
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully fetched teams", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)

}

type Fixture struct {
	Data       models.Fixture ` json:"data"`
	Message    string         ` json:"message"`
	StatusCode int            `json:"statusCode"`
}

type AggFixture struct {
	Data       fixtures.Fixture ` json:"data"`
	Message    string           ` json:"message"`
	StatusCode int              `json:"statusCode"`
}

func TestAddFixture(t *testing.T) {
	jwtToken := ""

	// Create a new request with the necessary payload
	payload := map[string]interface{}{
		"competition_id": "65fdfb3556a912cf828b39ec",
		"home_team_id":   "65fdd0c9193055e3c7a77af5",
		"away_team_id":   "65fdd0c9193055e3c7a77af6",
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
	payloadBytes, _ := json.Marshal(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://league-ybry.onrender.com/api/v1/fixtures", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Fixture
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully created fixture", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestEditFixture(t *testing.T) {
	jwtToken := ""

	// Create a new request with the necessary payload
	payload := map[string]string{
		"competition_id": "65fdfb3556a912cf828b39ec",
		"home_team_id":   "65fdd0c9193055e3c7a77af5",
		"away_team_id":   "65fdd0c9193055e3c7a77af6",
		"date":           "2022-03-28T15:04:05Z",
		"status":         "ongoing",
		"stadium":        "stadium",
		"referee":        "jill robertt",
		"unique_link":    "doOKx-eJeg",
	}
	payloadBytes, _ := json.Marshal(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("PATCH", "https://league-ybry.onrender.com/api/v1/teams/65fdab162a01779985f840f9", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response Team
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully updated team", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestViewFixture(t *testing.T) {
	jwtToken := ""
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "https://league-ybry.onrender.com/api/v1/fixtures/QIe0HHgoJA", nil)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response AggFixture
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type DeleteFixture struct {
	Data       interface{} ` json:"data"`
	Message    string      ` json:"message"`
	StatusCode int         `json:"statusCode"`
}

func TestRemoveFixture(t *testing.T) {
	jwtToken := ""

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "https://league-ybry.onrender.com/api/v1/fixtures/65ff5651f0dfef2f33dcec19", nil)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("create team request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: " + err.Error())
	}

	var response DeleteFixture
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Error decoding JSON: " + err.Error())
	}
	assert.Equal(t, "successfully deleted fixture", response.Message)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestGetPendingFixtures(t *testing.T) {
	resp, err := http.Get("https://league-ybry.onrender.com/api/v1/fixtures/status/pending")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetCompletedFixtures(t *testing.T) {
	resp, err := http.Get("https://league-ybry.onrender.com/api/v1/fixtures/status/completed")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

}

func TestGetAllFixtures(t *testing.T) {
	resp, err := http.Get("https://league-ybry.onrender.com/api/v1/fixtures")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

}
