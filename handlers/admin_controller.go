package handlers

import (
	"dearDoctor/auth"
	"dearDoctor/models"
	"dearDoctor/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (h handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read data, Try again"))
		return
	}
	admin := models.Admin{}
	err = json.Unmarshal(body, &admin)
	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read json data, Try again"))
		return
	}
	token, err := h.SignInAdmin(admin.Username, admin.Password)
	if err != nil {
		util.Respond(w, util.Message(false, "Invalid Username or Password"))

		return
	}
	util.Respond(w, util.Message(true, "LoggedIn Successfully"))
	json.NewEncoder(w).Encode(admin.Username)
	json.NewEncoder(w).Encode(token)

}
func (h handler) SignInAdmin(username, password string) (string, error) {

	var err error

	admin := models.Admin{}

	err = h.DB.Debug().Model(models.Admin{}).Where("username = ?", username).Take(&admin).Error
	if err != nil {
		return "", err
	}

	err = models.PasswordVerifyAdmin(admin.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(admin.Username, (os.Getenv("SECRET_ADMIN")))
}

func (h handler) AddDept(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read data, Try again"))
		log.Fatalln(err)
	}

	var dept models.Departments
	json.Unmarshal(body, &dept)

	if result := h.DB.Create(&dept); result.Error != nil {
		util.Respond(w, util.Message(false, "Failed to read json data, Try again"))
		fmt.Println(result.Error)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	util.Respond(w, util.Message(true, "Department Added"))
	json.NewEncoder(w).Encode(dept)
}

func (h handler) AppoveAndFee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email, _ := vars["email"]

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read data, Try again"))
		log.Fatalln(err)
	}

	var updatedDoctor models.Doctor
	json.Unmarshal(body, &updatedDoctor)

	var doctor models.Doctor

	result := h.DB.First(&doctor, email)
	if result == nil {
		util.Respond(w, util.Message(false, "Invalid Email, Try again"))
		json.NewEncoder(w).Encode(result.Error)
	} else {
		if err := h.DB.Where(models.Doctor{Email: updatedDoctor.Email}).
			Assign(models.Doctor{Email: updatedDoctor.Email, Approvel: updatedDoctor.Approvel, Fee: updatedDoctor.Fee}).
			FirstOrCreate(&models.Doctor{}).Error; err != nil {
			util.Respond(w, util.Message(false, "Failed to update data, Try again"))
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		util.Respond(w, util.Message(true, "Updated approvel/fee"))
		json.NewEncoder(w).Encode(doctor.Fee)
		json.NewEncoder(w).Encode(doctor.Approvel)
	}

}

func (h handler) ListAllUsers(w http.ResponseWriter, r *http.Request) {

	var users []models.User

	if result := h.DB.Find(&users); result.Error != nil {
		util.Respond(w, util.Message(false, "Failed to fetch data,Try again"))
		json.NewEncoder(w).Encode(result.Error)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.Respond(w, util.Message(true, "Listed All Users"))
	json.NewEncoder(w).Encode(users)
}

func (h handler) GetDoctors(w http.ResponseWriter, r *http.Request) {

	var doctors []models.Doctor

	if result := h.DB.Find(&doctors); result.Error != nil {
		util.Respond(w, util.Message(false, "Failed to fetch data,Try again"))
		json.NewEncoder(w).Encode(result.Error)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.Respond(w, util.Message(true, "Listed All Doctors"))
	json.NewEncoder(w).Encode(doctors)
}

func (h handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var usertodelete models.User

	if result := h.DB.First(&usertodelete, id); result.Error != nil {
		util.Respond(w, util.Message(false, "Incorrect details, User not found"))
		fmt.Println(result.Error)
	}

	h.DB.Delete(&usertodelete)

	w.WriteHeader(http.StatusOK)
	util.Respond(w, util.Message(true, "Deleted User"))
	json.NewEncoder(w).Encode(usertodelete)
}

func (h handler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var doctortodelete models.Doctor

	if result := h.DB.First(&doctortodelete, id); result.Error != nil {
		util.Respond(w, util.Message(false, "Incorrect details, Doctor not found"))
		fmt.Println(result.Error)
	}

	h.DB.Delete(&doctortodelete)

	w.WriteHeader(http.StatusOK)
	util.Respond(w, util.Message(true, "Deleted Doctor"))
	json.NewEncoder(w).Encode(doctortodelete)
}
