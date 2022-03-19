package db

import (
	"errors"

	"github.com/IceWreck/BetterBin/config"
)

var errCannotCreateShortenedLink = errors.New("cannot create shortened url")

// Table name is 'shortened_links'

// ShortenedLink is struct to resemble 'shortened_links' table
type ShortenedLink struct {
	ID           string `db:"id"`
	CompleteLink string `db:"complete_link"`
	Created      string `db:"created"`
}

// NewLink is the db operation to create a new shortened url in database
func NewLink(app *config.Application, id string, completeLink string) error {
	query := `INSERT INTO shortened_links (id, complete_link, created) VALUES ($1, $2, datetime('now'))`
	_, err := app.DB.Exec(query, id, completeLink)
	if err != nil {
		app.Logger.Error().Err(err).Msg("sql spew")
		return errCannotCreateShortenedLink
	}
	return nil
}

// GetLink is the db operation to fetch expanded url for shortened url
func GetLink(app *config.Application, id string) (ShortenedLink, error) {
	s := ShortenedLink{}
	err := app.DB.Get(&s, "SELECT * FROM shortened_links WHERE id=$1", id)
	if err != nil {
		app.Logger.Error().Str("id", id).Err(err).Msg("cannot fetch shortened link")
	}
	return s, err
}

// LinkIDExists returns true when id exists
func LinkIDExists(app *config.Application, id string) bool {
	s := 0
	// returns error when id does not exist
	err := app.DB.Get(&s, "SELECT 1 FROM shortened_links WHERE id=$1 LIMIT 1", id)
	if err == nil {
		app.Logger.Debug().Str("id", id).Msg("id already exists")
		return true
	}
	return false
}
