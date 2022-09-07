package routes

import (
	h "dearDoctor/handler"
	m "dearDoctor/middleware"

	"github.com/go-chi/chi"
)

type UserRoute interface {
	UserRouter(router chi.Router,
		authHandler h.AuthHandler,
		middleware m.Middleware,
		userHandler h.UserHandler,
	)
}

type userRoute struct{}

func NewUserRoute() UserRoute {
	return &userRoute{}
}

func (r *userRoute) UserRouter(routes chi.Router,
	authHandler h.AuthHandler,
	middleware m.Middleware,
	userHandler h.UserHandler) {

	routes.Post("/user/signup", authHandler.UserSignup())
	routes.Post("/user/login", authHandler.UserLogin())
	routes.Post("/user/send/verification", userHandler.SendVerificationMail())
	routes.Patch("/user/verify/account", userHandler.VerifyAccount())
	routes.Get("/user/payment/{Appointment_id}", userHandler.Payment())
	routes.Get("/payment-success", userHandler.PaymentSuccess())
	routes.Get("/success", userHandler.Success())
	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Post("/user/add/appointment", userHandler.ConfirmAppointment())
		r.Get("/user/token/refresh", authHandler.UserRefreshToken())
		r.Get("/user/appointments/confirmed", userHandler.ViewConfirmedAppointment())

	})

}
