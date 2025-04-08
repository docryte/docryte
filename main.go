package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func handleStatic(mux *http.ServeMux) {
	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))
}

func main() {
	mux := http.NewServeMux()
	handleStatic(mux)

	if temp, err := template.ParseFiles("./templates/main.html"); err != nil {
		log.Fatal("Error while loading template: ", err)
	} else {
		mux.HandleFunc("/", indexPage(temp))
	}

	if err := http.Serve(autocert.NewListener("docryte.site"), mux); err != nil {
		log.Fatal("Error while starting server: ", err)
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
