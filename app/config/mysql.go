package config

import (
	"os"
	"fmt"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
  "gorm.io/driver/mysql"

	"app/domain/model"
)

func NewMysqlConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_ID"), 	// シングルクォート => more than one character in rune literal
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ Logger: logger.Default.LogMode(logger.Info) })
	if err != nil {
		panic("failed to connect database")
	}

	return db, err
}

func ExecuteMigrate(db *gorm.DB) (error) {
	err := db.AutoMigrate(
		&model.PostCode{},
		&model.Tag{},
		&model.Todo{},
		&model.UserAddress{},
		&model.UserInfo{},
		&model.User{},
	)
	return err
}