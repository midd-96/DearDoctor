package models

import "github.com/jinzhu/gorm"

type Slote struct {
	gorm.Model

	Id            int    `json:"id" gorm:"primary_key"`
	Docter_id     int    `json:"docter_id" gorm:"not null"`
	Available_day string `json:"available_day" gorm:"not null"`
	Time_from     string `json:"staring_time" gorm:"not null"`
	Time_to       string `json:"ending_time" gorm:"not null"`
	Status        bool   `json:"booked" gorm:"default:false"`
}

func (s *Slote) SaveSlot(db *gorm.DB) (*Slote, error) {
	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return &Slote{}, err
	}
	return s, nil
}
