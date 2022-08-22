package handler

import (
	"dearDoctor/common/response"
	"dearDoctor/model"
	"dearDoctor/service"
	"dearDoctor/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type AdminHandler interface {
	ViewAllUsers() http.HandlerFunc
	AddDepartment() http.HandlerFunc
	ApprovelAndFee() http.HandlerFunc
	ViewAllDoctors() http.HandlerFunc
	ViewAllAppointments() http.HandlerFunc
}

type adminHandler struct {
	adminService  service.AdminService
	userService   service.UserService
	doctorService service.DoctorService
}

func NewAdminHandler(
	adminService service.AdminService,
	userService service.UserService,
	doctorService service.DoctorService,
) AdminHandler {
	return &adminHandler{
		adminService:  adminService,
		userService:   userService,
		doctorService: doctorService,
	}
}

func (c *adminHandler) ViewAllAppointments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		doc_id, _ := strconv.Atoi(r.URL.Query().Get("Doc_Id"))
		page, _ := strconv.Atoi(r.URL.Query().Get("Page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("Pagesize"))
		day := (r.URL.Query().Get("Day"))
		log.Println(page, " ", pageSize, " ", doc_id, " ", day)
		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		appointments, metadata, err := c.adminService.ViewAllAppointments(pagenation, doc_id, day)

		result := struct {
			AppointmentByDoctor *[]model.AppointmentByDoctor
			Meta                *utils.Metadata
		}{
			AppointmentByDoctor: appointments,
			Meta:                metadata,
		}

		if err != nil {

			response := response.ErrorResponse("could not process the request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "LIST OF ALL APPOINTMENTS", result)
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}

}

func (c *adminHandler) ViewAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		log.Println(page, "   ", pageSize)

		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		users, metadata, err := c.adminService.AllUsers(pagenation)

		result := struct {
			Users *[]model.UserResponse
			Meta  *utils.Metadata
		}{
			Users: users,
			Meta:  metadata,
		}

		if err != nil {
			response := response.ErrorResponse("error while getting users from database", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "Listed All Users", result)
		utils.ResponseJSON(w, response)

	}
}

func (c *adminHandler) ViewAllDoctors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		log.Println(page, "   ", pageSize)

		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		doctors, metadata, err := c.adminService.AllDoctors(pagenation)

		result := struct {
			Doctors *[]model.DoctorResponse
			Meta    *utils.Metadata
		}{
			Doctors: doctors,
			Meta:    metadata,
		}

		if err != nil {
			response := response.ErrorResponse("error while getting docters from database", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "Listed All Doctors", result)
		utils.ResponseJSON(w, response)

	}
}

func (c *adminHandler) AddDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newdepartment model.Departments

		json.NewDecoder(r.Body).Decode(&newdepartment)

		err := c.adminService.AddDept(newdepartment)

		if err != nil {
			response := response.ErrorResponse("failed to add new department", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "SUCCESS", "new department added")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *adminHandler) ApprovelAndFee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var datatoadd model.ApproveAndFee
		emailid := (r.URL.Query().Get("email"))

		json.NewDecoder(r.Body).Decode(&datatoadd)

		err := c.adminService.UpdateApproveFee(datatoadd, emailid)

		if err != nil {
			response := response.ErrorResponse("error while updating ", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "SUCCESS", "Updated approvel/fee")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}
