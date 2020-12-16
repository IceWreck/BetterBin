package handlers

import (
	"encoding/json"
	"net/http"
)

type errorJSON struct {
	Err string `json:"error"`
}

func renderError(w http.ResponseWriter, r *http.Request, err error, status int) {
	data := errorJSON{Err: err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
	return
}
