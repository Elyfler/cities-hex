package app

import "errors"

var (
	ErrCityNameTooLong = errors.New("city name too long")
	ErrEmptyCityName   = errors.New("empty city name")
)
