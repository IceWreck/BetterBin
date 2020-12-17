package handlers

import (
	"log"
	"net/http"
)

// Home Page
func Home(w http.ResponseWriter, r *http.Request) {
	err := templateCache["home"].Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
