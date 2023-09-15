package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Add a "/" handler function
func home(w http.ResponseWriter, r *http.Request) {
	if methodAllowed(w, r, http.MethodGet) {
		// Check if the current request URL path exactly matches "/".
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Snippet home page..."))
	}

}

// Add a SnippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	if methodAllowed(w, r, http.MethodGet) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

// Add a SnippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if methodAllowed(w, r, http.MethodPost) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Create a new snippet..."))
	}
}

// Check if requset is using allowed method
func methodAllowed(w http.ResponseWriter, r *http.Request, allowedRespone string) bool {
	if r.Method != allowedRespone {
		w.Header().Set("Allow", allowedRespone)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return false
	} else {
		return true
	}
}
