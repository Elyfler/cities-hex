package ports

import (
	"encoding/json"
	"net/http"

	"github.com/cities/internal/cities/app"
	"github.com/go-chi/chi"
)

type CityHandler interface {
	CreateCity(http.ResponseWriter, *http.Request)
	FindCityByUUID(http.ResponseWriter, *http.Request)
	FindCityByName(http.ResponseWriter, *http.Request)
	AllCities(http.ResponseWriter, *http.Request)
	DeleteCity(http.ResponseWriter, *http.Request)
}

type handler struct {
	cityService app.CityService
}

func NewHandler(cityService app.CityService) CityHandler {
	return &handler{cityService: cityService}
}

func (h *handler) CreateCity(w http.ResponseWriter, r *http.Request) {
	var c app.City
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		h.cityService.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.cityService.CreateCity(r.Context(), c)
	if err != nil {
		h.cityService.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) FindCityByUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	city, err := h.cityService.FindCityByUUID(r.Context(), uuid)
	if err != nil {
		h.cityService.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(city)
}

func (h *handler) FindCityByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	city, err := h.cityService.FindCityByName(r.Context(), name)
	if err != nil {
		h.cityService.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(city)
}
func (h *handler) AllCities(w http.ResponseWriter, r *http.Request) {

	cities, err := h.cityService.AllCities(r.Context())
	if err != nil {
		h.cityService.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func (h *handler) DeleteCity(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	err := h.cityService.DeleteCity(r.Context(), name)
	if err != nil {
		h.cityService.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}
