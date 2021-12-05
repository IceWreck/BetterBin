package db

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/IceWreck/BetterBin/logger"
)

var errCannotCreatePaste = errors.New("cannot create a new paste")

// Table name is 'pastes'

// Paste symbolizes a single paste
type Paste struct {
	ID              string `db:"id"`
	Title           string `db:"title"`
	Content         string `db:"content"`
	Password        string `db:"password"`
	Preview         string // One of `markdown`, `code` and `plain`.
	PreviewLanguage string // Code preview language extension.
	ContentHTML     template.HTML
	Expiry          string `db:"expiry"`
	Created         string `db:"created"`
	Burn            int    `db:"burn"`
}

// NewPaste is the db operation to create a new paste in database
func NewPaste(id string, title string, content string, expiry string, password string, burn int) error {
	query := `INSERT INTO pastes (id, title, content, password, expiry, created, burn) 
	VALUES ($1, $2, $3, $4, %s, datetime('now'), $5)`
	query = fmt.Sprintf(query, "datetime('now', '+"+expiry+"')")
	_, err := db.Exec(query, id, title, content, password, burn)
	if err != nil {
		logger.Error(err)
		return errCannotCreatePaste
	}
	return nil
}

// GetPaste is the db operation to fetch a single paste
func GetPaste(id string) (Paste, error) {
	p := Paste{}
	err := db.Get(&p, "SELECT * FROM pastes WHERE id=$1", id)
	if err != nil {
		logger.Error("cannot fetch paste", id, err)
	}
	return p, err
}

// BurnPaste removes a paste with given ID
func BurnPaste(id string) error {
	query := `DELETE FROM pastes WHERE id=$1 AND burn=1`
	_, err := db.Exec(query, id)
	if err != nil {
		logger.Error("cannot burn paste", err)
		return err
	}
	return nil
}
