package main

import "net/http"

// The routes method returns a servemux containing out application routes
func (app *appliaction) routes() *http.ServeMux {
	mux := http.NewServeMux()
	// Define file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	// Define handler functions
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
