package models

import (
	"gorm.io/gorm"
)

type Slotes struct {
	gorm.Model

	Id            int    `json:"id" gorm:"primary_key"`
	Docter_id     int    `json:"docter_id" gorm:"not null"`
	Available_day string `json:"available_day" gorm:"not null"`
	Time_from     string `json:"staring_time" gorm:"not null"`
	Time_to       string `json:"ending_time" gorm:"not null"`
	Status        bool   `json:"booked" gorm:"default:false"`
}
