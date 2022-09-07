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
	"strconv"
	"time"
)

type UserService interface {
	FindUser(email string) (*model.UserResponse, error)
	CreateUser(newUser model.User) error
	AddAppointment(confirm model.Confirmed) error
	ViewAppointments(email string) ([]model.Appointments, error)
	SendVerificationEmail(email string) error
	VerifyAccount(email string, code int) error
	ProcessingPayment(data model.PaymentDetails) (*model.PaymentDetails, error)
	AddPayment(data model.PaymentDetails) error
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

func (c *userService) ProcessingPayment(data model.PaymentDetails) (*model.PaymentDetails, error) {

	var err error
	data.Amount, err = c.userRepo.FindAppointmentById(data.Appointment_ID)

	if err != nil {
		return nil, errors.New("Unable to find appointment/ Already paid")
	}

	user, _ := c.userRepo.FindUserByAppointmentId(data.Appointment_ID)

	data.User_ID = user.ID

	data.Email = user.Email

	data.Full_Name = user.First_Name

	data.Phone_Number, _ = strconv.Atoi(user.Phone)

	return &data, nil

}

func (c *userService) AddPayment(data model.PaymentDetails) error {

	err := c.userRepo.Payment(data)

	if err != nil {
		return err
	}

	user, _ := c.userRepo.FindUserById(data.User_ID)

	message := fmt.Sprintf(
		"Hello, %s %s ..\nYour Appointment (Appointment Id: %d) has been confirmed.\n\nConsultation Fee has been paid.\n\nThanks for using dearDoctor Visit Again,\n\n with regards dearDoctor.",
		user.First_Name,
		user.Last_Name,
		data.Appointment_ID,
	)

	err = c.mailConfig.SendMail(user.Email, message)

	if err != nil {
		return err
	}

	return nil
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
		"\nThe verification code is:\n\n%d.\nUse to verify your account.\n Thank you for using dearDoctor.\n with regards Team dearDoctor.",
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

func (c *userService) ViewAppointments(email string) ([]model.Appointments, error) {

	appointments, err := c.userRepo.ViewAppointments(email)
	if err != nil {
		return nil, err
	}
	return appointments, nil
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
