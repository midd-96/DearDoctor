package model

type AdminResponse struct {
	ID       int    `json:"id"`
	Username string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role"`
	Token    string `json:"token,omitempty"`
}

// user schema for user table
type UserResponse struct {
	ID               int    `json:"id"`
	First_Name       string `json:"first_name"`
	Last_Name        string `json:"last_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Phone            string `json:"phone"`
	Last_appointment int    `json:"last_appointment"`
	Token            string `json:"token"`
}

type DoctorResponse struct {
	ID         int    `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Approvel   bool   `json:"approvel"`
	Token      string `json:"token,omitempty"`
}
