package db

import (
	"github.com/IceWreck/BetterBin/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Open("sqlite3", "./betterbin.sqlite")
	err = db.Ping()
	if err != nil {
		logger.Fatal("cannot connect to db")
	} else {
		logger.Info("connected to db")
	}
}
