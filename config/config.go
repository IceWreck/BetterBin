package config

import "flag"

// Command line flags and their default values
var (
	DatabasePath = flag.String("d", "./betterbin.sqlite", "Database Path")
	Port         = flag.Int("p", 8963, "Port")
)

func init() {
	flag.Parse()
}
