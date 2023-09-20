package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com.scottyfionnghall.snippetbox/internal/models"
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

	// Use the filepath.Glob() function to get a slice of all filepath that
	// math the patter "./ui/html/pages/*.html"
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// Loop through the page filepaths
	for _, page := range pages {
		// Extracte the file name and assign it to the new name variable
		name := filepath.Base(page)
		// Parse the base template file into a template set
		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		// Call ParseGlob() *on this template set* to add any patials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
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
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}
