package main

import (
	"net/http"

	"github.com/cities/internal/cities/adapters"
	"github.com/cities/internal/cities/app"
	"github.com/cities/internal/cities/ports"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"moul.io/chizap"
)

func main() {
	l, _ := zap.NewDevelopment()
	defer l.Sync()
	zap := l.Sugar()
	zap.Info("Starting application")

	db, err := adapters.NewPostgresConnection()
	if err != nil {
		zap.Panic("Could not reach DB")
	}

	repo := adapters.NewSQLCityRepository(db)
	service := app.NewCityService(repo, zap)

	handler := ports.NewHandler(service)

	r := chi.NewRouter()
	r.Use(chizap.New(l, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	r.Post("/cities", handler.CreateCity)
	r.Get("/cities/{uuid}", handler.FindCityByUUID)
	r.Get("/cities", handler.AllCities)
	r.Delete("/cities/{uuid}", handler.DeleteCity)

	http.ListenAndServe("127.0.0.1:8000", r)
}
