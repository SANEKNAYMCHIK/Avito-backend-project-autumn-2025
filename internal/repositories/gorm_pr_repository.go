package repositories

import (
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"gorm.io/gorm"
)

type GormPRRepository struct {
	db *gorm.DB
}

func NewGormPRRepository(db *gorm.DB) *GormPRRepository {
	return &GormPRRepository{db: db}
}

func (g *GormPRRepository) PRExists(id string) (bool, error) {
	var count int64
	res := g.db.Model(&models.PullRequest{}).Where("id = ?", id).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count > 0, nil
}

func (g *GormPRRepository) CreatePR(pr *models.PullRequest) error {
	exists, err := g.PRExists(pr.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.NewPRExists(pr.ID)
	}
	result := g.db.Create(pr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (g *GormPRRepository) GetPRByID(id string) (*models.PullRequest, error) {
	var pr models.PullRequest
	res := g.db.Where("id = ?", id).First(&pr)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound()
		}
		return nil, res.Error
	}
	return &pr, nil
}

func (r *GormPRRepository) UpdatePR(pr *models.PullRequest) error {
	res := r.db.Save(pr)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g *GormPRRepository) GetPRsByReviewer(userID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest
	res := g.db.Where("? = ANY(reviewers)", userID).Find(&prs)
	if res.Error != nil {
		return nil, res.Error
	}
	return prs, nil
}
