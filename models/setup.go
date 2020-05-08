package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func SetupModels() *gorm.DB {
	db, err := gorm.Open("postgres", "user=ajdelgados dbname=user-go sslmode=disable")

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})

	return db
}
