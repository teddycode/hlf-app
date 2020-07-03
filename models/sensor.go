package models

import "github.com/jinzhu/gorm"

type Sensor struct {
	ID       string `json:"id"`
	Province string `json:"province"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Code     string `json:"code"`
	TypeID   string `json:"type_id"`
	Desc     string `json:"desc"`
	PicID    string `json:"pic_id"`
}

// get all sensors
func FindAllSensors() (*[]Sensor, error) {
	var sensors []Sensor
	err := db.Select("code,desc").Find(&sensors).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &sensors, err
}

// get each sensor's data number
func CountSensorsData(id Sensor) (int64, error) {
	var cnt int64
	err := db.Model(&Sensor{}).Where("Code = ?", id).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return cnt, err
}
