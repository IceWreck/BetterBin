package handlers

import (
	"encoding/json"
	"net/http"
)

type errorJSON struct {
	Err string `json:"error"`
}

type pasteJSON struct {
	ID string `json:"paste_id"`
}

func renderError(w http.ResponseWriter, r *http.Request, err error, status int) {
	data := errorJSON{Err: err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
	return
}

func renderPaste(w http.ResponseWriter, r *http.Request, id string) {
	data := pasteJSON{ID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
}
