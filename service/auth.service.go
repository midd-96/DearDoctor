package service

import (
	"crypto/md5"
	"dearDoctor/repo"
	"errors"
	"fmt"
)

type AuthService interface {
	VerifyAdmin(email string, password string) error
	VerifyUser(email string, password string) error
	VerifyDoctor(email string, password string) error
}

type authService struct {
	adminRepo  repo.AdminRepository
	userRepo   repo.UserRepository
	doctorRepo repo.DoctorRepository
}

func NewAuthService(
	adminRepo repo.AdminRepository,
	userRepo repo.UserRepository,
	doctorRepo repo.DoctorRepository,
) AuthService {
	return &authService{
		adminRepo:  adminRepo,
		userRepo:   userRepo,
		doctorRepo: doctorRepo,
	}
}

func (c *authService) VerifyAdmin(email, password string) error {

	admin, err := c.adminRepo.FindAdmin(email)

	//_, err = c.adminRepo.FindAdmin(email)

	if err != nil {
		return errors.New("Invalid Username/ password, failed to login")
	}

	isValidPassword := VerifyPassword(password, admin.Password)
	if !isValidPassword {
		return errors.New("Invalid username/ Password, failed to login")
	}

	return nil
}

func (c *authService) VerifyUser(email string, password string) error {

	user, err := c.userRepo.FindUser(email)

	if err != nil {
		return errors.New("failed to login. check your email")
	}

	isValidPassword := VerifyPassword(password, user.Password)
	if !isValidPassword {
		return errors.New("failed to login. check your credential")
	}

	return nil
}

func (c *authService) VerifyDoctor(email string, password string) error {

	doctor, err := c.doctorRepo.FindDoctor(email)

	if err != nil {
		return errors.New("failed to login. check your email/password")
	}

	isValidPassword := VerifyPassword(password, doctor.Password)
	if !isValidPassword {
		return errors.New("failed to login. check your emial/password")
	}

	return nil
}

func VerifyPassword(requestPassword, dbPassword string) bool {

	requestPassword = fmt.Sprintf("%x", md5.Sum([]byte(requestPassword)))
	return requestPassword == dbPassword
}
