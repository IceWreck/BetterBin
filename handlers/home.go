package handlers

import (
	"net/http"

	"github.com/IceWreck/BetterBin/logger"
)

// Home Page
func Home(w http.ResponseWriter, r *http.Request) {
	err := templateCache["home"].Execute(w, nil)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
