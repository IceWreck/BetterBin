package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/IceWreck/BetterBin/config"
	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

var errPasteExpired = errors.New("paste has expired")
var errPasteNotFound = errors.New("paste not found")
var errNoContent = errors.New("no content for new paste")
var errPasswordRequired = errors.New("password required")

const timeLayout = "2006-01-02 15:04:05"

// newPasteForm (API) - new paste using a POST request
func newPasteForm(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		expiry := r.PostFormValue("expiry")
		password := r.PostFormValue("password")

		// string (technically boolean) input from forms should be converted to ints
		burn := 0

		// error out if content is empty
		if len(content) < 1 {
			renderError(w, r, errNoContent, http.StatusUnprocessableEntity)
			return
		}

		// convert expiry to a form suitable for sql
		// never place user input directly into the query
		switch expiry {
		case "year":
			expiry = "1 years"
		case "month":
			expiry = "1 month"
		case "week":
			expiry = "7 days"
		case "day":
			expiry = "1 days"
		case "hour":
			expiry = "1 hours"
		case "10min":
			expiry = "10 minutes"
		case "1min":
			expiry = "1 minutes"
		case "burn":
			expiry = "1 month"
			burn = 1
		default:
			// never expire
			expiry = "999 years"
		}
		pasteID := newID(10)
		logger.Info("creating new paste", title, pasteID)
		if err := db.NewPaste(app, pasteID, title, content, expiry, password, burn); err != nil {
			logger.Info("could not create a new paste")
			renderError(w, r, err, http.StatusInternalServerError)
			return
		}
		renderSuccess(w, r, pasteID)
	}
}

// newPastePage - webpage to make a new paste
func newPastePage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "new_paste", nil)
	}
}

// viewPastePage - webpage to view a paste
func viewPastePage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paste, err := getPaste(app, r)
		if err == errPasswordRequired {
			renderTemplate(w, "password_required", nil)
			return
		} else if err == errPasteExpired {
			renderTemplate(w, "paste_expired", nil)
			return
		} else if err != nil {
			renderTemplate(w, "paste_not_found", nil)
			return
		}

		if paste.Preview == "code" {
			htmlw := new(strings.Builder)
			lexer := lexers.Match("foo." + paste.PreviewLanguage)
			// if that didn't math try auto analyze
			if lexer == nil {
				lexer = lexers.Analyse(paste.Content)
			}
			// only proceed if that worked
			if lexer != nil {
				logger.Debug("found lexer", lexer.Config().Name)
				lexer = chroma.Coalesce(lexer)
				formatter := html.New(html.WithLineNumbers(true), html.LineNumbersInTable(true))
				iterator, err := lexer.Tokenise(nil, paste.Content)
				if err == nil {
					err = formatter.Format(htmlw, styles.Colorful, iterator)
					if err != nil {
						logger.Error("syntax formatting error", err)
					} else {
						// logger.Debug("syntax highlighted HTML is ", htmlw.String())
						paste.ContentHTML = template.HTML(htmlw.String())
					}
				} else {
					logger.Error("syntax tokenization error", err)
				}
			} else {
				// just use the plaintext preview type if no lexer was found
				paste.Preview = "plain"
			}
		} else if paste.Preview == "markdown" {
			htmlw := new(strings.Builder)
			md := goldmark.New(
				goldmark.WithExtensions(extension.GFM),
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
				),
				goldmark.WithRendererOptions(
					goldmarkhtml.WithHardWraps(),
					goldmarkhtml.WithXHTML(),
				),
			)
			if err := md.Convert([]byte(paste.Content), htmlw); err != nil {
				logger.Error("markdown conversion error", err)
				paste.Preview = "plain"
			}
			paste.ContentHTML = template.HTML(htmlw.String())
		}

		renderTemplate(w, "view_paste", paste)
	}
}

// viewPasteRaw - curl/wget friendly raw paste contents
func viewPasteRaw(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paste, err := getPaste(app, r)
		if err != nil {
			renderError(w, r, err, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(paste.Content))
	}
}

// getPaste has common stuff required in both ViewPastePage and ViewPasteRaw
func getPaste(app *config.Application, r *http.Request) (db.Paste, error) {
	pasteID := chi.URLParam(r, "pasteID")
	// password can be either URL query parameters or POST/PUT values
	password := r.FormValue("password")

	logger.Debug("fetching paste", pasteID)
	paste, err := db.GetPaste(app, pasteID)
	if err != nil {
		// error prolly means that it can't find paste
		return paste, errPasteNotFound
	}

	// expiry check
	expiry, err := time.Parse(timeLayout, paste.Expiry)
	if err != nil {
		logger.Error("unable to parse sql datetime", err)
	}
	if expiry.Sub(time.Now()) < 0 {
		// expired
		return paste, errPasteExpired
	}

	// password check
	if password != paste.Password {
		return paste, errPasswordRequired
	}

	// burn paste for future fetches if burn=1
	if paste.Burn == 1 {
		logger.Info("burning paste", pasteID)
		db.BurnPaste(app, pasteID)
	}

	// preview information from query paramater
	preview := r.FormValue("preview")
	if preview == "code" {
		paste.Preview = "code"
	} else if preview == "markdown" {
		paste.Preview = "markdown"
	}

	// check if the lang query parameter exists
	_, langParamExists := r.Form["lang"]
	if langParamExists {
		paste.Preview = "code"
		paste.PreviewLanguage = r.FormValue("lang")
	}

	return paste, nil
}
