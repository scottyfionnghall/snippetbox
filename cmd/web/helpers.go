package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// decodePostForm() helper method. Made for a specific type of error
// If app.formDecoder.Decode() gets something that isn't a non-nil pointer, then
// Decode() will return a form.InvalidDecodeError.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination
	// as the first paramter
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var InvalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &InvalidDecodeError) {
			panic(err)
		}
		// For all other errors, we return them as normal
		return err
	}
	return nil
}

// Create a newTemplateData() helper, which returns a pointer to a
// templateData struct initialize with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	err = app.errorLog.Output(2, trace)
	if err != nil {
		app.serverError(w, err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding
// description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The notFound helper sends a 404 Not Founds response to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
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
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}
}

// Return true if the current request if rom an authenticated user,
// otherwise return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
