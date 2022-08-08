package models

import (
	"dd/project/api/utils"
	"dd/project/api/utils/formaterror"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Id               int        `json:"id" gorm:"primary_key"`
	First_name       string     `json:"first_name" gorm:"not null"`
	Last_name        string     `json:"last_name" gorm:"not null"`
	Email            string     `json:"email" gorm:"unique;not null" validate:"email"`
	Phone            string     `json:"phone" gorm:"unique;not null"`
	Password         string     `json:"password" gorm:"not null" validate:"min=6)"`
	Last_appointment time.Month `json:"last_appointment"`
	Role             int        `json:"role" gorm:"default:3"`
}

func (u *User) Prepare() {
	u.Id = 0
	u.First_name = html.EscapeString(strings.TrimSpace(u.First_name))
	u.Last_name = html.EscapeString(strings.TrimSpace(u.Last_name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.Last_appointment = 0 //time.Now().Month()
	u.Role = 3
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":

		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.First_name == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		formaterror.FormatError("hashPassword")
	}
	u.Password = hashedPassword

	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"first_name": u.First_name,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
