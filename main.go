package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IceWreck/BetterBin/config"
	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/handlers"
	"github.com/rs/zerolog"
)

func main() {

	app := &config.Application{
		Logger: zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC822,
			},
		).With().Timestamp().Logger(),
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	app.Config = config.LoadConfig(app)

	db.ConnectDB(app)

	// TODO: cron like timer to purge expired pastes
	// until then server admin can manually view database to see expired pastes

	// Initialize Router
	r := handlers.Routes(app)

	app.Logger.Info().Int("port", app.Config.Port).Msg("Starting")
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), r)
	if err != nil {
		app.Logger.Error().Err(err).Msg("")
	}
}
