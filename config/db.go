package config

import (
	"dearDoctor/model"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	// dbURL := "postgresql://postgres:secret@:5432/postgres" //driver://postgres:password@localhost:5432/dbname
	dbURL := os.Getenv("DB_SOURCE")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.Admin{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Departments{})
	db.AutoMigrate(&model.Doctor{})
	db.AutoMigrate(&model.Slotes{})
	db.AutoMigrate(&model.Confirmed{})
	db.AutoMigrate(&model.Verification{})
	db.AutoMigrate(&model.Payout{})
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.Payment{})
	db.Exec("INSERT INTO admins (username,password) VALUES ('admin@gmail.com','admin')")
	return db
}
