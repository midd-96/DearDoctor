package models

import "github.com/jinzhu/gorm"

type Department struct {
	gorm.Model

	Dep_Id string `json:"dep_id" gorm:"primary_key"`
	Name   string `json:"dep_name" gorm:"not null;unique"`
}

func (d *Department) SaveDept(db *gorm.DB) (*Department, error) {
	var err error
	err = db.Debug().Create(&d).Error
	if err != nil {
		return &Department{}, err
	}
	return d, nil
}
