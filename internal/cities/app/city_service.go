package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type cityRepository interface {
	CreateCity(ctx context.Context, city City) error
	FindCityByUUID(ctx context.Context, uuid string) (City, error)
	AllCities(ctx context.Context) ([]City, error)
	DeleteCity(ctx context.Context, cityUUID string) error
}

type CityService struct {
	repo cityRepository
}

func NewCityService(
	repo cityRepository,
) CityService {
	if repo == nil {
		panic("missong cityRepository")
	}

	return CityService{repo: repo}
}

func (c CityService) CreateCity(ctx context.Context, city City) error {
	// Small sanity check
	if len(city.Name) > 25 {
		return errors.New("City name too big")
	}

	city.UUID = uuid.New().String()

	return c.repo.CreateCity(ctx, city)
}

func (c CityService) FindCityByUUID(ctx context.Context, uuid string) (City, error) {
	return c.repo.FindCityByUUID(ctx, uuid)
}

func (c CityService) AllCities(ctx context.Context) ([]City, error) {
	return c.repo.AllCities(ctx)
}

func (c CityService) DeleteCity(ctx context.Context, uuid string) error {
	return c.repo.DeleteCity(ctx, uuid)
}
