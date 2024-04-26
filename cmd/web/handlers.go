package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Info(w http.ResponseWriter, r *http.Request) {
	car := Car{
		RegNum: "123XXX567",
		Mark:   "lada",
		Model:  "vesta",
		Owner: Owner{
			Name:       "alex",
			Surname:    "v",
			Patronymic: "what the hell is patronymic",
		},
	}
	json.NewEncoder(w).Encode(car)

}
func GetCars(w http.ResponseWriter, r *http.Request) {
	var searchCar Car
	err := json.NewDecoder(r.Body).Decode(&searchCar)
	if err != nil {
		app.ErrorLog.Println("error while decoding car object in PatchCars")
		return
	}
	matchedCars, err := GetCarsByFilters(searchCar)
	if err != nil {
		app.ErrorLog.Println("error while getting cars by filters:", err)
		return
	}
	app.InfoLog.Printf("found %d cars matching the request\n", len(matchedCars))
}

func PostCars(w http.ResponseWriter, r *http.Request) {
	type jsonReq struct {
		RegNums []string `json:"regNums"`
	}
	var carNumbers jsonReq
	err := json.NewDecoder(r.Body).Decode(&carNumbers)
	if err != nil {
		app.ErrorLog.Println("error while decoding regNums[] in PostCars")
		return
	}

	for _, newCarID := range carNumbers.RegNums {
		newCar, err := getCarInfoByRegNum(newCarID)
		if err != nil {
			app.ErrorLog.Println("error while getting car info via API:", err)
			continue
		}
		_, err = InsertCarInfo(newCar)

		if err != nil {
			app.ErrorLog.Println("error while inserting car:", err)
		} else {
			app.InfoLog.Println("inserted car with regNum:", newCarID)
		}
	}

}

func PatchCars(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		app.ErrorLog.Println("error while decoding car object in PatchCars")
		return
	}
	err = UpdateCar(newCar)
	if err != nil {
		app.ErrorLog.Println("error while updating car:", err)
	} else {
		app.InfoLog.Println("updated car with regNum:", newCar.RegNum)
	}
}

func DeleteCars(w http.ResponseWriter, r *http.Request) {
	regNumToDelete := chi.URLParam(r, "carID")
	err := DeleteCarByRegNum(regNumToDelete)
	if err != nil {
		app.ErrorLog.Printf("error while deleting car with regNum:%s: %s\n", regNumToDelete, err)
	} else {
		app.InfoLog.Printf("deleted car with regNum:%s\n", regNumToDelete)
	}

}
