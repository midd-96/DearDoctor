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

func (h handler) SignupUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var user models.User
	json.Unmarshal(body, &user)

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		json.NewEncoder(w).Encode("Internal server error can't convert to hash password")
		return
	}
	user.Password = hashedPassword

	valid_email := util.ValidateEmail(user.Email)
	if valid_email == false {
		json.NewEncoder(w).Encode("Email not valid")
		return
	}
	valid_phone := util.ValidatePhone(user.Phone)
	if valid_phone == false {
		json.NewEncoder(w).Encode("Phone number is not valid")
		return
	}

	result := h.DB.Create(&user)

	if result.Error != nil {
		json.NewEncoder(w).Encode("Can't signup")
		json.NewEncoder(w).Encode(result.Error)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Created")
		json.NewEncoder(w).Encode(user.Last_name)
	}

}

func (h handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		//responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		//responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	h.SignInDoctor(user.Email, user.Password)

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

func (h handler) SignInUser(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = h.DB.Debug().Model(models.Doctor{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = util.CheckPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	//return auth.CreateToken(uint32(doctor.Id))
	return user.First_name, nil
}

func (h handler) ConfirmAppointment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var confirm models.Confirmed
	json.Unmarshal(body, &confirm)

	if result := h.DB.Create(&confirm); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Booking confirmed")
	json.NewEncoder(w).Encode(confirm)
}
