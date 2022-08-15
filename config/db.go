package config

import (
	"dearDoctor/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dbURL := "postgres://postgres:secret@localhost:5432/postgres"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Departments{})
	db.AutoMigrate(&model.Doctor{})
	//db.AutoMigrate(&model.Guests{})
	db.AutoMigrate(&model.Slotes{})
	db.AutoMigrate(&model.Confirmed{})
	return db
}
