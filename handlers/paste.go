package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi"
)

var errPasteBurned = errors.New("paste has been burned")
var errPasteExpired = errors.New("paste has expired")
var errPasteNotFound = errors.New("paste not found")
var errNoContent = errors.New("no content for new paste")

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
		renderError(w, r, errNoContent, http.StatusUnprocessableEntity)
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
		renderError(w, r, err, http.StatusInternalServerError)
		return
	}
	renderPaste(w, r, pasteID)
	return
}

// NewPastePage - webpage to make a new paste
func NewPastePage(w http.ResponseWriter, r *http.Request) {
	err := templateCache["new_paste"].Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// ViewPastePage - webpage to view a paste
func ViewPastePage(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err != nil {
		renderError(w, r, err, http.StatusNotFound)
		return
	}
	err = templateCache["view_paste"].Execute(w, paste)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// ViewPasteRaw - curl/wget friendly raw paste contents
func ViewPasteRaw(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err != nil {
		renderError(w, r, err, http.StatusNotFound)
		return
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
