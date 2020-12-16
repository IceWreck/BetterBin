package main

import (
	"net/http"
	"time"

	"github.com/IceWreck/BetterBin/handlers"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	// TODO: cron like timer to purge expired/burned pastes and drops

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handlers.Home)
	// Paste Bin
	r.Get("/paste/new", handlers.NewPastePage)
	r.Post("/paste/new", handlers.NewPasteForm)
	r.Get("/paste/view/{pasteID}", handlers.ViewPastePage)
	r.Get("/paste/raw/{pasteID}", handlers.ViewPasteRaw)
	// Link Shortner
	r.Get("/shortner/new", handlers.Home)
	r.Post("/shortner/new", handlers.NewLinkForm)
	r.Get("/s/{linkID}", handlers.RedirectLink)
	// File Drop
	r.Get("/drop/new", handlers.Home)
	r.Post("/drop/new", handlers.Home)
	r.Get("/drop/dl/{dropID}", handlers.Home)

	fileServer(r, "/static", http.Dir("./static"))

	logger.Info("Starting at :8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		logger.Error(err)
	}
}
