package models

import (
	"time"
)

type Task struct {
	ID          uint `gorm:"column:id;primaryKey" json:"id"`
	UserID      uint `gorm:"column:user_id" json:"user_id"`
	User        User
	Title       string    `gorm:"column:title" json:"title" validate:"required,min=1,max=40"`
	Description *string   `gorm:"column:description" json:"description" validate:"omitempty,max=255"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}
