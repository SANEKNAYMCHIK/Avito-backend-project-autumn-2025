package models

import (
	"time"

	"github.com/lib/pq"
)

type Team struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
	Name      string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time ``
	Users     []User    `gorm:"foreignKey:TeamID"`
}

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
	Username  string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	TeamID    string    `gorm:"not null;type:varchar(255)"`
	Team      Team      `gorm:"foreignKey:TeamID"`
	CreatedAt time.Time ``
}

type PullRequest struct {
	ID        string         `gorm:"primaryKey;type:varchar(255)"`
	Title     string         `gorm:"not null"`
	AuthorID  string         `gorm:"not null;type:varchar(255)"`
	Author    User           `gorm:"foreignKey:AuthorID"`
	Status    string         `gorm:"default:'OPEN'"`
	Reviewers pq.StringArray `gorm:"type:text[]"`
	MergedAt  *time.Time     ``
}
