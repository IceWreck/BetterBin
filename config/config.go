package config

import (
	"flag"
)

// Config struct to define settings
type Config struct {
	DatabasePath string
	Port         int
}

func LoadConfig(app *Application) Config {
	var settings Config

	// Command line flags and their default values
	settings.DatabasePath = *flag.String("d", "./betterbin.sqlite", "Database Path")
	settings.Port = *flag.Int("p", 8963, "Port")

	flag.Parse()

	return settings
}
