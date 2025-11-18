package services

import "github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"

type TeamService interface {
	CreateTeam(teamName string, members []models.TeamMember) (*models.TeamResponse, error)
	GetTeam(teamName string) (*models.TeamResponse, error)
}

type UserService interface {
	SetUserActive(userID string, isActive bool) (*models.UserResponse, error)
	GetUserReviews(userID string) (*models.UserPRsResponse, error)
}

type PRService interface {
	CreatePR(prID, title, authorID string) (*models.PullRequestShort, error)
	MergePR(prID string) (*models.PullRequestResponse, error)
	ReassignReviewer(prID, oldReviewerID string) (*models.ReassignResponse, error)
}

type ReviewService interface {
	TeamService
	UserService
	PRService
}
