package models

import (
	"time"
)

type Task struct {
	ID          uint       `gorm:"column:id;primaryKey" json:"id"`
	UserID      uint       `gorm:"column:user_id" json:"-"`
	Title       string     `gorm:"column:title" json:"title" validate:"required,min=1,max=40"`
	Description *string    `gorm:"column:description" json:"description" validate:"omitempty,max=255"`
	DoneAt      *time.Time `gorm:"column:done_at;index" json:"done_at"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
