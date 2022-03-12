package db

import (
	"github.com/IceWreck/BetterBin/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

func ConnectDB(app *config.Application) {
	db, err := sqlx.Open("sqlite3", app.Config.DatabasePath)
	err = db.Ping()
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Cannot connect to db")
	} else {
		app.DB = db
		app.Logger.Info().Msg("Connected to db")
	}
}
