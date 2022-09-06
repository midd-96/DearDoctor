package handler

import (
	"dearDoctor/common/response"
	"dearDoctor/model"
	"dearDoctor/service"
	"dearDoctor/utils"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/razorpay/razorpay-go"
)

type UserHandler interface {
	ConfirmAppointment() http.HandlerFunc
	SendVerificationMail() http.HandlerFunc
	VerifyAccount() http.HandlerFunc
	Payment() http.HandlerFunc
	PaymentSuccess() http.HandlerFunc
	Success() http.HandlerFunc
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (c *userHandler) Success() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		parsedTemplate, _ := template.ParseFiles("template/success.html")
		parsedTemplate.Execute(w, nil)

	}
}

func (c *userHandler) PaymentSuccess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data model.PaymentDetails

		// json.NewDecoder(r.Body).Decode(&data)
		log.Println("here it comes")
		data.User_ID, _ = strconv.Atoi(r.URL.Query().Get("user_id"))

		data.Razorpay_payment_id = r.URL.Query().Get("payment_id")

		data.Razorpay_order_id = r.URL.Query().Get("order_id")

		data.Razorpay_signature = r.URL.Query().Get("signature")

		data.Appointment_ID, _ = strconv.Atoi(r.URL.Query().Get("id"))

		data.PaymentType = "razor_pay"

		data.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		data.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		err := c.userService.AddPayment(data)

		if err != nil {
			response := response.ErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		response := response.SuccessResponse(true, "OK!", "Payment successful")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *userHandler) Payment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type PageVariables struct {
			UserID           int
			OrderIdCreated   string
			TotalPrice       float64
			Name             string
			Email            string
			Phone_Number     int
			AmountInSubUnits float64
			OrderId          int
		}

		var requestData model.PaymentDetails

		client := razorpay.NewClient("rzp_test_kt3cXZneHJI2uV", "OUu6OR0p6chSb32gmuPzjW9o")

		requestData.Appointment_ID, _ = strconv.Atoi(chi.URLParam(r, "Appointment_id"))

		//err := c.userService.CheckAlreadyPaid(requestData.Appointment_ID)

		paymentData, err := c.userService.ProcessingPayment(requestData)

		if err != nil {
			response := response.ErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		amountInSubUnits := paymentData.Amount * 100

		data := map[string]interface{}{
			"amount":          amountInSubUnits,
			"currency":        "INR",
			"receipt":         "some_receipt_id",
			"payment_capture": 1,
		}

		body, err := client.Order.Create(data, nil)

		if err != nil {
			response := response.ErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		val := body["id"]

		orderIDCreated := val.(string)
		pageVariables := PageVariables{
			UserID:           paymentData.User_ID,
			OrderIdCreated:   orderIDCreated,
			TotalPrice:       paymentData.Amount,
			AmountInSubUnits: amountInSubUnits,
			Name:             paymentData.Full_Name,
			Email:            paymentData.Email,
			Phone_Number:     paymentData.Phone_Number,
			OrderId:          paymentData.Appointment_ID,
		}

		parsedTemplate, err := template.ParseFiles("template/app.html")
		parsedTemplate.Execute(w, pageVariables)

		if err != nil {
			response := response.ErrorResponse("error processing request", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}
	}
}

func (c *userHandler) VerifyAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")
		code, _ := strconv.Atoi(r.URL.Query().Get("Code"))

		err := c.userService.VerifyAccount(email, code)
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

func (c *userHandler) SendVerificationMail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("Email")

		_, err := c.userService.FindUser(email)

		if err == nil {
			err = c.userService.SendVerificationEmail(email)
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
		response := response.SuccessResponse(true, "SUCCESS", "Appointment Confirmed")
		w.Header().Add("Content-Type", "application/json")
		utils.ResponseJSON(w, response)

	}
}
