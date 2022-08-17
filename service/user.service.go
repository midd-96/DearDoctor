package service

import (
	"crypto/md5"
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/repo"
	"errors"
	"fmt"
)

type UserService interface {
	FindUser(email string) (*model.UserResponse, error)
	CreateUser(newUser model.User) error
	AddAppointment(confirm model.Confirmed) error
}

type userService struct {
	userRepo  repo.UserRepository
	adminRepo repo.AdminRepository
}

func NewUserService(
	userRepo repo.UserRepository,
	adminRepo repo.AdminRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		adminRepo: adminRepo,
	}
}

func (c *userService) FindUser(email string) (*model.UserResponse, error) {
	user, err := c.userRepo.FindUser(email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *userService) CreateUser(newUser model.User) error {

	_, err := c.userRepo.FindUser(newUser.Email)

	if err == nil {
		return errors.New("Username already exists")
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	//hashing password
	newUser.Password = HashPassword(newUser.Password)

	_, err = c.userRepo.InsertUser(newUser)
	if err != nil {
		return err
	}
	return nil

}
func (c *userService) AddAppointment(confirm model.Confirmed) error {
	_, err := c.userRepo.AddAppointment(confirm)

	if err != nil {
		return err
	}

	return nil
}

// func (c *userService) GetAllProducts(filter model.Filter, user_id int, pagenation utils.Filter) (*[]model.GetProduct, *utils.Metadata, error) {

// 	products, metadata, err := c.productRepo.GetAllProducts(filter, user_id, pagenation)

// 	if err != nil {
// 		return nil, &metadata, err

// 	}

// 	return &products, &metadata, nil

// }

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}
