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

	// TODO: cron like timer to purge expired pastes

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
	r.Post("/paste/view/{pasteID}", handlers.ViewPastePage) // required for password
	r.Post("/paste/raw/{pasteID}", handlers.ViewPasteRaw)   // required for password
	// Link Shortner
	r.Get("/shortner/new", handlers.NewLinkPage)
	r.Post("/shortner/new", handlers.NewLinkForm)
	r.Get("/s/{linkID}", handlers.RedirectLink)
	// File Drop
	r.Get("/drop/new", handlers.NewDropPage)
	r.Post("/drop/new", handlers.UploadFile)
	r.Get("/drop/dl/{dropID}", handlers.Home)
	fileServer(r, "/drops", http.Dir("./drops")) // download drops
	// Static Files (CSS/JS/Images)
	fileServer(r, "/static", http.Dir("./static"))

	logger.Info("Starting at :8963")
	err := http.ListenAndServe(":8963", r)
	if err != nil {
		logger.Error(err)
	}
}
