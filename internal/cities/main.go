package main

import (
	"log"
	"net/http"

	"github.com/cities/internal/cities/adapters"
	"github.com/cities/internal/cities/app"
	"github.com/cities/internal/cities/ports"
	"github.com/go-chi/chi"
)

func main() {

	db, err := adapters.NewPostgresConnection()
	if err != nil {
		log.Panic("Could not reach DB")
	}

	repo := adapters.NewSQLCityRepository(db)
	service := app.NewCityService(repo)

	handler := ports.NewHandler(service)

	r := chi.NewRouter()

	r.Post("/cities", handler.CreateCity)
	r.Get("/cities/{uuid}", handler.FindCityByUUID)
	r.Get("/cities", handler.AllCities)
	r.Delete("/cities/{uuid}", handler.DeleteCity)

	http.ListenAndServe("127.0.0.1:8000", r)
}
