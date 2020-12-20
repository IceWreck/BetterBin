package handlers

import (
	"net/http"

	"github.com/IceWreck/BetterBin/logger"
)

// NewDropPage - webpage to make a new file drop
func NewDropPage(w http.ResponseWriter, r *http.Request) {
	err := templateCache["new_drop"].Execute(w, nil)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
