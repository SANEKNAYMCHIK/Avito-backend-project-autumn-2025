package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock сервиса
type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) CreateTeam(teamName string, members []models.TeamMember) (*models.TeamResponse, error) {
	args := m.Called(teamName, members)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TeamResponse), args.Error(1)
}

func (m *MockReviewService) GetTeam(teamName string) (*models.TeamResponse, error) {
	args := m.Called(teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TeamResponse), args.Error(1)
}

func (m *MockReviewService) SetUserActive(userID string, isActive bool) (*models.UserResponse, error) {
	args := m.Called(userID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockReviewService) CreatePR(prID, title, authorID string) (*models.PullRequestResponse, error) {
	args := m.Called(prID, title, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PullRequestResponse), args.Error(1)
}

func (m *MockReviewService) MergePR(prID string) (*models.PullRequestResponse, error) {
	args := m.Called(prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PullRequestResponse), args.Error(1)
}

func (m *MockReviewService) ReassignReviewer(prID, oldReviewerID string) (*models.ReassignResponse, error) {
	args := m.Called(prID, oldReviewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReassignResponse), args.Error(1)
}

func (m *MockReviewService) GetUserReviews(userID string) (*models.UserPRsResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPRsResponse), args.Error(1)
}

func TestHandler_CreateTeam_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewHandler(mockService)

	requestBody := models.CreateTeamRequest{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
		},
	}

	expectedResponse := &models.TeamResponse{
		TeamName: "backend",
		Members:  requestBody.Members,
	}

	// Mock expectations
	mockService.On("CreateTeam", "backend", requestBody.Members).Return(expectedResponse, nil)

	// Create request
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.POST("/team/add", handler.CreateTeam)

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "team")
	mockService.AssertExpectations(t)
}

func TestHandler_CreateTeam_InvalidJSON(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewHandler(mockService)

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/team/add", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.POST("/team/add", handler.CreateTeam)

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	errorObj := response["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_INPUT", errorObj["code"])
}

func TestHandler_CreateTeam_ServiceError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewHandler(mockService)

	requestBody := models.CreateTeamRequest{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
		},
	}

	// Mock expectations
	mockService.On("CreateTeam", "backend", requestBody.Members).Return(
		nil, errors.NewTeamExists("backend"),
	)

	// Create request
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.POST("/team/add", handler.CreateTeam)

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	errorObj := response["error"].(map[string]interface{})
	assert.Equal(t, "TEAM_EXISTS", errorObj["code"])

	mockService.AssertExpectations(t)
}

func TestHandler_GetTeam_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewHandler(mockService)

	expectedResponse := &models.TeamResponse{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
		},
	}

	// Mock expectations
	mockService.On("GetTeam", "backend").Return(expectedResponse, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/team/get?team_name=backend", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.GET("/team/get", handler.GetTeam)

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TeamResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "backend", response.TeamName)
	assert.Len(t, response.Members, 1)

	mockService.AssertExpectations(t)
}

func TestHandler_GetTeam_MissingParameter(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewHandler(mockService)

	// Create request without team_name parameter
	req, _ := http.NewRequest("GET", "/team/get", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.GET("/team/get", handler.GetTeam)

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
	errorObj := response["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_INPUT", errorObj["code"])
}

func TestHandler_SetUserActive_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewHandler(mockService)

	requestBody := models.SetActiveRequest{
		UserID:   "u1",
		IsActive: false,
	}

	expectedResponse := &models.UserResponse{
		UserId:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: false,
	}

	// Mock expectations
	mockService.On("SetUserActive", "u1", false).Return(expectedResponse, nil)

	// Create request
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.POST("/users/setIsActive", handler.SetUserActive)

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "user")
	mockService.AssertExpectations(t)
}
