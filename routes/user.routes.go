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
	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Post("/user/add/appointment", userHandler.ConfirmAppointment())

	})

}
