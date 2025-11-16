package repositories

import (
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"gorm.io/gorm"
)

type GormTeamRepository struct {
	db *gorm.DB
}

func NewGormTeamRepository(db *gorm.DB) TeamRepository {
	return &GormTeamRepository{db: db}
}

func (g *GormTeamRepository) CreateTeam(team *models.Team, users []models.User) error {
	exists, err := g.TeamExists(team.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.NewTeamExists(team.Name)
	}
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		for i := range users {
			users[i].TeamID = team.ID
			if err := tx.Create(&users[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GormTeamRepository) GetTeamByName(name string) (*models.Team, error) {
	var team models.Team
	res := g.db.Where("name = ?", name).First(&team)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound()
		}
		return nil, res.Error
	}
	return &team, nil
}

func (g *GormTeamRepository) GetTeamUsers(teamID string) ([]models.User, error) {
	var users []models.User
	res := g.db.Where("team_id = ?", teamID).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (g *GormTeamRepository) TeamExists(name string) (bool, error) {
	var ans int64
	res := g.db.Model(&models.Team{}).Where("name = ?", name).Count(&ans)
	if res.Error != nil {
		return false, res.Error
	}
	return ans > 0, nil
}
