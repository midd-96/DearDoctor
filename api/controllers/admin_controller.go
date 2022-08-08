package controllers

import (
	"dd/project/api/auth"
	"dd/project/api/models"
	"dd/project/api/responses"
	"dd/project/api/utils/formaterror"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// to login
func (server *Server) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to login, Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	admin := models.Admin{}
	err = json.Unmarshal(body, &admin)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to login, Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignInAdmin(admin.Username, admin.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Invalid Username or Password")
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	json.NewEncoder(w).Encode("Status:Succuess")
	json.NewEncoder(w).Encode("Welcome,Logged in successfully")
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignInAdmin(username, password string) (string, error) {

	var err error

	admin := models.Admin{}

	err = server.DB.Debug().Model(models.User{}).Where("username = ?", username).Take(&admin).Error
	if err != nil {
		return "", err
	}
	err = models.PasswordVerifyAdmin(admin.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(uint32(admin.Role))
}

// to create or add new department

func (server *Server) CreateDept(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json ,Try again")
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	dept := models.Department{}
	err = json.Unmarshal(body, &dept)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	deptCreated, err := dept.SaveDept(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to create department,Try again")
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	json.NewEncoder(w).Encode("Status:Success")
	json.NewEncoder(w).Encode("New Department Created")
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, dept.ID))
	responses.JSON(w, http.StatusCreated, deptCreated)
}
