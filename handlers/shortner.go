package handlers

import (
	"net/http"
)

// NewLinkForm - new shortened url using a POST request
func NewLinkForm(w http.ResponseWriter, r *http.Request) {
	longLink := r.PostFormValue("url")

	// error out if link is empty
	if len(longLink) < 1 {
		w.Write([]byte("no link"))
		return
	}

	w.Write([]byte("shortend url"))
	return
}
