package routes

import (
	h "dearDoctor/handler"
	m "dearDoctor/middleware"

	"github.com/go-chi/chi"
)

type AdminRoute interface {
	AdminRouter(routes chi.Router,
		authHandler h.AuthHandler,
		adminHandler h.AdminHandler,
		middleware m.Middleware,
		doctorHandler h.DoctorHandler)
}

type adminRoute struct{}

func NewAdminRoute() AdminRoute {
	return &adminRoute{}
}

// to handle admin routes
func (r *adminRoute) AdminRouter(routes chi.Router,
	authHandler h.AuthHandler,
	adminHandler h.AdminHandler,
	middleware m.Middleware,
	doctorHandler h.DoctorHandler) {

	routes.Post("/admin/login", authHandler.AdminLogin())
	routes.Get("/admin/view/all/appointments", adminHandler.ViewAllAppointments())

	routes.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)

		r.Get("/admin/view/users", adminHandler.ViewAllUsers())
		r.Get("/admin/view/doctors", adminHandler.ViewAllDoctors())
		r.Post("/admin/add/dept", adminHandler.AddDepartment())
		r.Patch("/admin/approve/doctor", adminHandler.ApprovelAndFee())

		r.Get("/admin/token/refresh", authHandler.AdminRefreshToken())

	})

}
