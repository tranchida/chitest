package main

import (
	"compress/flate"
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

type TemplateHandler struct {
	template *template.Template
}

func main() {

	static, _ := fs.Sub(files, "static")
	template, err := template.ParseFS(files, "template/*.html")
	if err != nil {
		panic(err)
	}

	handler := &TemplateHandler{template: template}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.NewCompressor(flate.DefaultCompression).Handler)
	r.Handle("/*", http.FileServer(http.FS(static)))

	r.Get("/hello", handler.hello)

	fmt.Println("Starting server on :8080")

	if err = http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}

func (h *TemplateHandler) hello(w http.ResponseWriter, r *http.Request) {
	h.template.ExecuteTemplate(w, "hello", nil)
}
