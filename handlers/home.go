package handlers

import "net/http"

// Home Page
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home"))
}
