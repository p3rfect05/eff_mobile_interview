package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func router(app *AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/info", Info)
	//mux.Get("/Cars/{CarID}", CarOne)

	mux.Get("/cars", Cars)
	mux.Post("/cars", PostCars)

	mux.Patch("/cars", PatchCars)

	mux.Delete("/cars/{carID}", DeleteCars)
	return mux
}
