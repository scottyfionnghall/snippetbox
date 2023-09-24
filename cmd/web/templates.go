package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com.scottyfionnghall.snippetbox/internal/models"
	"github.com.scottyfionnghall.snippetbox/ui"
)

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filessystem wich match the pattern 'html/pages/*.html'.
	pages, err := fs.Glob(ui.Files, "/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// Loop through the page filepaths
	for _, page := range pages {
		// Extracte the file name and assign it to the new name variable
		name := filepath.Base(page)
		// Create a slice containing the filepath for the templates
		// we want to parse
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		// Add the template set to the map, using the name of the page
		// as the key
		cache[name] = ts
	}
	return cache, nil
}

// Define a templateData type to act as the holding structure
// for any dynamic data that we want to pass
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}
