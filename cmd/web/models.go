package main

import (
	"errors"

	"gorm.io/gorm"
)

var db *gorm.DB

type Models struct {
	Car Car
}

type Owner struct {
	ID         int    `json:"-"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Car struct {
	RegNum  string `json:"regNum" gorm:"primaryKey"`
	Mark    string `json:"mark"`
	Model   string `json:"model"`
	Year    int    `json:"year,omitempty"`
	Owner   Owner  `json:"owner"`
	OwnerID int    `json:"-"`
}

func InitConn(dbPool *gorm.DB) {
	db = dbPool
}

func MigrateModels() (bool, error) {
	if !db.Migrator().HasTable(&Car{}) {
		app.InfoLog.Println("Table Car does not exist, creating...")
		if err := db.AutoMigrate(&Car{}); err != nil {
			return false, err
		}
	} else {
		app.InfoLog.Println("Table Car already exists skipping...")
	}
	if !db.Migrator().HasTable(&Owner{}) {
		app.InfoLog.Println("Table Owner does not exist, creating...")
		if err := db.AutoMigrate(&Owner{}); err != nil {
			return false, err
		}
	} else {
		app.InfoLog.Println("Table Owner already exists skipping...")
	}
	return true, nil
}

func DeleteCarByRegNum(regNum string) error {
	res := db.Delete(&Car{}, regNum)
	if res.Error != nil {
		return res.Error
	}
	return nil

}

func InsertCarInfo(car Car) (string, error) {

	res := db.Create(&car) 

	if res.Error != nil {
		return "", res.Error
	}

	return car.RegNum, nil
}

func GetCarByRegNum(regNum string) (*Car, error) {
	var car Car
	res := db.First(&car, regNum)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {

		return nil, res.Error
	}

	return &car, nil
}

func UpdateCar(car Car) error {
	res := db.Model(&car).Updates(car)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
