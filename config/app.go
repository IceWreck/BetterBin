package config

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// Application struct to hold the dependencies for our application.
type Application struct {
	Config Config
	Logger zerolog.Logger

	DB *sqlx.DB
}
