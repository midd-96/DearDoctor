package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	gorm.Model

	Username string `json:"username" gorm:"primarykey"`
	Password string `json:"password" gorm:"not null"`
	Role     int    `json:"role" gorm:"defualt:1"`
}

func PasswordVerifyAdmin(ogPassword, password string) error {
	res := ogPassword == password
	if res != true {
		return errors.New("Incorrect Password")
	}
	return nil
}
