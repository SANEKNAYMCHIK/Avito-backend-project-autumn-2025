package repositories

import (
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (g *GormUserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	res := g.db.Where("id = ?", id).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound()
		}
		return nil, res.Error
	}
	return &user, nil
}

func (g *GormUserRepository) UpdateUser(userID string, isActive bool) error {
	res := g.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", isActive)
	return res.Error
}

func (g *GormUserRepository) GetActiveUsersByTeam(teamID string) ([]models.User, error) {
	var users []models.User
	res := g.db.Where("team_id = ? AND is_active = ?", teamID, true).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (g *GormUserRepository) GetUserTeam(userID string) (*models.Team, error) {
	var user models.User
	var team models.Team
	res := g.db.Where("id = ?", userID).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound()
		}
		return nil, res.Error
	}
	res = g.db.Where("id = ?", user.TeamID).First(&team)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound()
		}
		return nil, res.Error
	}
	return &team, nil
}
