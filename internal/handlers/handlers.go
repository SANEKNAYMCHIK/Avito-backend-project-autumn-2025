package handlers

import (
	"net/http"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.ReviewService
}

func NewHandler(service services.ReviewService) *Handler {
	return &Handler{service: service}
}

// POST /team/add
func (h *Handler) CreateTeam(c *gin.Context) {
	var req models.CreateTeamRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid request body"))
		return
	}
	result, err := h.service.CreateTeam(req.TeamName, req.Members)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"team": result,
	})
}

// GET /team/get?team_name=<team name>
func (h *Handler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		_ = c.Error(errors.NewInvalidInput("team_name parameter is required"))
		return
	}
	result, err := h.service.GetTeam(teamName)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /users/setIsActive
func (h *Handler) SetUserActive(c *gin.Context) {
	var req models.SetActiveRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid request body"))
		return
	}
	result, err := h.service.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": result,
	})
}

// POST /pullRequest/create
func (h *Handler) CreatePR(c *gin.Context) {
	var req models.CreatePRRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid request body"))
		return
	}
	result, err := h.service.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"pr": result,
	})
}

// POST /pullRequest/merge
func (h *Handler) MergePR(c *gin.Context) {
	var req models.MergePRRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid request body"))
		return
	}
	result, err := h.service.MergePR(req.PullRequestID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pr": result,
	})
}

// POST /pullRequest/reassign
func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req models.ReassignRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid request body"))
		return
	}
	result, err := h.service.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /users/getReview?user_id=<user id>
func (h *Handler) GetUserReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		_ = c.Error(errors.NewInvalidInput("user_id parameter is required"))
		return
	}
	result, err := h.service.GetUserReviews(userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}
