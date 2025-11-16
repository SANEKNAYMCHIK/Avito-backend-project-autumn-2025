package repositories

import "github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"

type TeamRepository interface {
	CreateTeam(team *models.Team, users []models.User) error
	GetTeamByName(name string) (*models.Team, error)
	GetTeamUsers(teamID string) ([]models.User, error)
	TeamExists(name string) (bool, error)
}

type UserRepository interface {
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetActiveUsersByTeam(teamID string) ([]models.User, error)
	GetUserTeam(userID string) (*models.Team, error)
}

type PRRepository interface {
	CreatePR(pr *models.PullRequest) error
	GetPRByID(id string) (*models.PullRequest, error)
	UpdatePR(pr *models.PullRequest) error
	GetPRsByReviewer(userID string) ([]models.PullRequest, error)
	PRExists(id string) (bool, error)
}
