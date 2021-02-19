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
		"password_required": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/password_required.html",
		}...)),
		"new_paste": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/new_paste.html",
		}...)),
		"new_link": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/new_link.html",
		}...)),
		"new_drop": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/new_drop.html",
		}...)),
		"view_drop": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/view_drop.html",
		}...)),
		"paste_expired": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/paste_expired.html",
		}...)),
		"paste_not_found": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/paste_not_found.html",
		}...)),
		"drop_not_found": template.Must(template.ParseFiles([]string{
			"./templates/layout.html",
			"./templates/pages/drop_not_found.html",
		}...)),
	}
	return cache
}
