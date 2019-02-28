package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// disable cache
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// set hostname (used for demo)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprint(w, "Error:", err)
	}

	data := struct {
		Title         string
		Hostname      string
		BuildDate     string
		BuildreVision string
	}{
		Title:         "Kubernetes Pod Load Balancer Demo (refresh page)",
		Hostname:      hostname,
		BuildDate:     builddate,
		BuildreVision: buildrevision,
	}

	t, err := template.New("index.html").ParseFiles("index.html")

	if err != nil {
		fmt.Fprint(w, "Error:", err)
		fmt.Println("Error:", err)
		return
	}

	err = t.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		fmt.Fprint(w, "Error:", err)
		fmt.Println("Error:", err)
	}

}

// used to dump headers for debugging
func debugHandler(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	// disable cache
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// set hostname (used for demo)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprint(w, "Error:", err)
	}

	// dump headers
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "%v", string(requestDump))
	fmt.Fprintf(w, "Served-By: %v\n", hostname)
	fmt.Fprintf(w, "Serving-Time: %s\n", time.Now().Sub(startTime))
	fmt.Fprintf(w, "Build-Date: %v\n", builddate)
	fmt.Fprintf(w, "Build-Revision: %v", buildrevision)
	return

}

func bootstrapHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// for data passing between middleware
		ctx := r.Context()

		// get current url for redirect if needed
		path := r.URL.Path
		if path == "" { // should not happen, but avoid a crash if it does
			path = "/"
		}

		// force www.sysadmincasts.com -> sysadmincasts.com
		if r.Host == "www.sysadmindemo.com" {
			http.Redirect(w, r, "https://sysadmindemo.com"+path, 301)
			return
		}

		// force http -> https upgrade
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			http.Redirect(w, r, "https://sysadmindemo.com"+path, 301)
			return
		}

		next(w, r.WithContext(ctx))

	}
}

var (
	builddate     = "2019-02-27"
	buildrevision = ""
)

// mux
var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", bootstrapHandler(indexHandler))
	router.HandleFunc("/debug", bootstrapHandler(debugHandler))
	http.Handle("/", router)

	fmt.Println("Listening on port 5005...")
	//http.ListenAndServe("localhost:5005", router)
	http.ListenAndServe(":5005", handlers.CompressHandler(router))

}
