package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// Add a "/" handler function
func (app *appliaction) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/".
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/base.hml",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Uses the serverError() helper
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err) // Uses the serverError() helper
		return
	}

}

// Add a SnippetView handler function.
func (app *appliaction) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Uses notFound() helper
		return
	}
	files := []string{
		"./ui/html/partials/nav.html",
		"./ui/html/base.html",
		"./ui/html/pages/view.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", id)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

// Add a SnippetCreate handler function.
func (app *appliaction) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Create a new snippet..."))

}
