package db

import (
	"github.com/amirdaraby/go-todo-list-api/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {

	// register auto migrate models here
	err := db.AutoMigrate(
		models.User{},
		models.Task{},
	)

	return err
}
