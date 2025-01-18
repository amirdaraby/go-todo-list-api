package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	UserName  string    `gorm:"column:user_name;unique;index"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Tasks     []Task
}
