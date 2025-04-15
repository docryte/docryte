package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func redirectToTls(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func handleStatic(mux *http.ServeMux) {
	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))
}

func main() {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "DEBUG"
	}

	mux := http.NewServeMux()
	handleStatic(mux)

	if temp, err := template.ParseFiles("./templates/main.html"); err != nil {
		log.Fatal("Error while loading template: ", err)
	} else {
		mux.HandleFunc("/", indexPage(temp))
	}

	if mode == "PROD" {
		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
				log.Fatalf("ListenAndServe error: %v", err)
			}
		}()
		http.Serve(autocert.NewListener("docryte.site"), mux)
	} else if mode == "DEBUG" {
		http.ListenAndServe(":80", mux)
	} else {
		log.Fatal("MODE environment variable should be either PROD or DEBUG: ", mode)
	}
}

func indexPage(template *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := template.Execute(w, nil); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}
}
