package service

import (
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/repo"
	"errors"
)

type DoctorService interface {
	AddSlotes(slote model.Slotes) error
	FindDoctor(email string) (*model.DoctorResponse, error)
	CreateDoctor(newDoctor model.Doctor) error
}

type doctorService struct {
	doctorRepo repo.DoctorRepository
	userRepo   repo.UserRepository
}

func NewDoctorService(
	doctorRepo repo.DoctorRepository,
	userRepo repo.UserRepository,
) DoctorService {
	return &doctorService{
		doctorRepo: doctorRepo,
		userRepo:   userRepo,
	}
}

func (c *doctorService) AddSlotes(slote model.Slotes) error {

	_, err := c.doctorRepo.AddSlotes(slote)

	if err != nil {
		return err
	}

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
