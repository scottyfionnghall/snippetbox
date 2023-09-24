package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// Function to recover from panic
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a defered function to recover from panic
		defer func() {
			// Use the builtin recover function to check if there has been
			// a panic or not.
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Function to log every request into infoLog
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// Function to add security layers to every request
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		w.Header().Set(
			"Referrer-Policy",
			"origin-when-cross-origin",
		)
		w.Header().Set(
			"X-Content-Type-Options",
			"nosniff",
		)
		w.Header().Set(
			"X-Frame-Options",
			"deny",
		)
		w.Header().Set(
			"X-XSS-Protection",
			"0",
		)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise, set the "Cache-Control: no store" header so that pages
		// require authentication are not stored in the users browser cache
		w.Header().Add("Cache-Control", "no-store")
		// And call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// NoSurf middleware to prevent cross-site request forgery attacks.
func noSurf(next http.Handler) http.Handler {
	crsfHandler := nosurf.New(next)
	crsfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return crsfHandler
}
