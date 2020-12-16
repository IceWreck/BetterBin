package handlers

import (
	"errors"
	"net/http"

	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi"
)

var errPasteBurned = errors.New("paste has been burned")
var errPasteExpired = errors.New("paste has expired")
var errPasteNotFound = errors.New("paste not found")

// NewPasteForm - new paste using a POST request
func NewPasteForm(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	expiry := r.PostFormValue("expiry")
	burnStr := r.PostFormValue("burn")
	discussStr := r.PostFormValue("discuss")
	password := r.PostFormValue("password")

	// string (technically boolean) input from forms should be converted to ints
	burn := 0
	discuss := 0
	if burnStr == "1" {
		burn = 1
	}
	if discussStr == "1" {
		discuss = 1
	}

	// error out if content is empty
	if len(content) < 1 {
		w.Write([]byte("no content"))
		return
	}

	// convert expiry to a form suitable for sql
	switch expiry {
	case "year":
	case "month":
	case "week":
	case "day":
	case "hour":
	case "10min":
	case "1min":
	default:
		// never expire
		expiry = "never"
	}
	pasteID := newID(10)
	logger.Info("creating new paste", title, pasteID)
	if err := db.NewPaste(pasteID, title, content, expiry, password, burn, discuss); err != nil {
		logger.Info("could not create a new paste")
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("new paste form"))
	return
}

// NewPastePage - webpage to make a new paste
func NewPastePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("new paste page"))
}

// ViewPastePage - webpage to view a paste
func ViewPastePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("view paste page"))
}

// ViewPasteRaw - curl/wget friendly raw paste contents
func ViewPasteRaw(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(paste.Content))
}

// getPaste has common stuff required in both ViewPastePage and ViewPasteRaw
func getPaste(r *http.Request) (db.Paste, error) {
	pasteID := chi.URLParam(r, "pasteID")
	logger.Debug("fetching paste", pasteID)
	paste, err := db.GetPaste(pasteID)
	if err != nil {
		// error prolly means that it can't find paste
		return paste, errPasteNotFound
	}
	// TODO: add expired/burned checks
	return paste, nil
}
