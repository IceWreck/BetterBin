package handlers

import (
	"net/http"
)

// Home Page
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", nil)
}
