package models

import (
	"time"
)

type Task struct {
	ID          uint `gorm:"column:id;primaryKey"`
	UserID      uint `gorm:"column:user_id"`
	User        User
	Title       string    `gorm:"column:title"`
	Description *string   `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
