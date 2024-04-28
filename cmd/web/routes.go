package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// router устанавливает REST ресурсы с соответствующими HTTP-методами
func router() http.Handler {
	mux := chi.NewRouter()
	domain := "http://localhost:80"
	mux.Get("/swagger/*", httpSwagger.Handler(

		httpSwagger.URL(domain+"/swagger/doc.json"), //The url pointing to API definition
	))

	mux.Get("/api/v1/cars", GetCars)

	mux.Post("/api/v1/cars", PostCars)

	mux.Patch("/api/v1/cars", PatchCars)

	mux.Delete("/api/v1/cars/{carID}", DeleteCars)

	return mux
}
