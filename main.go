package main

import (
	//"dearDoctor/db"
	//"dearDoctor/handlers"
	"dearDoctor/routers"
	//"log"
	//"net/http"
	//"github.com/gorilla/mux"
)

func main() {

	// DB := db.Init()
	// h := handlers.New(DB)
	// router := mux.NewRouter()
	routers.InitializeAdminRouter()

	// router.HandleFunc("/signupUser", h.SignupUser).Methods(http.MethodPost)
	// router.HandleFunc("/guestBooking", h.GuestBooking).Methods(http.MethodPost)

	// router.HandleFunc("/lists/{id}", h.DeleteList).Methods(http.MethodDelete)

	// router.HandleFunc("/booking/confirm", h.ConfirmAppointment).Methods(http.MethodPost)

	// log.Println("API is running")
	// http.ListenAndServe(":4000", router)
}
