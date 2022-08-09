package models

import "gorm.io/gorm"

type Confirmed struct {
	gorm.Model

	Id             int    `json:"id" gorm:"primary key"`
	Day_consult    string `json:"consulting_day" gorm:"not null"`
	Time_consult   string `json:"consulting_time" gorm:"not null"`
	Payment_mode   string `json:"payment_mode" gorm:"not null"`
	Payment_status bool   `json:"payment_status" gorm:"not null"`
	Fee            int    `json:"fee" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`
}
