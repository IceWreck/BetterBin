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
	flag.StringVar(&settings.DatabasePath, "d", "./betterbin.sqlite", "Database Path")
	flag.IntVar(&settings.Port, "p", 8963, "Port")

	flag.Parse()

	return settings
}
