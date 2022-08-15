package handler

import (
	"dearDoctor/common/response"
	"dearDoctor/model"
	"dearDoctor/service"
	"dearDoctor/utils"
	"encoding/json"
	"net/http"
)

type UserHandler interface {
	ConfirmAppointment() http.HandlerFunc
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (c *userHandler) ConfirmAppointment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appointment model.Confirmed

		json.NewDecoder(r.Body).Decode(&appointment)

		err := c.userService.AddAppointment(appointment)

		if err != nil {
			response := response.ErrorResponse("Appointment not confirmed", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}
		response := response.SuccessResponse(true, "OK", "Appointment Confirmed")
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}
}
