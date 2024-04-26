package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const externalApiURL = "http://localhost/info"

func getCarInfoByRegNum(regNum string, app *AppConfig) {

	req, _ := http.NewRequest("GET", externalApiURL, nil)
	req.URL.RawQuery = url.Values{
		"regNum": {regNum},
	}.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		app.ErrorLog.Printf("error while getting %s?regNum=%s: %d\n", externalApiURL, regNum, resp.StatusCode)
		return
	}
	car := Car{}
	err = json.NewDecoder(resp.Body).Decode(&car)

	if err != nil {
		app.ErrorLog.Println("could not decode response:", err)
	}

	app.InfoLog.Println("response acquired successfully")

}
