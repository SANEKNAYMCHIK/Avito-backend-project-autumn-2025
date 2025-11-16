package repositories

import "gorm.io/gorm"

type Repository struct {
	Team TeamRepository
	User UserRepository
	PR   PRRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Team: NewGormTeamRepository(db),
		User: NewGormUserRepository(db),
		PR:   NewGormPRRepository(db),
	}
}
