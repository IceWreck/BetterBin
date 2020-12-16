package db

import (
	"errors"

	"github.com/IceWreck/BetterBin/logger"
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
func NewLink(id string, completeLink string) error {
	query := `INSERT INTO shortened_links (id, complete_link, created) VALUES ($1, $2, datetime('now'))`
	_, err := db.Exec(query, id, completeLink)
	if err != nil {
		logger.Error("sql spews", err)
		return errCannotCreateShortenedLink
	}
	return nil
}

// GetLink is the db operation to fetch expanded url for shortened url
func GetLink(id string) (ShortenedLink, error) {
	s := ShortenedLink{}
	err := db.Get(&s, "SELECT * FROM shortened_links WHERE id=$1", id)
	if err != nil {
		logger.Error("cannot fetch shortened link", id, err)
	}
	return s, err
}
