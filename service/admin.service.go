package service

import (
	"dearDoctor/model"
	"dearDoctor/repo"
	"errors"
	"log"
)

type AdminService interface {
	FindAdmin(username string) (*model.AdminResponse, error)
	AllUsers() (*[]model.UserResponse, error)
	UpdateApproveFee(approvel model.ApproveAndFee) error
	AddDept(department model.Departments) error
}

type adminService struct {
	adminRepo repo.AdminRepository
	userRepo  repo.UserRepository
}

func NewAdminService(
	adminRepo repo.AdminRepository,
	userRepo repo.UserRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}

func (c *adminService) FindAdmin(username string) (*model.AdminResponse, error) {
	admin, err := c.adminRepo.FindAdmin(username)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (c *adminService) AllUsers() (*[]model.UserResponse, error) {

	users, err := c.userRepo.AllUsers()
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (c *adminService) AddDept(department model.Departments) error {

	err := c.adminRepo.AddDept(department)

	if err != nil {
		log.Println(err)
		return errors.New("error in adding new department")
	}

	return nil
}

func (c *adminService) UpdateApproveFee(approvel model.ApproveAndFee) error {

	err := c.adminRepo.UpdateApproveFee(approvel)

	if err != nil {
		log.Println(err)
		return errors.New("error while updating approvel/fee")
	}

	return nil
}
