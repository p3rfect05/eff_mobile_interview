package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const externalApiURL = "https://<your-domain>/info"


// GetCarInfoByRegNum - функция для получения информации об автомобиле от внешнего API
func GetCarInfoByRegNum(regNum string) (Car, error) {
	car := Car{}
	req, _ := http.NewRequest("GET", externalApiURL, nil)
	req.URL.RawQuery = url.Values{
		"regNum": {regNum},
	}.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		app.ErrorLog.Println(err)
		return car, err
	}

	if resp.StatusCode != 200 {
		app.ErrorLog.Printf("error while getting %s?regNum=%s: %d\n", externalApiURL, regNum, resp.StatusCode)
		return car, err
	}

	err = json.NewDecoder(resp.Body).Decode(&car)

	if err != nil {
		app.ErrorLog.Println("could not decode response:", err)
		return car, err
	}

	app.InfoLog.Println("response acquired successfully")

	return car, nil

}


// TestGetCarInfoByRegNum - имитация функции GetCarInfoByRegNum
func TestGetCarInfoByRegNum(regNum string) (Car, error) {
	car := Car{
		RegNum: regNum,
		Mark:   "lada",
		Model:  "vesta",
		Owner: Owner{
			Name:       "alex",
			Surname:    "v",
			Patronymic: "what the hell is patronymic",
		},
	}

	return car, nil

}
