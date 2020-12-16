package logger

import (
	"fmt"
	"log"
)

func Error(v ...interface{}) {
	log.Print("ERROR ", fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	log.Print("INFO ", fmt.Sprintln(v...))
}

func Debug(v ...interface{}) {
	log.Print("DEBUG ", fmt.Sprintln(v...))
}

func Fatal(v ...interface{}) {
	log.Fatal("FATAL ", fmt.Sprintln(v...))
}
