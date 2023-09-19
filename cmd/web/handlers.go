package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com.scottyfionnghall.snippetbox/internal/models"
)

// Add a "/" handler function
func (app *appliaction) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/".
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.html", data)
}

// Add a SnippetView handler function.
func (app *appliaction) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Uses notFound() helper
		return
	}
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.html", data)
}

// Add a SnippetCreate handler function.
func (app *appliaction) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	var p models.Snippet
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.badRequest(w)
		return
	}
	if p.Title == "" || p.Content == "" {
		app.badRequest(w)
		return
	}
	// Pass the data to the SnippetModel.Insert() method, reciving the
	// ID of the new record back
	id, err := app.snippets.Insert(p.Title, p.Content, 7)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *appliaction) snippetDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.badRequest(w)
		return
	}
	err = app.snippets.Delete(id)
	if err != nil {
		app.badRequest(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
