package routers

import (
	"dearDoctor/db"
	"dearDoctor/handlers"

	//m "dearDoctor/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeDoctorRouter() {

	DB := db.Init()
	h := handlers.New(DB)
	router := mux.NewRouter()

	router.HandleFunc("/signupDoctor", h.SignupDoctor).Methods(http.MethodPost)
	router.HandleFunc("/doctor/addslotes", h.MarkAvilability).Methods(http.MethodPost)

	log.Println("API is running")
	http.ListenAndServe(":4000", router)
}
