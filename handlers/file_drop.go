package handlers

import (
	"net/http"
)

// NewDropPage - webpage to make a new file drop
func NewDropPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_drop", nil)
}

// UploadFile - upload a new file
func UploadFile(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_drop", nil)
}
