package db

import (
	"fmt"
	"log"
	"os"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Connect() (*DB, error) {
	connStr := os.Getenv("CONN_STR")

	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected")

	// Используем автомиграции для создания схем бд
	err = gormDB.AutoMigrate(&models.Team{}, &models.User{}, &models.PullRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations completed")

	// Создаём кастомные индексы на голом SQL
	if err := createCustomIndexes(gormDB); err != nil {
		return nil, fmt.Errorf("failed to create custom indexes: %w", err)
	}
	log.Println("Indexes created")

	return &DB{gormDB}, nil
}

func createCustomIndexes(db *gorm.DB) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_team_id ON users(team_id)`,
		`CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_reviewers ON pull_requests USING GIN(reviewers)`,
	}
	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return fmt.Errorf("failed to create index %s: %w", indexSQL, err)
		}
	}
	log.Println("Custom indexes created")
	return nil
}

func (db *DB) Close() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Printf("Error with sql.DB in gorm.DB:%s", err)
	}
	if err = sqlDB.Close(); err != nil {
		log.Printf("Database close error:%s", err)
	}
}
