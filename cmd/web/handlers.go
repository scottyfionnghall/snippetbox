package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com.scottyfionnghall.snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

// This handler returns home page.
func (app *appliaction) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/".

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.html", data)
}

// This handler allow to show the user particular snippet based on the passed ID.
func (app *appliaction) snippetView(w http.ResponseWriter, r *http.Request) {
	// httprouter extracts all parameters passed in the request in a form
	// of a slice
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

// This handler handels all request to create a new snippet
func (app *appliaction) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *appliaction) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
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
