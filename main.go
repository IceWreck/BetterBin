package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IceWreck/BetterBin/config"
	"github.com/IceWreck/BetterBin/handlers"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	// TODO: cron like timer to purge expired pastes
	// until then server admin can manually view database to see expired pastes

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/paste/new", http.StatusTemporaryRedirect)
	})
	r.Get("/home", handlers.Home)
	// Paste Bin
	r.Get("/paste/new", handlers.NewPastePage)
	r.Post("/paste/new", handlers.NewPasteForm)
	r.Get("/paste/view/{pasteID}", handlers.ViewPastePage)
	r.Get("/paste/raw/{pasteID}", handlers.ViewPasteRaw)
	r.Post("/paste/view/{pasteID}", handlers.ViewPastePage) // required for password
	r.Post("/paste/raw/{pasteID}", handlers.ViewPasteRaw)   // required for password
	// Link Shortner
	r.Get("/shortner/new", handlers.NewLinkPage)
	r.Post("/shortner/new", handlers.NewLinkForm)
	r.Get("/s/{linkID}", handlers.RedirectLink)
	// File Drop
	r.Get("/drop/new", handlers.NewDropPage)
	r.Post("/drop/new", handlers.UploadFile)
	r.Get("/drop/dl/{dropID}", handlers.ViewDrop)
	handlers.FileServer(r, "/drops", http.Dir("./drops")) // download drops
	// Static Files (CSS/JS/Images)
	handlers.FileServer(r, "/static", http.Dir("./static"))

	logger.Info("Starting at port", *config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *config.Port), r)
	if err != nil {
		logger.Error(err)
	}
}
