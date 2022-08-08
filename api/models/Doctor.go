package models

import (
	"dd/project/api/utils"
	"dd/project/api/utils/formaterror"
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

type Doctor struct {
	gorm.Model

	Id         int    `json:"id" gorm:"primary_key"`
	First_name string `json:"first_name" gorm:"not null"`
	Last_name  string `json:"last_name" gorm:"not null"`
	Email      string `json:"email" gorm:"unique;not null" valid:"email"`
	Phone      string `json:"phone" gorm:"unique;not null"`
	Password   string `json:"password" gorm:"not null" valid:"length(6|20)"`
	//Dep_code       string `json:"dep_code"`
	Department     string `json:"department" gorm:"not null"`
	Specialization string `json:"specialization"`
	Approvel       bool   `json:"approvel" gorm:"default:false"`
	Fee            int    `json:"fee"`
	Role           int    `gorm:"default:3" `
}

func (d *Doctor) Validate(action string) error {
	switch strings.ToLower(action) {
	case "email":
		if !strings.Contains(d.Email, "@") && strings.Contains(d.Email, ".") {
			return errors.New("Email is not valid")
		}
		return nil
	case "phone":
		if !strings.Contains(d.Email, "@") && strings.Contains(d.Email, ".") {
			return errors.New("Phone Number is not valid")
		}
		return nil

	default:
		if !strings.Contains(d.Email, "@") && strings.Contains(d.Email, ".") {
			return errors.New("Email is not valid")
		}
		if !strings.Contains(d.Email, "@") && strings.Contains(d.Email, ".") {
			return errors.New("Phone Number is not valid")
		}
		return nil
	}

}

func (d *Doctor) SaveDoctor(db *gorm.DB) (*Doctor, error) {
	hashedPassword, err := utils.HashPassword(d.Password)
	if err != nil {
		formaterror.FormatError("hashPassword")
	}
	d.Password = hashedPassword
	err = db.Debug().Create(&d).Error
	if err != nil {
		return &Doctor{}, err
	}
	return d, nil
}

func (d *Doctor) FindAllDoctors(db *gorm.DB) (*[]Doctor, error) {
	var err error
	doctors := []Doctor{}
	err = db.Debug().Model(&Doctor{}).Find(&doctors).Error
	if err != nil {
		return &[]Doctor{}, err
	}
	return &doctors, err
}

func (d *Doctor) UpdateADoctor(db *gorm.DB) (*Doctor, error) {

	var err error

	err = db.Debug().Model(&Doctor{}).Where("email = ?", d.Email).Updates(Doctor{Approvel: d.Approvel, Fee: d.Fee}).Error
	if err != nil {
		return &Doctor{}, err
	}
	if d.Email == "" {
		err = db.Debug().Model(&Doctor{}).Where("email = ?", d.Email).Take(&d.Approvel, &d.Fee).Error
		if err != nil {
			return &Doctor{}, err
		}
	}
	return d, nil
}
