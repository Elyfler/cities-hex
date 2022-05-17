package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type cityRepository interface {
	CreateCity(ctx context.Context, city City) error
	FindCityByUUID(ctx context.Context, uuid string) (City, error)
	FindCityByName(ctx context.Context, name string) (City, error)
	AllCities(ctx context.Context) ([]City, error)
	DeleteCity(ctx context.Context, name string) error
}

type CityService struct {
	Logger *zap.SugaredLogger
	repo   cityRepository
}

func NewCityService(
	repo cityRepository,
	log *zap.SugaredLogger,
) CityService {
	log.Info("Creating new city service")
	if repo == nil {
		log.Panic("missing cityRepository")
	}

	return CityService{repo: repo, Logger: log}
}

func (c CityService) CreateCity(ctx context.Context, city City) error {
	// Small sanity check
	if len(city.Name) > 25 {
		err := errors.New("City name too big")
		return err
	}
	city.UUID = uuid.New().String()

	c.Logger.Info("Creating city", zap.String("UUID", city.UUID))

	return c.repo.CreateCity(ctx, city)
}

func (c CityService) FindCityByUUID(ctx context.Context, uuid string) (City, error) {
	c.Logger.Info("Finding city", zap.String("UUID", uuid))
	return c.repo.FindCityByUUID(ctx, uuid)
}

func (c CityService) FindCityByName(ctx context.Context, name string) (City, error) {
	c.Logger.Info("Finding city", zap.String("Name", name))
	return c.repo.FindCityByName(ctx, name)
}

func (c CityService) AllCities(ctx context.Context) ([]City, error) {
	c.Logger.Info("Finding all cities")
	return c.repo.AllCities(ctx)
}

func (c CityService) DeleteCity(ctx context.Context, name string) error {
	c.Logger.Info("Deleting city", zap.String("Name", name))
	return c.repo.DeleteCity(ctx, name)
}
