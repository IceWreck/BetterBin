package db

import (
	"errors"

	"github.com/IceWreck/BetterBin/config"
	"github.com/IceWreck/BetterBin/logger"
)

var errCannotCreateDrop = errors.New("cannot create file drop")

// Table name is 'file_drops'

// FileDrop is struct to resemble 'file_drops' table
type FileDrop struct {
	ID       string `db:"id"`
	Title    string `db:"title"`
	Created  string `db:"created"`
	FileName string `db:"filename"`
	Preview  string
}

// NewDrop is the db operation to create a new file drop in database
func NewDrop(app *config.Application, id string, title string, filename string) error {
	query := `INSERT INTO file_drops (id, title, created, filename) VALUES ($1, $2, datetime('now'), $3)`
	_, err := app.DB.Exec(query, id, title, filename)
	if err != nil {
		logger.Error("sql spews", err)
		return errCannotCreateDrop
	}
	return nil
}

// GetDrop is the db operation to fetch details for the file drop
func GetDrop(app *config.Application, id string) (FileDrop, error) {
	d := FileDrop{}
	err := app.DB.Get(&d, "SELECT * FROM file_drops WHERE id=$1", id)
	if err != nil {
		logger.Error("cannot fetch file drop", id, err)
	}
	return d, err
}
