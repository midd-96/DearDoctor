package handlers

import (
	// "crypto/rand"

	"dearDoctor/auth"
	"dearDoctor/models"
	"dearDoctor/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func (h handler) SignupDoctor(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read data, Try again"))
		log.Fatalln(err)
	}

	var doctor models.Doctor
	json.Unmarshal(body, &doctor)
	hashedPassword, err := util.HashPassword(doctor.Password)
	if err != nil {
		json.NewEncoder(w).Encode("Internal server error can't convert to hash password")
		return
	}
	doctor.Password = hashedPassword

	valid_email := util.ValidateEmail(doctor.Email)
	if valid_email == false {
		json.NewEncoder(w).Encode("Email not valid")
		return
	}

	valid_phone := util.ValidatePhone(doctor.Phone)
	if valid_phone == false {
		json.NewEncoder(w).Encode("Phone number is not valid")
		return
	}

	result := h.DB.Create(&doctor)

	if result.Error != nil {
		json.NewEncoder(w).Encode("Can't signup")
		json.NewEncoder(w).Encode(result.Error)

	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Created")
		json.NewEncoder(w).Encode(doctor.Last_name)

	}

}
func (h handler) LoginDoctor(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read data, Try again"))
		return
	}
	doctor := models.Doctor{}
	err = json.Unmarshal(body, &doctor)
	if err != nil {
		util.Respond(w, util.Message(false, "Failed to read json data, Try again"))
		return
	}

	token, err := h.SignInDoctor(doctor.Email, doctor.Password)
	if err != nil {
		util.Respond(w, util.Message(false, "Invalid Username or Password"))
		return
	}
	util.Respond(w, util.Message(true, "LoggedIn Successfully"))
	json.NewEncoder(w).Encode(doctor.Email)
	json.NewEncoder(w).Encode(token)
}

func (h handler) SignInDoctor(email, password string) (string, error) {

	var err error

	doctor := models.Doctor{}

	err = h.DB.Debug().Model(models.Doctor{}).Where("email = ?", email).Take(&doctor).Error
	if err != nil {
		return "", err
	}
	err = util.CheckPassword(doctor.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(doctor.Email, (os.Getenv("SECRET_DOCTOR")))

}

func (h handler) MarkAvilability(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var slot models.Slotes
	json.Unmarshal(body, &slot)

	if result := h.DB.Create(&slot); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Slotes updated")
	json.NewEncoder(w).Encode(slot)
}
