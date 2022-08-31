package service

import (
	"database/sql"
	"dearDoctor/config"
	"dearDoctor/model"
	"dearDoctor/repo"
	"dearDoctor/utils"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type DoctorService interface {
	AddSlotes(slote model.Slotes) error
	FindDoctor(email string) (*model.DoctorResponse, error)
	CreateDoctor(newDoctor model.Doctor) error
	AppointmentsByDoctor(pagenation utils.Filter, docId int) (*[]model.Appointments, *utils.Metadata, error)
	SendVerificationEmail(email string) error
	VerifyAccount(email string, code int) error
	RequestForPayout(email string, requestAmount float64) (float64, error)
	AddBankAccountDetails(bankAccount model.Account) error
}

type doctorService struct {
	doctorRepo repo.DoctorRepository
	userRepo   repo.UserRepository
	mailConfig config.MailConfig
}

func NewDoctorService(
	doctorRepo repo.DoctorRepository,
	userRepo repo.UserRepository,
	mailConfig config.MailConfig,
) DoctorService {
	return &doctorService{
		doctorRepo: doctorRepo,
		userRepo:   userRepo,
		mailConfig: mailConfig,
	}
}

func (c *doctorService) AddBankAccountDetails(bankAccount model.Account) error {

	bank_Ac := c.doctorRepo.AddBankAccountDetails(bankAccount)

	if bank_Ac != nil {
		return bank_Ac
	}

	return nil
}

func (c *doctorService) RequestForPayout(email string, requestAmount float64) (float64, error) {

	amount, err := c.doctorRepo.RequestForPayout(email, requestAmount)

	if err != nil {
		log.Println("Error from doctor service :", err)
		return amount, err

	}

	return amount, nil
}

func (c *doctorService) VerifyAccount(email string, code int) error {

	err := c.doctorRepo.VerifyAccount(email, code)

	if err != nil {
		return err
	}
	return nil
}

func (c *doctorService) SendVerificationEmail(email string) error {
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

	err := c.doctorRepo.StoreVerificationDetails(email, code)

	if err != nil {
		return err
	}

	return nil
}

func (c *doctorService) AppointmentsByDoctor(pagenation utils.Filter, docId int) (*[]model.Appointments, *utils.Metadata, error) {
	appointments, metadata, err := c.doctorRepo.ListAppointments(pagenation, docId)
	if err != nil {
		return nil, &metadata, err
	}

	return &appointments, &metadata, nil
}

func (c *doctorService) AddSlotes(slote model.Slotes) error {

	slot, err := c.doctorRepo.AddSlotes(slote)

	if err != nil {
		return err
	}
	log.Println(slot)
	return nil
}

func (c *doctorService) CreateDoctor(newDoctor model.Doctor) error {

	_, err := c.userRepo.FindUser(newDoctor.Email)

	if err == nil {
		return errors.New("Username already exists")
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	//hashing password
	newDoctor.Password = HashPassword(newDoctor.Password)

	_, err = c.doctorRepo.InsertDoctor(newDoctor)
	if err != nil {
		return errors.New("Error inserting doctor details in the database")
	}
	return nil

}

func (c *doctorService) FindDoctor(email string) (*model.DoctorResponse, error) {
	doctor, err := c.doctorRepo.FindDoctor(email)

	if err != nil {
		return nil, err
	}

	return &doctor, nil
}
