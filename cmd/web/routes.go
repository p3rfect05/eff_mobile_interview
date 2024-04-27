package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// router устанавливает REST ресурсы с соответствующими HTTP-методами
func router() http.Handler {
	mux := chi.NewRouter()

	//mux.Get("/info", Info)

	mux.Get("/cars", GetCars)

	mux.Post("/cars", PostCars)

	mux.Patch("/cars", PatchCars)

	mux.Delete("/cars/{carID}", DeleteCars)
	return mux
}
