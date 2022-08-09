package handlers

import (
	"dearDoctor/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (h handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to login, Try again")
		//responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	admin := models.Admin{}
	err = json.Unmarshal(body, &admin)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to login, Try again")
		//responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	h.SignInAdmin(admin.Username, admin.Password)
	//token, err := h.SignInAdmin(admin.Username, admin.Password)
	// if err != nil {
	// 	formattedError := formaterror.FormatError(err.Error())
	// 	json.NewEncoder(w).Encode("Status:Failure")
	// 	json.NewEncoder(w).Encode("Invalid Username or Password")
	// 	//responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
	// 	return
	// }
	json.NewEncoder(w).Encode("Status:Succuess")
	json.NewEncoder(w).Encode("Welcome,Logged in successfully")
	//responses.JSON(w, http.StatusOK, token)
}
func (h handler) SignInAdmin(username, password string) (string, error) {

	var err error

	admin := models.Admin{}

	err = h.DB.Debug().Model(models.User{}).Where("username = ?", username).Take(&admin).Error
	if err != nil {
		return "", err
	}
	err = models.PasswordVerifyAdmin(admin.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return admin.Username, err
}

func (h handler) AddDept(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var dept models.Departments
	json.Unmarshal(body, &dept)

	if result := h.DB.Create(&dept); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("DeptCreated")
	json.NewEncoder(w).Encode(dept)
}

func (h handler) AppoveAndFee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email, _ := vars["email"]

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var updatedDoctor models.Doctor
	json.Unmarshal(body, &updatedDoctor)

	var doctor models.Doctor

	result := h.DB.First(&doctor, email)
	if result == nil {
		json.NewEncoder(w).Encode("Invalid Email id")
		json.NewEncoder(w).Encode(result.Error)
	} else {
		if err := h.DB.Where(models.Doctor{Email: updatedDoctor.Email}).
			Assign(models.Doctor{Email: updatedDoctor.Email, Approvel: updatedDoctor.Approvel, Fee: updatedDoctor.Fee}).
			FirstOrCreate(&models.Doctor{}).Error; err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Updated")
	}

}

func (h handler) ListAllUsers(w http.ResponseWriter, r *http.Request) {

	var users []models.User

	if result := h.DB.Find(&users); result.Error != nil {
		json.NewEncoder(w).Encode(result.Error)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h handler) GetDoctors(w http.ResponseWriter, r *http.Request) {

	var doctors []models.Doctor

	if result := h.DB.Find(&doctors); result.Error != nil {
		json.NewEncoder(w).Encode(result.Error)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(doctors)
}
