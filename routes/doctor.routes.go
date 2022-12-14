package routes

import (
	h "dearDoctor/handler"
	m "dearDoctor/middleware"

	"github.com/go-chi/chi"
)

type DoctorRoute interface {
	DoctorRouter(router chi.Router,
		authHandler h.AuthHandler,
		middleware m.Middleware,
		doctorHandler h.DoctorHandler,
	)
}

type doctorRoute struct{}

func NewDoctorRoute() DoctorRoute {
	return &doctorRoute{}
}

func (r *doctorRoute) DoctorRouter(routes chi.Router,
	authHandler h.AuthHandler,
	middleware m.Middleware,
	doctorHandler h.DoctorHandler) {

	routes.Post("/doctor/signup", authHandler.DoctorSignup())
	routes.Post("/doctor/login", authHandler.DoctorLogin())
	routes.Post("/doctor/send/verification", doctorHandler.SendVerificationMail())
	routes.Patch("/doctor/verify/account", doctorHandler.VerifyAccount())
	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Post("/doctor/add/availability", doctorHandler.MarkAvailability())
		r.Patch("/doctor/request/payout", doctorHandler.RequestForPayout())
		r.Get("/doctor/token/refresh", authHandler.DoctorRefreshToken())
		r.Get("/doctor/list/allappointments", doctorHandler.AppointmentsByDoctor())
		r.Post("/doctor/add/bankac", doctorHandler.AddBankAccountDetails())
	})

}
