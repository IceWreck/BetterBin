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
var errNoContent = errors.New("no content for new paste")
var errPasswordRequired = errors.New("password required")

// NewPasteForm (API) - new paste using a POST request
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
	// never place user input directly into the query
	switch expiry {
	case "year":
		expiry = "1 years"
	case "month":
		expiry = "1 month"
	case "week":
		expiry = "7 days"
	case "day":
		expiry = "1 days"
	case "hour":
		expiry = "1 hours"
	case "10min":
		expiry = "10 minutes"
	case "1min":
		expiry = "1 minutes"
	default:
		// never expire
		expiry = "999 years"
	}
	pasteID := newID(10)
	logger.Info("creating new paste", title, pasteID)
	if err := db.NewPaste(pasteID, title, content, expiry, password, burn, discuss); err != nil {
		logger.Info("could not create a new paste")
		renderError(w, r, err, http.StatusInternalServerError)
		return
	}
	renderSuccess(w, r, pasteID)
	return
}

// NewPastePage - webpage to make a new paste
func NewPastePage(w http.ResponseWriter, r *http.Request) {
	err := templateCache["new_paste"].Execute(w, nil)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// ViewPastePage - webpage to view a paste
func ViewPastePage(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err == errPasswordRequired {
		renderTemplate(w, "password_required", nil)
		return
	} else if err != nil {
		renderError(w, r, err, http.StatusNotFound)
		return
	}
	err = templateCache["view_paste"].Execute(w, paste)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// ViewPasteRaw - curl/wget friendly raw paste contents
func ViewPasteRaw(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err == errPasswordRequired {
		renderError(w, r, err, http.StatusNotFound)
		return
	} else if err != nil {
		renderError(w, r, err, http.StatusNotFound)
		return
	}
	w.Write([]byte(paste.Content))
}

// getPaste has common stuff required in both ViewPastePage and ViewPasteRaw
func getPaste(r *http.Request) (db.Paste, error) {
	pasteID := chi.URLParam(r, "pasteID")
	// password can be either URL query parameters or POST/PUT values
	password := r.FormValue("password")

	logger.Debug("fetching paste", pasteID)
	paste, err := db.GetPaste(pasteID)
	if err != nil {
		// error prolly means that it can't find paste
		return paste, errPasteNotFound
	}
	// TODO: add expired/burned checks

	// password check
	if password != paste.Password {
		return paste, errPasswordRequired
	}

	return paste, nil
}
