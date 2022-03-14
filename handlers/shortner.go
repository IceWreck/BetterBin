package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/IceWreck/BetterBin/config"
	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi/v5"
)

var errInvalidLink = errors.New("invalid link")
var errLinkNotFound = errors.New("expanded URL for this id not found")
var errInvalidPreferredID = errors.New("preferred id is invalid")
var errPreferredIDExists = errors.New("preferred id is already in use")

// newLinkForm (API) - new shortened url using a POST request
func newLinkForm(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		preferredID := strings.Replace(r.PostFormValue("id"), " ", "", -1)
		longLink := r.PostFormValue("url")

		// error out if link is empty
		if len(longLink) < 1 {
			renderError(w, r, errInvalidLink, http.StatusUnprocessableEntity)
			return
		}

		// validate input link
		_, err := url.ParseRequestURI(longLink)
		if err != nil {
			renderError(w, r, errInvalidLink, http.StatusUnprocessableEntity)
			return
		}

		linkID := newID(10)

		if preferredID != "" {
			// error out if preferred id doesn't pass validation
			if len(preferredID) < 5 || len(preferredID) > 30 {
				renderError(w, r, errInvalidPreferredID, http.StatusUnprocessableEntity)
				return
			}

			if db.LinkIDExists(app, preferredID) {
				renderError(w, r, errPreferredIDExists, http.StatusUnprocessableEntity)
				return
			}

			linkID = preferredID
		}

		if err = db.NewLink(app, linkID, longLink); err != nil {
			logger.Info(err)
			renderError(w, r, err, http.StatusInternalServerError)
			return
		}
		renderSuccess(w, r, linkID)
	}
}

// newLinkPage - webpage to make a new shortened link
func newLinkPage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "new_link", nil)
	}
}

// redirectLink - redirect to complete link
func redirectLink(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkID := chi.URLParam(r, "linkID")
		logger.Debug("fetching link", linkID)
		slink, err := db.GetLink(app, linkID)
		if err != nil {
			renderError(w, r, errLinkNotFound, http.StatusNotFound)
			return
		}
		http.Redirect(w, r, slink.CompleteLink, http.StatusFound)
	}
}
