package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"ID"`
	UserName  string    `gorm:"column:user_name;unique;index" json:"user_name" validate:"required,min=2,max=255"`
	Password  string    `gorm:"column:password" json:"password" validate:"required,min=8,max=255"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Tasks     []Task    `json:"tasks"`
}
