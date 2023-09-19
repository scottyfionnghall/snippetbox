package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Create a newTemplateData() helper, which returns a pointer to a
// templateData struct initialize with the current year.
func (app *appliaction) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *appliaction) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding
// description to the user.
func (app *appliaction) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The notFound helper sends a 404 Not Founds response to the user
func (app *appliaction) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *appliaction) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *appliaction) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrive the appropriate template set from the cache based on the page
	// name. If no entry exists, the create a new error.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return

	}
	// Initialize a new buffer
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Write out the provided HTTP status code
	w.WriteHeader(status)
	// Write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}
