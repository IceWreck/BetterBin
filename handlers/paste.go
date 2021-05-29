package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi"
)

var errPasteExpired = errors.New("paste has expired")
var errPasteNotFound = errors.New("paste not found")
var errNoContent = errors.New("no content for new paste")
var errPasswordRequired = errors.New("password required")

const timeLayout = "2006-01-02 15:04:05"

// NewPasteForm (API) - new paste using a POST request
func NewPasteForm(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	expiry := r.PostFormValue("expiry")
	password := r.PostFormValue("password")

	// string (technically boolean) input from forms should be converted to ints
	burn := 0

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
	case "burn":
		expiry = "1 month"
		burn = 1
	default:
		// never expire
		expiry = "999 years"
	}
	pasteID := newID(10)
	logger.Info("creating new paste", title, pasteID)
	if err := db.NewPaste(pasteID, title, content, expiry, password, burn); err != nil {
		logger.Info("could not create a new paste")
		renderError(w, r, err, http.StatusInternalServerError)
		return
	}
	renderSuccess(w, r, pasteID)
	return
}

// NewPastePage - webpage to make a new paste
func NewPastePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_paste", nil)
}

// ViewPastePage - webpage to view a paste
func ViewPastePage(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err == errPasswordRequired {
		renderTemplate(w, "password_required", nil)
		return
	} else if err == errPasteExpired {
		renderTemplate(w, "paste_expired", nil)
		return
	} else if err != nil {
		renderTemplate(w, "paste_not_found", nil)
		return
	}
	renderTemplate(w, "view_paste", paste)
}

// ViewPasteRaw - curl/wget friendly raw paste contents
func ViewPasteRaw(w http.ResponseWriter, r *http.Request) {
	paste, err := getPaste(r)
	if err != nil {
		renderError(w, r, err, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
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

	// expiry check
	expiry, err := time.Parse(timeLayout, paste.Expiry)
	if err != nil {
		logger.Error("unable to parse sql datetime", err)
	}
	if expiry.Sub(time.Now()) < 0 {
		// expired
		return paste, errPasteExpired
	}

	// password check
	if password != paste.Password {
		return paste, errPasswordRequired
	}

	// burn paste for future fetches if burn=1
	if paste.Burn == 1 {
		logger.Info("burning paste", pasteID)
		db.BurnPaste(pasteID)
	}

	// preview information from query paramater
	preview := r.FormValue("preview")
	if preview == "code" {
		paste.Preview = "code"
	} else if preview == "markdown" {
		paste.Preview = "markdown"
	}

	return paste, nil
}
