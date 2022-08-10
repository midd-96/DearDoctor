package routers

import (
	"dearDoctor/db"
	"dearDoctor/handlers"
	m "dearDoctor/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeAdminRouter() {

	DB := db.Init()
	h := handlers.New(DB)
	router := mux.NewRouter()

	router.HandleFunc("/admin/login", h.LoginAdmin).Methods(http.MethodPost)
	router.HandleFunc("/listAllUsers", m.SetMiddlewareAuthentication(h.ListAllUsers)).Methods(http.MethodGet)
	router.HandleFunc("/doctors/listall", m.SetMiddlewareAuthentication(h.GetDoctors)).Methods(http.MethodGet)
	router.HandleFunc("/lists/{id}", h.GetList).Methods(http.MethodGet)
	router.HandleFunc("/addDepartment", h.AddDept).Methods(http.MethodPost)
	router.HandleFunc("/approveAndFee/{email}", h.AppoveAndFee).Methods(http.MethodPatch)

	log.Println("API is running")
	http.ListenAndServe(":4000", router)
}
