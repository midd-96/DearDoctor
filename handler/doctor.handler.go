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
}

type doctorHandler struct {
	doctorService service.DoctorService
}

func NewDoctorHandler(doctorService service.DoctorService) DoctorHandler {
	return &doctorHandler{
		doctorService: doctorService,
	}
}

func (c *doctorHandler) AppointmentsByDoctor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

		docId, _ := strconv.Atoi(r.URL.Query().Get("docEmail"))

		log.Println(page, "   ", pageSize)

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

		err := c.doctorService.AddSlotes(addslotes)
		log.Println(err)
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
