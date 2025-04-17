package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"

	"docryte/src/config"
	"docryte/src/types"
)

func Init(cfg config.Config) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fileServer(r, "/static", http.Dir("./static"))
	r.Get("/", indexPage)
	r.Post("/contact", contactRequest(cfg.TelegramToken, cfg.UserId))

	go func() { http.ListenAndServe(":80", r) }()
	http.Serve(autocert.NewListener(cfg.Domain), r)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/main.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func contactRequest(token string, userId string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var cr types.ContactRequest
		decoder.Decode(&cr)
		go notifyNewContact(&cr, token, userId)
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
