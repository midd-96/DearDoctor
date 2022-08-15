package handler

import (
	"dearDoctor/common/response"
	"dearDoctor/model"
	"dearDoctor/service"
	"dearDoctor/utils"
	"encoding/json"
	"net/http"
	//"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	AdminLogin() http.HandlerFunc
	UserLogin() http.HandlerFunc
	UserSignup() http.HandlerFunc
	DoctorSignup() http.HandlerFunc
	DoctorLogin() http.HandlerFunc
}

type authHandler struct {
	jwtAdminService  service.JWTService
	jwtUserService   service.JWTService
	jwtDoctorService service.JWTService
	authService      service.AuthService
	adminService     service.AdminService
	userService      service.UserService
	doctorService    service.DoctorService
	//validate        *validator.Validate
}

func NewAuthHandler(
	jwtAdminService service.JWTService,
	jwtUserService service.JWTService,
	jwtDoctorService service.JWTService,
	authService service.AuthService,
	adminService service.AdminService,
	userService service.UserService,
	doctorService service.DoctorService,
	//validate *validator.Validate,

) AuthHandler {
	return &authHandler{
		jwtAdminService:  jwtAdminService,
		jwtUserService:   jwtUserService,
		jwtDoctorService: jwtDoctorService,
		authService:      authService,
		adminService:     adminService,
		userService:      userService,
		doctorService:    doctorService,
		//validate:        validate,
	}
}

func (c *authHandler) AdminLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var adminLogin model.Admin

		json.NewDecoder(r.Body).Decode(&adminLogin)

		//verifying  admin credentials
		err := c.authService.VerifyAdmin(adminLogin.Username, adminLogin.Password)

		if err != nil {
			response := response.ErrorResponse("Failed to login", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		//getting admin values
		admin, _ := c.adminService.FindAdmin(adminLogin.Username)
		token := c.jwtAdminService.GenerateToken(admin.ID, admin.Username, "admin")
		admin.Password = ""
		admin.Token = token
		response := response.SuccessResponse(true, "OK!", admin.Token)
		utils.ResponseJSON(w, response)
	}

}

func (c *authHandler) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userLogin model.User

		json.NewDecoder(r.Body).Decode(&userLogin)

		//verify User details
		err := c.authService.VerifyUser(userLogin.Email, userLogin.Password)

		if err != nil {
			response := response.ErrorResponse("Failed to login", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		//fetching user details
		user, _ := c.userService.FindUser(userLogin.Email)
		token := c.jwtUserService.GenerateToken(user.ID, user.Email, "user")
		user.Password = ""
		user.Token = token
		response := response.SuccessResponse(true, "OK", user.Token)
		utils.ResponseJSON(w, response)
	}
}

func (c *authHandler) UserSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newUser model.User

		//fetching data
		json.NewDecoder(r.Body).Decode(&newUser)

		err := c.userService.CreateUser(newUser)

		if err != nil {
			response := response.ErrorResponse("Failed to create user", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		user, _ := c.userService.FindUser(newUser.Email)
		user.Password = ""
		response := response.SuccessResponse(true, "OK", user)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *authHandler) DoctorSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var newDoctor model.Doctor

		//fetching data
		json.NewDecoder(r.Body).Decode(&newDoctor)

		err := c.doctorService.CreateDoctor(newDoctor)

		if err != nil {
			response := response.ErrorResponse("Failed to signup", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			utils.ResponseJSON(w, response)
			return
		}

		doctor, _ := c.doctorService.FindDoctor(newDoctor.Email)
		doctor.Password = ""
		response := response.SuccessResponse(true, "OK", doctor)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utils.ResponseJSON(w, response)
	}
}

func (c *authHandler) DoctorLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var doctorLogin model.Doctor

		json.NewDecoder(r.Body).Decode(&doctorLogin)

		//verify doctor Credentials
		err := c.authService.VerifyDoctor(doctorLogin.Email, doctorLogin.Password)

		if err != nil {
			response := response.ErrorResponse("Failed to login", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		//fetching doctor details
		doctor, _ := c.doctorService.FindDoctor(doctorLogin.Email)
		token := c.jwtDoctorService.GenerateToken(doctor.ID, doctor.Email, "doctor")
		doctor.Password = ""
		doctor.Token = token
		response := response.SuccessResponse(true, "OK", doctor.Token)
		utils.ResponseJSON(w, response)
	}
}
