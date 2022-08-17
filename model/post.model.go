package model

import (
	"time"

	"gorm.io/gorm"
)

// user schema for user table to get listed all users
type User struct {
	Id               int    `json:"user_id"`
	First_Name       string `json:"first_name"`
	Last_Name        string `json:"last_name"`
	Email            string `json:"email" gorm:"not null;unique"`
	Phone            int    `json:"phone_number"`
	Password         string `json:"password"`
	Last_appointment int    `json:"last_appointment"`
	Role             bool   `json:"role"`
}

// confirmed appointments shema for confirmed table
type Confirmed struct {
	Id             int    `json:"id" gorm:"primary key"`
	Day_consult    string `json:"consulting_day" gorm:"not null"`
	Time_consult   string `json:"consulting_time" gorm:"not null"`
	Payment_mode   string `json:"payment_mode" gorm:"not null"`
	Payment_status bool   `json:"payment_status" gorm:"not null"`
	Fee            int    `json:"fee" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`
	Doctor_id      int    `json:"doc_id"`
}

//  to add new department as departments as schema

type Departments struct {
	gorm.Model

	Dep_Id string `json:"dep_id" gorm:"primary_key"`
	Name   string `json:"dep_name" gorm:"not null;unique"`
}

//table schema to update doctor table by admin to approve and add fee

type ApproveAndFee struct {
	Approve   bool `json:"approvel"`
	Fee       int  `json:"fee"`
	Doctor_id int  `json:"doc_id"`
}

//table schema to mark availability by doctor

type Slotes struct {
	gorm.Model

	Id            int    `json:"id" gorm:"primary_key"`
	Docter_id     int    `json:"docter_id" gorm:"not null"`
	Available_day string `json:"available_day" gorm:"not null"`
	Time_from     string `json:"staring_time" gorm:"not null"`
	Time_to       string `json:"ending_time" gorm:"not null"`
	Status        bool   `json:"booked" gorm:"default:false"`
}

//table schema to signup by doctors

type Doctor struct {
	gorm.Model

	Id         int    `json:"id" gorm:"primary_key"`
	First_name string `json:"first_name" gorm:"not null"`
	Last_name  string `json:"last_name" gorm:"not null"`
	Email      string `json:"email" gorm:"unique;not null" valid:"email"`
	Phone      string `json:"phone" gorm:"unique;not null"`
	Password   string `json:"password" gorm:"not null" valid:"length(6|20)"`
	//Dep_code       string `json:"dep_code"`
	Department     string `json:"department" gorm:"not null"`
	Specialization string `json:"specialization"`
	Approvel       bool   `json:"approvel" gorm:"default:false"`
	Fee            int    `json:"fee"`
}

//table schema for user address
type Address struct {
	User_id     int    `json:"user_id,omitempty"`
	AddressType string `json:"address_type"`
	HouseName   string `json:"house_name"`
	StreetName  string `json:"stree_name"`
	Landmark    string `json:"landmark"`
	District    string `json:"district"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PinCode     int    `json:"pincode"`
	Created_At  time.Time
}

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Filter struct {
	Category    []string `json:"category"`
	Brand       []string `json:"brand"`
	Color       []string `json:"color"`
	Processor   []string `json:"processor"`
	Name        []string `json:"name"`
	ProductCode []string `json:"product_code"`
}