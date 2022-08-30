package service

import (
	"crypto/md5"
	"database/sql"
	"dearDoctor/config"
	"dearDoctor/model"
	"dearDoctor/repo"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type UserService interface {
	FindUser(email string) (*model.UserResponse, error)
	CreateUser(newUser model.User) error
	AddAppointment(confirm model.Confirmed) error
	SendVerificationEmail(email string) error
	VerifyAccount(email string, code int) error
}

type userService struct {
	userRepo   repo.UserRepository
	adminRepo  repo.AdminRepository
	mailConfig config.MailConfig
}

func NewUserService(
	userRepo repo.UserRepository,
	adminRepo repo.AdminRepository,
	mailConfig config.MailConfig) UserService {
	return &userService{
		userRepo:   userRepo,
		adminRepo:  adminRepo,
		mailConfig: mailConfig,
	}
}

func (c *userService) VerifyAccount(email string, code int) error {

	err := c.userRepo.VerifyAccount(email, code)

	if err != nil {
		return err
	}
	return nil
}

func (c *userService) SendVerificationEmail(email string) error {
	//to generate random code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(100000)

	message := fmt.Sprintf(
		"\nThe verification code is:\n\n%d\nUseto verify your account.\n\n dearDoctor.",
		code,
	)

	// send random code to user's email
	if err := c.mailConfig.SendMail(email, message); err != nil {
		return err
	}

	err := c.userRepo.StoreVerificationDetails(email, code)

	if err != nil {
		return err
	}

	return nil
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

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}
