package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

// Define an application struct to hold the application-wide dependencies.
type appliaction struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Command-line argument "addr" to define address on wich the server
	// will be listening.
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()
	//Add info and error logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//Initialize a new instance of our applicaiton struct, containing
	// the dependencies
	app := &appliaction{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	// Initialize a new http.Server struct
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// Start server
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
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
