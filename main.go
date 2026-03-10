package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

//go:embed static template
var files embed.FS

type server struct {
	router   *chi.Mux
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

// Middleware pour logger les headers en entrée et en sortie
func logHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logger les headers en entrée
		fmt.Println("\n=== Headers en entrée ===")
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		fmt.Print("========================\n\n")

		// Wrapper pour capturer les headers en sortie
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		// Logger les headers en sortie
		fmt.Println("\n=== Headers en sortie ===")
		for name, values := range ww.Header() {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		fmt.Print("========================\n\n")
	})
}

// Middleware pour ajouter un identifiant unique à chaque requête
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
		router:   chi.NewRouter(),
		template: template,
	}

	// Ajout du middleware requestID en premier
	s.router.Use(requestID)
	s.router.Use(logHeaders)
	s.router.Use(middleware.Logger)

	s.router.Handle("/*", http.FileServer(http.FS(static)))

	s.router.Get("/hello", s.hello)

	return s, nil
}

func (s *server) hello(w http.ResponseWriter, r *http.Request) {
	message := "ceci est un message de test"
	s.template.ExecuteTemplate(w, "hello", map[string]string{"Message": message})
}
