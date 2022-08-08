package controllers

import (
	"dd/project/api/auth"
	"dd/project/api/models"
	"dd/project/api/responses"
	"dd/project/api/utils"
	"dd/project/api/utils/formaterror"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//to signup by doctor
func (server *Server) CreateDoctor(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	doctor := models.Doctor{}
	err = json.Unmarshal(body, &doctor)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	//doctor.Prepare()
	err = doctor.Validate("")
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to Vlidate data")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	doctorCreated, err := doctor.SaveDoctor(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to Signup,Try again")
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	json.NewEncoder(w).Encode("Status:Success")
	json.NewEncoder(w).Encode("Signed up Successfully")
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, doctorCreated.Id))
	responses.JSON(w, http.StatusCreated, doctorCreated)
}

func (server *Server) GetDoctors(w http.ResponseWriter, r *http.Request) {

	doctor := models.Doctor{}

	doctors, err := doctor.FindAllDoctors(server.DB)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to fetch data,Try again")
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, doctors)
	json.NewEncoder(w).Encode("Status:Success")
	json.NewEncoder(w).Encode("Found User")
}

func (server *Server) UpdateDoctor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	em := vars["email"]

	// Check if the post exist
	doctor := models.Doctor{}
	err := server.DB.Debug().Model(models.Post{}).Where("email = ?", em).Take(&doctor).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Doctor not found"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	postUpdate := models.Doctor{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = postUpdate.Validate("email")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.Email = doctor.Email //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdateADoctor(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

//login by doctor

func (server *Server) LoginDoctor(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to validate,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignInDoctor(user.Email, user.Password)
	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Invalid Username / Password")
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignInDoctor(email, password string) (string, error) {

	var err error

	doctor := models.Doctor{}

	err = server.DB.Debug().Model(models.Doctor{}).Where("email = ?", email).Take(&doctor).Error
	if err != nil {
		return "", err
	}
	err = utils.VerifyPassword(doctor.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(uint32(doctor.Id))
}

//to add available slotes

func (server *Server) CreateAailableSlot(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	slot := models.Slote{}
	err = json.Unmarshal(body, &slot)
	if err != nil {
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to read Json,Try again")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	slotCreated, err := slot.SaveSlot(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())
		json.NewEncoder(w).Encode("Status:Failure")
		json.NewEncoder(w).Encode("Failed to add slotes,Try again")
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	json.NewEncoder(w).Encode("Status:Success")
	json.NewEncoder(w).Encode("Added slotes successfully")
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, slot.ID))
	responses.JSON(w, http.StatusCreated, slotCreated)
}
