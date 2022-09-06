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
	CalculatePayout() http.HandlerFunc
	ViewSingleUser() http.HandlerFunc
	ViewSingleDoctor() http.HandlerFunc
	ApprovePayout() http.HandlerFunc
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

func (c *adminHandler) ApprovePayout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")

		amount, err := c.adminService.ApprovePayout(email)

		if err != nil {

			response := response.ErrorResponse("Request not accepted", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "Amount Transferred", amount)
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}
}

func (c *adminHandler) ViewSingleDoctor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		doc_Id, _ := strconv.Atoi(r.URL.Query().Get("Doctor_Id"))

		user, err := c.adminService.ViewSingleDoctor(doc_Id)

		if err != nil {

			response := response.ErrorResponse("could not process the request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "DOCTOR DETAILS", user)
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}
}
func (c *adminHandler) ViewSingleUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		doc_Id, _ := strconv.Atoi(r.URL.Query().Get("User_Id"))

		user, err := c.adminService.ViewSingleUser(doc_Id)

		if err != nil {

			response := response.ErrorResponse("could not process the request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "USER DETAILS", user)
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}
}
func (c *adminHandler) CalculatePayout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doc_Id, _ := strconv.Atoi(r.URL.Query().Get("Doctor_Id"))

		amount, err := c.adminService.CalculatePayout(doc_Id)

		if err != nil {
			response := response.ErrorResponse("failed to calculate payout", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "SUCCESS", "Payout calculated")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, amount)
		utils.ResponseJSON(w, response)

	}
}

func (c *adminHandler) ViewAllAppointments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//doc_id, _ := strconv.Atoi(r.URL.Query().Get("Doc_Id"))
		page, _ := strconv.Atoi(r.URL.Query().Get("Page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("Pagesize"))
		//day := (r.URL.Query().Get("Day"))
		log.Println(page, " ", pageSize)
		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		var filters model.Filter

		json.NewDecoder(r.Body).Decode(&filters)

		appointments, metadata, err := c.adminService.ViewAllAppointments(pagenation, filters)

		result := struct {
			FilteredAppointment *[]model.AppointmentByDoctor
			Meta                *utils.Metadata
		}{
			FilteredAppointment: appointments,
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
