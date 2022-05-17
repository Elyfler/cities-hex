package app_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/adrianbrad/psqldocker"
	"github.com/adrianbrad/psqltest"
	"github.com/cities/internal/cities/adapters"
	"github.com/cities/internal/cities/app"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {

	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	c, err := psqldocker.NewContainer(
		usr, password, dbName,
		psqldocker.WithContainerName(containerName),
		psqldocker.WithSql(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE cities
		(
			name VARCHAR(25) NOT NULL,
			post_code INT NOT NULL,
			uuid uuid DEFAULT uuid_generate_v4 (),
			PRIMARY KEY (UUID)
		);`),
	)
	if err != nil {
		log.Fatalf("err while creating new psql container: %s", err)
	}

	var ret int

	defer func() {

		err = c.Close()
		if err != nil {
			log.Printf("err while tearing down db container: %s", err)
		}

		os.Exit(ret)
	}()

	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable ",
		usr,
		password,
		dbName,
		c.Port(),
	)

	psqltest.Register(dsn)

	ret = m.Run()
}

func setupService() app.CityService {
	db, _ := adapters.NewPostgresConnection()
	repo := adapters.NewSQLCityRepository(db)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	return app.NewCityService(repo, logger)
}

func TestCreateCity(t *testing.T) {
	ctx := context.Background()
	svc := setupService()

	t.Run("A city saved should be found by its name", func(t *testing.T) {
		city := app.City{Name: "city", PostCode: 12345}
		err := svc.CreateCity(ctx, city)
		if err != nil {
			t.Error(err)
		}

		cityRet, err := svc.FindCityByName(ctx, city.Name)
		if err != nil {
			t.Error(err)
		}
		if city.Name != cityRet.Name {
			t.Errorf("got: %v, expected: %v", cityRet.Name, city.Name)
		}
		if city.PostCode != cityRet.PostCode {
			t.Errorf("got: %v, expected: %v", cityRet.PostCode, city.PostCode)
		}

	})
}

func TestDeleteCity(t *testing.T) {
	ctx := context.Background()
	svc := setupService()

	t.Run("A city saved could be deleted", func(t *testing.T) {
		city := app.City{Name: "city", PostCode: 12345}
		err := svc.CreateCity(ctx, city)
		if err != nil {
			t.Error(err)
		}

		err = svc.DeleteCity(ctx, city.Name)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

	})
}
