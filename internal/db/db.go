package db

import (
	"fmt"

	"github.com/amirdaraby/go-todo-list-api/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(c config.Config) (db *gorm.DB, err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return
	}

	err = AutoMigrate(db)

	return
}
