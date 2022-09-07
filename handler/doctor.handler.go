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

type DoctorHandler interface {
	MarkAvailability() http.HandlerFunc
	AppointmentsByDoctor() http.HandlerFunc
	SendVerificationMail() http.HandlerFunc
	VerifyAccount() http.HandlerFunc
	RequestForPayout() http.HandlerFunc
	AddBankAccountDetails() http.HandlerFunc
}

type doctorHandler struct {
	doctorService service.DoctorService
}

func NewDoctorHandler(doctorService service.DoctorService) DoctorHandler {
	return &doctorHandler{
		doctorService: doctorService,
	}
}

func (c *doctorHandler) AddBankAccountDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bankAccount model.Account

		json.NewDecoder(r.Body).Decode(&bankAccount)

		bankAccount.Email = r.Header.Get("email")

		err := c.doctorService.AddBankAccountDetails(bankAccount)

		if err != nil {
			response := response.ErrorResponse("failed to add bank account details", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "SUCCESS", "Bank account added")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *doctorHandler) RequestForPayout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("email")
		requestAmount, err := strconv.ParseFloat(r.URL.Query().Get("Amount"), 64)

		amount, err := c.doctorService.RequestForPayout(email, requestAmount)

		if err != nil {
			response := response.ErrorResponse("error while requesting for payout", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "Requested successfully", amount)
		utils.ResponseJSON(w, response)

	}
}

func (c *doctorHandler) VerifyAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")
		code, _ := strconv.Atoi(r.URL.Query().Get("Code"))

		err := c.doctorService.VerifyAccount(email, code)
		log.Println(err)

		if err != nil {
			response := response.ErrorResponse("Verification failed, Invalid OTP", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "Account verified successfully", email)
		utils.ResponseJSON(w, response)
	}
}

func (c *doctorHandler) SendVerificationMail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")

		_, err := c.doctorService.FindDoctor(email)

		if err == nil {
			err = c.doctorService.SendVerificationEmail(email)
		}

		if err != nil {
			response := response.ErrorResponse("Error while sending verification mail", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "Verification mail sent successfully", email)
		utils.ResponseJSON(w, response)
	}
}
func (c *doctorHandler) AppointmentsByDoctor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		docId, _ := strconv.Atoi(r.URL.Query().Get("docId"))

		pagenation := utils.Filter{
			Page:     page,
			PageSize: pageSize,
		}

		appointments, metadata, err := c.doctorService.AppointmentsByDoctor(pagenation, docId)

		result := struct {
			Appointments *[]model.Appointments
			Meta         *utils.Metadata
		}{
			Appointments: appointments,
			Meta:         metadata,
		}

		if err != nil {
			response := response.ErrorResponse("error while getting appointments from database", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "Listed All Appointments", result)
		utils.ResponseJSON(w, response)

	}
}

func (c *doctorHandler) MarkAvailability() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var addslotes model.Slotes

		json.NewDecoder(r.Body).Decode(&addslotes)

		addslotes.Docter_id, _ = strconv.Atoi(r.Header.Get("user_id"))

		err := c.doctorService.AddSlotes(addslotes)

		if err != nil {
			response := response.ErrorResponse("failed to add avialable slotes", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "SUCCESS", "available slotes added")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}
