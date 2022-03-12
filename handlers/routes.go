package handlers

import (
	"net/http"
	"time"

	"github.com/IceWreck/BetterBin/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(app *config.Application) http.Handler {
	r := chi.NewRouter()

	// custom error handlers
	r.NotFound(notFoundResponse(app))
	r.MethodNotAllowed(methodNotAllowedResponse(app))

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/paste/new", http.StatusTemporaryRedirect)
	})
	r.Get("/home", home(app))
	// Paste Bin
	r.Get("/paste/new", newPastePage(app))
	r.Post("/paste/new", newPasteForm(app))
	r.Get("/paste/view/{pasteID}", viewPastePage(app))
	r.Get("/paste/raw/{pasteID}", viewPasteRaw(app))
	r.Post("/paste/view/{pasteID}", viewPastePage(app)) // required for password
	r.Post("/paste/raw/{pasteID}", viewPasteRaw(app))   // required for password
	// Link Shortner
	r.Get("/shortner/new", newLinkPage(app))
	r.Post("/shortner/new", newLinkForm(app))
	r.Get("/s/{linkID}", redirectLink(app))
	// File Drop
	r.Get("/drop/new", newDropPage(app))
	r.Post("/drop/new", uploadFile(app))
	r.Get("/drop/dl/{dropID}", viewDrop(app))
	fileServer(app, r, "/drops", http.Dir("./drops")) // download drops
	// Static Files (CSS/JS/Images)
	fileServer(app, r, "/static", http.Dir("./static"))

	return r
}
