package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IceWreck/BetterBin/logger"
)

type errorJSON struct {
	Err string `json:"error"`
}

type pasteJSON struct {
	ID string `json:"id"`
}

type dropJSON struct {
	ID       string `json:"id"`
	FileName string `json:"filename"`
}

func renderError(w http.ResponseWriter, r *http.Request, err error, status int) {
	data := errorJSON{Err: err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
	return
}

func renderSuccess(w http.ResponseWriter, r *http.Request, id string) {
	data := pasteJSON{ID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
}

func renderDropSuccess(w http.ResponseWriter, r *http.Request, id string, filename string) {
	data := dropJSON{ID: id, FileName: filename}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
}

func renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	err := templateCache[template].Execute(w, data)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	return
}
