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
	ViewAllAppointments(pagenation utils.Filter, filters model.Filter) (*[]model.AppointmentByDoctor, *utils.Metadata, error)
	CalculatePayout(doc_Id int) (string, error)
	ViewSingleUser(user_Id int) (*model.UserResponse, error)
	ViewSingleDoctor(doc_Id int) (*model.DoctorResponse, error)
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

func (c *adminService) ViewSingleDoctor(doc_Id int) (*model.DoctorResponse, error) {
	doctor, err := c.adminRepo.ViewSingleDoctor(doc_Id)

	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (c *adminService) ViewSingleUser(user_Id int) (*model.UserResponse, error) {
	user, err := c.adminRepo.ViewSingleUser(user_Id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *adminService) CalculatePayout(doc_Id int) (string, error) {

	amount, err := c.adminRepo.CalculatePayout(doc_Id)

	if err != nil {
		return "", err
	}

	return amount, nil

}

func (c *adminService) ViewAllAppointments(pagenation utils.Filter, filters model.Filter) (*[]model.AppointmentByDoctor, *utils.Metadata, error) {

	appointments, metadata, err := c.adminRepo.ViewAllAppointments(pagenation, filters)

	if err != nil {
		return nil, &metadata, err
	}

	return &appointments, &metadata, nil
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
