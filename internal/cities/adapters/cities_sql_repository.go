package adapters

import (
	"context"
	"fmt"

	"github.com/cities/internal/cities/app"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type cityModel struct {
	UUID     string `db:"uuid"`
	Name     string `db:"name"`
	PostCode int    `db:"post_code"`
}

type SQLCityRepository struct {
	db *sqlx.DB
}

func (c cityModel) toCity() app.City {
	return app.City{
		UUID:     c.UUID,
		Name:     c.Name,
		PostCode: c.PostCode,
	}
}
func NewSQLCityRepository(db *sqlx.DB) *SQLCityRepository {
	if db == nil {
		panic("missing DB")
	}

	return &SQLCityRepository{db: db}
}

func (s SQLCityRepository) CreateCity(ctx context.Context, city app.City) error {
	dbCity := cityModel{
		UUID:     city.UUID,
		Name:     city.Name,
		PostCode: city.PostCode,
	}
	query := `INSERT INTO 
					cities (name, post_code, uuid)
				VALUES
				(:name, :post_code, :uuid)`

	_, err := s.db.NamedExecContext(ctx, query, dbCity)
	if err != nil {
		return err
	}
	return nil
}

func (s SQLCityRepository) FindCityByUUID(ctx context.Context, uuid string) (app.City, error) {

	rows, err := s.db.NamedQueryContext(ctx, "SELECT * FROM cities WHERE post_code = $1", uuid)
	if err != nil {
		return app.City{}, err
	}

	var dbCity cityModel

	for rows.Next() {
		err = rows.StructScan(&dbCity)
		if err != nil {
			return app.City{}, err
		}
	}

	return dbCity.toCity(), nil
}

func (s SQLCityRepository) AllCities(ctx context.Context) ([]app.City, error) {

	rows, err := s.db.QueryxContext(ctx, "SELECT * FROM cities")
	if err != nil {
		return []app.City{}, err
	}

	var dbCities []cityModel
	err = sqlx.StructScan(rows, &dbCities)
	if err != nil {
		return []app.City{}, err
	}

	var cities []app.City

	for _, dbCity := range dbCities {
		cities = append(cities, dbCity.toCity())
	}

	return cities, nil
}

func (s SQLCityRepository) DeleteCity(ctx context.Context, uuid string) error {

	_, err := s.db.QueryxContext(ctx, `DELETE FROM cities WHERE uuid = $1`, uuid)
	if err != nil {
		return err
	}

	return nil
}

func NewPostgresConnection() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s", "postgres", "postgres", "127.0.0.1", "5432"))
}
