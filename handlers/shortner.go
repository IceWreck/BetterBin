package handlers

import (
	"errors"
	"net/http"

	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi"
)

var errInvalidLink = errors.New("invalid link")
var errLinkNotFound = errors.New("expanded URL for this id not found")

// NewLinkForm - new shortened url using a POST request
func NewLinkForm(w http.ResponseWriter, r *http.Request) {
	longLink := r.PostFormValue("url")

	// error out if link is empty
	if len(longLink) < 1 {
		renderError(w, r, errInvalidLink, http.StatusUnprocessableEntity)
		return
	}
	// TODO: regex to validate input link

	linkID := newID(10)
	if err := db.NewLink(linkID, longLink); err != nil {
		logger.Info(err)
		renderError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(linkID))
	return
}

// RedirectLink - redirect to complete link
func RedirectLink(w http.ResponseWriter, r *http.Request) {
	linkID := chi.URLParam(r, "linkID")
	logger.Debug("fetching link", linkID)
	slink, err := db.GetLink(linkID)
	if err != nil {
		renderError(w, r, errLinkNotFound, http.StatusNotFound)
		return
	}
	http.Redirect(w, r, slink.CompleteLink, http.StatusFound)
}
