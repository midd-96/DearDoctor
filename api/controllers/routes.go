package controllers

import (
	"dd/project/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//admin routes
	s.Router.HandleFunc("/admin/createdept", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateDept))).Methods("POST")
	s.Router.HandleFunc("/admin/login", middlewares.SetMiddlewareJSON(s.LoginAdmin)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/user/login", middlewares.SetMiddlewareJSON(s.LoginUser)).Methods("POST")
	s.Router.HandleFunc("/user/signup", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/user/listall", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/user/listone/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/user/booking/confirm", s.ConfirmReservation).Methods("POST")

	//Doctors routes
	s.Router.HandleFunc("/doctor/login", middlewares.SetMiddlewareJSON(s.LoginDoctor)).Methods("POST")
	s.Router.HandleFunc("/doctor/signup", middlewares.SetMiddlewareJSON(s.CreateDoctor)).Methods("POST")
	s.Router.HandleFunc("/doctors/update/{email}", middlewares.SetMiddlewareJSON(s.UpdateDoctor)).Methods("PATCH")
	s.Router.HandleFunc("/doctors/listall", middlewares.SetMiddlewareJSON(s.GetDoctors)).Methods("GET")
	s.Router.HandleFunc("/doctor/addslotes", middlewares.SetMiddlewareJSON(s.CreateAailableSlot)).Methods("POST")
	//Posts routes
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")
}
