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

type Appointments struct {
	Day_consult    string `json:"consulting_day" gorm:"not null"`
	Time_consult   string `json:"consulting_time" gorm:"not null"`
	Payment_mode   string `json:"payment_mode" gorm:"not null"`
	Payment_status bool   `json:"payment_status" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`
}
