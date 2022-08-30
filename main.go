package main

import (
	"database/sql"
	"dearDoctor/config"
	h "dearDoctor/handler"
	m "dearDoctor/middleware"
	"dearDoctor/repo"
	"dearDoctor/routes"
	"dearDoctor/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	//"github.com/go-playground/validator/v10"
	"github.com/subosito/gotenv"
)

//to call functions before main functions
func init() {
	gotenv.Load()
}

func main() {

	//Loading value from env file
	port := os.Getenv("PORT")

	//For making log file
	file, err := os.OpenFile("Logging Details", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Logging in File not done")
	}
	log.SetOutput(file)

	// creating an instance of chi r
	router := chi.NewRouter()

	// using logger to display each request
	router.Use(middleware.Logger)

	config.Init()

	var (
		db         *sql.DB           = config.ConnectDB()
		mailConfig config.MailConfig = config.NewMailConfig()
		//validate    *validator.Validate    = validator.New()
		adminRepo  repo.AdminRepository  = repo.NewAdminRepo(db)
		userRepo   repo.UserRepository   = repo.NewUserRepo(db)
		doctorRepo repo.DoctorRepository = repo.NewDoctorRepo(db)

		jwtAdminService  service.JWTService    = service.NewJWTAdminService()
		jwtUserService   service.JWTService    = service.NewJWTUserService()
		jwtDoctorService service.JWTService    = service.NewJWTDoctorService()
		authService      service.AuthService   = service.NewAuthService(adminRepo, userRepo, doctorRepo)
		adminService     service.AdminService  = service.NewAdminService(adminRepo, userRepo, doctorRepo)
		userService      service.UserService   = service.NewUserService(userRepo, adminRepo, mailConfig)
		doctorService    service.DoctorService = service.NewDoctorService(doctorRepo, userRepo, mailConfig)

		authHandler h.AuthHandler = h.NewAuthHandler(jwtAdminService,
			jwtUserService, jwtDoctorService, authService,
			adminService,
			userService,
			doctorService)
		//validate)
		adminMiddleware  m.Middleware    = m.NewMiddlewareAdmin(jwtAdminService)
		userMiddleware   m.Middleware    = m.NewMiddlewareUser(jwtUserService)
		doctorMiddleware m.Middleware    = m.NewMiddlewareDoctors(jwtDoctorService)
		adminHandler     h.AdminHandler  = h.NewAdminHandler(adminService, userService, doctorService)
		userHandler      h.UserHandler   = h.NewUserHandler(userService)
		doctorHandler    h.DoctorHandler = h.NewDoctorHandler(doctorService)

		adminRoute  routes.AdminRoute  = routes.NewAdminRoute()
		userRoute   routes.UserRoute   = routes.NewUserRoute()
		doctorRoute routes.DoctorRoute = routes.NewDoctorRoute()
	)

	//routing
	adminRoute.AdminRouter(router,
		authHandler,
		adminHandler,
		adminMiddleware,
		doctorHandler)

	userRoute.UserRouter(router,
		authHandler,
		userMiddleware,
		userHandler)

	doctorRoute.DoctorRouter(router,
		authHandler,
		doctorMiddleware,
		doctorHandler)

	log.Println("Api is listening on port:", port)
	// Starting server
	http.ListenAndServe(":"+port, router)

}
