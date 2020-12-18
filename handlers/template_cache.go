package handlers

import (
	"html/template"

	"github.com/IceWreck/BetterBin/logger"
)

var templateCache = cacheTemplates()

// we want to parse and cache templates on program start instead of parsing them per request
func cacheTemplates() map[string](*template.Template) {
	logger.Info("parsing templates")
	cache := map[string](*template.Template){
		"home": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/home.html",
		}...)),
		"view_paste": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/view_paste.html",
		}...)),
		"new_paste": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/new_paste.html",
		}...)),
	}
	return cache
}
