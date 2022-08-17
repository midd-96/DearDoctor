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
}

type adminHandler struct {
	adminService service.AdminService
	userService  service.UserService
}

func NewAdminHandler(
	adminService service.AdminService,
	userService service.UserService,
) AdminHandler {
	return &adminHandler{
		adminService: adminService,
		userService:  userService,
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
		//log.Print(result)
		//users, err := c.adminService.AllUsers()

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

		response := response.SuccessResponse(true, "OK!", "new department added")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *adminHandler) ApprovelAndFee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var datatoadd model.ApproveAndFee

		json.NewDecoder(r.Body).Decode(&datatoadd)

		err := c.adminService.UpdateApproveFee(datatoadd)

		if err != nil {
			response := response.ErrorResponse("error while updating ", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "OK!", "Updated approvel/fee")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}
