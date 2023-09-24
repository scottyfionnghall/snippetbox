package main

import (
	"net/http"
	"strings"

	"github.com.scottyfionnghall.snippetbox/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes method returns a servemux containing out application routes
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Define file server
	// Tak the ui.Files embedded filesystem and convert it to a http.FS type
	// so that it satisfues the http.FileSystem interface.
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", neuter(fileServer))

	router.HandlerFunc(http.MethodGet, "/ping", ping)
	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	// Define handlers containing dynamic iddlware chain
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodDelete, "/snippet/delete", dynamic.ThenFunc(app.snippetDelete))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	return standard.Then(router)
}

// Disable directory listing for static directory
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
