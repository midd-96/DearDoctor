package main

import (
	"dearDoctor/db"
	"dearDoctor/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	DB := db.Init()
	h := handlers.New(DB)
	router := mux.NewRouter()

	router.HandleFunc("/listAllUsers", h.ListAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/lists/{id}", h.GetList).Methods(http.MethodGet)
	router.HandleFunc("/signupUser", h.SignupUser).Methods(http.MethodPost)
	router.HandleFunc("/guestBooking", h.GuestBooking).Methods(http.MethodPost)
	router.HandleFunc("/addDepartment", h.AddDept).Methods(http.MethodPost)
	router.HandleFunc("/signupDoctor", h.SignupDoctor).Methods(http.MethodPost)
	router.HandleFunc("/approveAndFee/{email}", h.AppoveAndFee).Methods(http.MethodPatch)
	router.HandleFunc("/lists/{id}", h.DeleteList).Methods(http.MethodDelete)
	router.HandleFunc("/doctor/addslotes", h.MarkAvilability).Methods(http.MethodPost)
	router.HandleFunc("/booking/confirm", h.ConfirmAppointment).Methods(http.MethodPost)

	log.Println("API is running")
	http.ListenAndServe(":4000", router)
}
