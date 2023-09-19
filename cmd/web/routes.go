package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes method returns a servemux containing out application routes
func (app *appliaction) routes() http.Handler {
	mux := http.NewServeMux()
	// Define file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	// Define handler functions
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	mux.HandleFunc("/snippet/delete", app.snippetDelete)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	return standard.Then(mux)
}
