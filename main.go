package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/render"
)

//go:embed static template
var files embed.FS

type server struct {
	router *chi.Mux
	template *template.Template
}

func main() {

	s, err := CreateServer()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		panic(err)
	}

}

func CreateServer() (*server, error) {

	static, err := fs.Sub(files, "static")
	if err != nil {
		return nil, err
	}

	template, err := template.ParseFS(files, "template/*.html")
	if err != nil {
		return nil, err
	}

	s := &server{
		router: chi.NewRouter(),
		template: template,
	}

	s.router.Use(middleware.Logger)

	s.router.Handle("/*", http.FileServer(http.FS(static)))

	s.router.Get("/hello", s.hello)

	return s, nil
}

func (s *server) hello(w http.ResponseWriter, r *http.Request) {
	s.template.ExecuteTemplate(w, "hello", nil)
}
