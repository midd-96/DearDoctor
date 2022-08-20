package handler

import (
	"dearDoctor/common/response"
	"dearDoctor/model"
	"dearDoctor/service"
	"dearDoctor/utils"
	"encoding/json"
	"log"
	"net/http"
)

type DoctorHandler interface {
	MarkAvailability() http.HandlerFunc
}

type doctorHandler struct {
	doctorService service.DoctorService
}

func NewDoctorHandler(doctorService service.DoctorService) DoctorHandler {
	return &doctorHandler{
		doctorService: doctorService,
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
