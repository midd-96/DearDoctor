package service

import (
	"dearDoctor/model"
	"dearDoctor/repo"
	"dearDoctor/utils"
	"errors"
	"log"
)

type AdminService interface {
	FindAdmin(username string) (*model.AdminResponse, error)
	AllUsers(pagenation utils.Filter) (*[]model.UserResponse, *utils.Metadata, error)
	UpdateApproveFee(approvel model.ApproveAndFee, emailid string) error
	AddDept(department model.Departments) error
	AllDoctors(pagenation utils.Filter) (*[]model.DoctorResponse, *utils.Metadata, error)
}

type adminService struct {
	adminRepo  repo.AdminRepository
	userRepo   repo.UserRepository
	doctorRepo repo.DoctorRepository
}

func NewAdminService(
	adminRepo repo.AdminRepository,
	userRepo repo.UserRepository,
	doctorRepo repo.DoctorRepository) AdminService {
	return &adminService{
		adminRepo:  adminRepo,
		userRepo:   userRepo,
		doctorRepo: doctorRepo,
	}
}

func (c *adminService) FindAdmin(username string) (*model.AdminResponse, error) {
	admin, err := c.adminRepo.FindAdmin(username)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (c *adminService) AllUsers(pagenation utils.Filter) (*[]model.UserResponse, *utils.Metadata, error) {

	users, metadata, err := c.userRepo.AllUsers(pagenation)
	if err != nil {
		return nil, &metadata, err
	}

	return &users, &metadata, nil
}

func (c *adminService) AllDoctors(pagenation utils.Filter) (*[]model.DoctorResponse, *utils.Metadata, error) {

	doctors, metadata, err := c.doctorRepo.AllDoctors(pagenation)
	if err != nil {
		return nil, &metadata, err
	}

	return &doctors, &metadata, nil
}

func (c *adminService) AddDept(department model.Departments) error {

	err := c.adminRepo.AddDept(department)

	if err != nil {
		log.Println(err)
		return errors.New("error in adding new department")
	}

	return nil
}

func (c *adminService) UpdateApproveFee(approvel model.ApproveAndFee, emailid string) error {

	err := c.adminRepo.UpdateApproveFee(approvel, emailid)

	if err != nil {
		log.Println(err)
		return errors.New("error while updating approvel/fee")
	}

	return nil
}
