package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/p3rfect05/effecive_mobile/cmd/web/docs"
)

// @title			Cars Swagger API
// @version		1.0
// @description	Это описание API для работы с каталогом автомобилей
// @host			localhost:80
// @BasePath		/api/v1

// infoJson - ответ при успешном выполнении запроса/операции
type infoJson struct {
	InfoMessage string `json:"info_message"`
}

// errorJson - ответ при неуспешном выполнении запроса/операции

type errorJson struct {
	ErrorMessage string `json:"error_message"`
}

type updatedCarJson struct {
	Car   Car `json:"car"`
	Total int `json:"total"`
}

type postCarReq struct {
	RegNums []string `json:"reg_nums"`
}

type insertedCarsJson struct {
	Cars  []Car `json:"cars"`
	Total int   `json:"total"`
}

type foundCarsJsons struct {
	Total       int   `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	CarsPerPage int   `json:"per_page"`
	Cars        []Car `json:"cars"`
}

// writeError возвращает json ответ с описанием ошибки
func writeError(w http.ResponseWriter, err error, status_code int) {
	w.WriteHeader(status_code)
	w.Header().Set("Content-Type", "application/json")
	error_json, err := json.MarshalIndent(errorJson{
		ErrorMessage: err.Error(),
	}, "", "\t")
	if err != nil {
		app.ErrorLog.Println("error while writing error json")
		return
	}
	w.Write(error_json)

}

// writeInfo возвращает json ответ с описанием успешно выполненной операции
func writeInfo(w http.ResponseWriter, info_msg string, status_code int) {
	w.WriteHeader(status_code)
	w.Header().Set("Content-Type", "application/json")
	info_json, err := json.MarshalIndent(infoJson{
		InfoMessage: info_msg,
	}, "", "\t")
	if err != nil {
		app.ErrorLog.Println("error while writing info json")
		return
	}
	w.Write(info_json)
}

// Info - функция, имитирующая работу внешнего API
// Может понадобиться в процессе разработки, но на данный момент не используется
// func Info(w http.ResponseWriter, r *http.Request) {
// 	regNum := r.URL.Query().Get("regNum")
// 	car := Car{
// 		RegNum: regNum,
// 		Mark:   "lada",
// 		Model:  "vesta",
// 		Owner: Owner{
// 			Name:       "alex",
// 			Surname:    "v",
// 			Patronymic: "what the hell is patronymic",
// 		},
// 	}
// 	json.NewEncoder(w).Encode(car)

// }

// GetCars возвращает список объектов Car, которые совпадают с переданнами URL-параметрами
// возвращает определенную страницу(по умолчанию первую)
// с количеством объектов на каждой странице PerPage(по умолчанию 10)
// возвращает максимум limit объектов(по умолчанию 10)
// PatchCars godoc
//
//	@Summary		Get lists of request-matched car
//	@Description	GetCars возвращает список объектов Car, которые совпадают с переданнами URL-параметрами
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			reg_num				query		string	false	"Рег. номер"
//	@Param			mark				query		string	false	"Марка машины"
//	@Param			model				query		string	false	"Модель машины"
//	@Param			year				query		int		false	"Год производства"
//	@Param			page				query		int		false	"Страница (по умолчанию 1)"
//	@Param			limit				query		int		false	"Макс. записей с указанной страницы"
//	@Param			owner_name			query		string	false	"Имя владельца"
//	@Param			owner_surname		query		string	false	"Фамилия владельца"
//	@Param			owner_patronymic	query		string	false	"Отчество владельца"
//	@Success		200					{object}	foundCarsJsons
//	@Failure		400					{object}	errorJson
//	@Failure		500					{object}	errorJson
//	@Router			/cars      [get]
func GetCars(w http.ResponseWriter, r *http.Request) {
	PageLimit := 10
	Page := 1
	PerPage := 10
	URLPageLimit := r.URL.Query().Get("limit")

	if len(URLPageLimit) > 0 {
		var intPageLimit int
		intPageLimit, err := strconv.Atoi(URLPageLimit)
		if err != nil {
			app.ErrorLog.Println("error while getting limit value:", err.Error())
			writeError(w, fmt.Errorf("invalid limit value on a page: must be positive integer"), http.StatusBadRequest)
			return
		}
		if intPageLimit <= 0 {
			app.ErrorLog.Println("invalid limit value: must be positive, got:", intPageLimit)
			writeError(w, fmt.Errorf("invalid limit value on a page: must be positive integer"), http.StatusBadRequest)
			return
		}
		PageLimit = intPageLimit
	}

	URLPage := r.URL.Query().Get("page")

	if len(URLPage) > 0 {
		var intPage int
		intPage, err := strconv.Atoi(URLPage)
		if err != nil {
			app.ErrorLog.Println("error while getting page value:", err.Error())
			writeError(w, fmt.Errorf("invalid page value: must be positive integer"), http.StatusBadRequest)
			return
		}
		if intPage <= 0 {
			app.ErrorLog.Println("invalid page value: must be positive, got:", intPage)
			writeError(w, fmt.Errorf("invalid page value: must be positive integer"), http.StatusBadRequest)
			return
		}
		Page = intPage
	}

	var year int
	regNum := r.URL.Query().Get("reg_num")
	URLyear := r.URL.Query().Get("year")
	if len(URLyear) > 0 {
		intYear, err := strconv.Atoi(URLyear)
		if err != nil {
			app.ErrorLog.Println("error while getting year value:", err.Error())
			writeError(w, fmt.Errorf("invalid year value: must be positive integer"), http.StatusBadRequest)
			return
		}
		if intYear <= 0 {
			app.ErrorLog.Println("invalid year value: must be positive, got:", intYear)
			writeError(w, fmt.Errorf("invalid year value: must be positive integer"), http.StatusBadRequest)
			return
		}
		year = intYear
	}
	model := r.URL.Query().Get("model")
	mark := r.URL.Query().Get("mark")
	ownerName := r.URL.Query().Get("owner_name")
	ownerSurname := r.URL.Query().Get("owner_surname")
	ownerPatr := r.URL.Query().Get("owner_patronymic")
	searchCar := Car{
		RegNum: regNum,
		Mark:   mark,
		Model:  model,
		Year:   year,
		Owner: Owner{
			Name:       ownerName,
			Surname:    ownerSurname,
			Patronymic: ownerPatr,
		},
	}
	// dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()
	// err := dec.Decode(&searchCar)
	// if err != nil {
	// 	app.ErrorLog.Println("error while decoding car object in GetCars")
	// 	writeError(w, fmt.Errorf("error while decoding car object in GetCars"), http.StatusBadRequest)
	// 	return
	// }
	matchedCars, err := GetCarsByFilters(searchCar)
	if err != nil {
		app.ErrorLog.Println("error while getting cars by filters:", err)
		writeError(w, fmt.Errorf("error while getting cars by filters"), http.StatusInternalServerError)
		return
	}
	app.InfoLog.Printf("found %d cars matching the request\n", len(matchedCars))

	totalPages := int(math.Ceil(float64(len(matchedCars)) / float64(PerPage)))

	Page = min(Page, totalPages) // чтобы не отображать несуществующую страницу
	// например если totalPages=3, а пользователь дал Page=10, то Page = 3

	startIndex := (Page - 1) * PerPage
	endIndex := min(startIndex+PageLimit-1, Page*PerPage-1, len(matchedCars)-1)

	var resCars []Car
	if len(matchedCars) > 0 {
		resCars = matchedCars[startIndex : endIndex+1]
	} else {
		resCars = make([]Car, 0)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	res_json, err := json.MarshalIndent(foundCarsJsons{
		Total: len(resCars),
		Cars:  resCars,
		Page:  Page,
		Limit: endIndex - startIndex + 1, // фактический limit
	}, "", "\t")
	if err != nil {
		app.ErrorLog.Println("error while writing info json")
		return
	}

	w.Write(res_json)
}

// PostCars godoc
//
//	@Summary		Creates a car
//	@Description	PostCars добавляет объекты Car, с номерами, указанными в списке поля reg_nums
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			reg_nums	body		postCarReq	true	"Рег. номер"
//	@Success		200			{object}	insertedCarsJson
//	@Failure		400			{object}	errorJson
//	@Failure		500			{object}	errorJson
//	@Router			/api/v1/cars      [post]
func PostCars(w http.ResponseWriter, r *http.Request) {

	carNumbers := postCarReq{
		RegNums: make([]string, 0),
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&carNumbers)
	if err != nil {
		app.ErrorLog.Printf("error while decoding regNums[] in PostCars:%s\n", err.Error())
		writeError(w, fmt.Errorf("error while decoding regNums in PostCars"), http.StatusBadRequest)
		return
	}
	if len(carNumbers.RegNums) == 0 {
		app.ErrorLog.Printf("empty or not existing regNums")
		writeError(w, fmt.Errorf("empty or not existing regNums"), http.StatusBadRequest)
		return
	}

	if len(carNumbers.RegNums) > 30 {
		app.ErrorLog.Printf("length of regNums should not exceed 30")
		writeError(w, fmt.Errorf("length of regNums should not exceed 30"), http.StatusBadRequest)
		return
	}
	insertedCars := make([]Car, 0, 30)
	for _, newCarID := range carNumbers.RegNums {
		// ЗАМЕНИТЬ НА GetCarInfoByRegNum(newCarID)
		// В ТЗ не указан домен, поэтому здесь используется затычка, выдающая одинаковые значения
		newCar, err := TestGetCarInfoByRegNum(newCarID)
		if err != nil {
			app.ErrorLog.Printf("error while getting car info via API:%s, number of inserted cars: %d\n", err.Error(), len(insertedCars))
			writeError(w, fmt.Errorf("error while getting car info via API with regNum:%s, number of inserted cars: %d",
				newCarID, len(insertedCars)), http.StatusInternalServerError)
			return
		}
		regNum, err := InsertCarInfo(newCar)

		if err != nil {
			app.ErrorLog.Printf("error while inserting car with regNum:%s: %s\n", newCarID, err.Error())
			writeError(w, fmt.Errorf(
				"error while inserting car with regNum:%s, number of inserted cars:%d", newCarID, len(insertedCars),
			), http.StatusInternalServerError)
			return
		}
		insertedCar, err := GetCarByRegNum(regNum)
		if err != nil {
			app.ErrorLog.Printf("error while gettinginserted car with regNum:%s: %s\n", newCarID, err.Error())
			writeError(w, fmt.Errorf(
				"error while getting inserted car with regNum:%s, number of inserted cars:%d", newCarID, len(insertedCars),
			), http.StatusInternalServerError)
			return
		}
		insertedCars = append(insertedCars, *insertedCar)
	}
	app.InfoLog.Printf("successfully inserted %d cars\n", len(insertedCars))
	//writeInfo(w, fmt.Sprintf("successfully inserted %d cars", insertedCars), http.StatusOK)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	res_json, err := json.MarshalIndent(insertedCarsJson{
		Total: len(insertedCars),
		Cars:  insertedCars,
	}, "", "\t")
	if err != nil {
		app.ErrorLog.Println("error while writing info json in PostCars")
		return
	}

	w.Write(res_json)
}

// PatchCars godoc
//
//	@Summary		Updates a car
//	@Description	PatchCars изменяет указанные поля объекта Car по переданному номеру regNum
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			reg_num	body		int		true	"Рег. номер"
//	@Param			mark	body		string	false	"Марка машины"
//	@Param			model	body		string	false	"Модель машины"
//	@Param			year	body		int		false	"Год производства"
//	@Success		200		{object}	updatedCarJson
//	@Failure		400		{object}	errorJson
//	@Failure		500		{object}	errorJson
//	@Router			/api/v1/cars     [patch]
func PatchCars(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&newCar)
	if err != nil {
		app.ErrorLog.Println("error while decoding car object in PatchCars")
		writeError(w, fmt.Errorf("error while decoding car object with regNum:%s",
			newCar.RegNum), http.StatusBadRequest)
		return
	}
	newCar.Owner = Owner{} // не допускаем обновление владельца, оставляя его пустым
	if newCar.RegNum == "" {
		app.ErrorLog.Printf("empty regNum in PatchCars")
		writeError(w, fmt.Errorf("empty regNum not allowed while updating"), http.StatusBadRequest)
		return
	}
	err = UpdateCar(newCar) // обновляюся только поля с non-zero значением (отличным от 0, "", nil)
	if err != nil {
		app.ErrorLog.Printf("error while updating car with regNum:%s: %s", newCar.RegNum, err.Error())
		writeError(w, fmt.Errorf("error while updating car with regNum:%s",
			newCar.RegNum), http.StatusInternalServerError)
		return
	} else {
		app.InfoLog.Println("updated car with regNum:", newCar.RegNum)
	}

	//writeInfo(w, fmt.Sprintf("updated car with regNum:%s", newCar.RegNum), http.StatusOK)
	updatedCar, err := GetCarByRegNum(newCar.RegNum)
	if err != nil {
		app.ErrorLog.Printf("error while getting updated car with regNum:%s: %s\n", newCar.RegNum, err.Error())
		writeError(w, fmt.Errorf(
			"error while getting updated car with regNum:%s: ", newCar.RegNum,
		), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	res_json, err := json.MarshalIndent(updatedCarJson{
		Car: *updatedCar,
	}, "", "\t")
	if err != nil {
		app.ErrorLog.Println("error while writing info json in PatchCars")
		return
	}

	w.Write(res_json)
}

// DeleteCars godoc
//
//	@Summary		Delete a car
//	@Description	DeleteCars удаляет объект Car с номером regNum, переданным в URL params
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			reg_num	path		int	true	"Registration number"
//	@Success		200		{object}	infoJson
//	@Failure		400		{object}	errorJson
//	@Failure		500		{object}	errorJson
//	@Router			/api/v1/cars/{carID}      [delete]
func DeleteCars(w http.ResponseWriter, r *http.Request) {
	regNumToDelete := chi.URLParam(r, "carID")
	if len(regNumToDelete) == 0 {
		app.ErrorLog.Printf("empty regNum in PostCars")
		writeError(w, fmt.Errorf("empty regNum not allowed while deleting"), http.StatusBadRequest)
		return
	}
	res, err := DeleteCarByRegNum(regNumToDelete)
	if err != nil {
		app.ErrorLog.Printf("error while deleting car with regNum:%s: %s\n", regNumToDelete, err.Error())
		writeError(w, fmt.Errorf("error while deleting car with regNum:%s", regNumToDelete), http.StatusInternalServerError)
		return
	} else {
		app.InfoLog.Printf("deleted car with regNum:%s\n", regNumToDelete)
	}
	writeInfo(w, fmt.Sprintf("deleted %d casr with regNum:%s", res, regNumToDelete), http.StatusOK)

}
