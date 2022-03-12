package handlers

import (
	"net/http"

	"github.com/IceWreck/BetterBin/config"
)

// home Page
func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "home", nil)
	}
}
