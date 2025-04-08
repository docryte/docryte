package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func handleStatic() {
	static := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))
}

func main() {
	handleStatic()
	temp, err := template.ParseFiles("./templates/main.html")
	if err != nil {
		log.Fatal("Error while loading template: ", err)
	}

	http.HandleFunc("/", indexPage(temp))

	err = http.Serve(autocert.NewListener("docryte.site"), nil)
	if err != nil {
		log.Fatal("Error while starting server: ", err)
	}
}

func indexPage(template *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := template.Execute(w, nil)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}
}
