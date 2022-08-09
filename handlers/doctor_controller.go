package handlers

import (
	// "crypto/rand"

	"dearDoctor/models"
	"dearDoctor/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h handler) SignupDoctor(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
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
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		//responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	doctor := models.Doctor{}
	err = json.Unmarshal(body, &doctor)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		//responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	h.SignInDoctor(doctor.Email, doctor.Password)

	// token, err := server.SignInDoctor(user.Email, user.Password)
	// if err != nil {

	// 	formattedError := formaterror.FormatError(err.Error())
	// 	json.NewEncoder(w).Encode("Status:Failure")
	// 	json.NewEncoder(w).Encode("Invalid Username / Password")
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
	// 	return
	// }
	// responses.JSON(w, http.StatusOK, token)
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
	//return auth.CreateToken(uint32(doctor.Id))
	return doctor.First_name, nil
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
