package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

//go:embed static template
var files embed.FS

type server struct {
	router   *echo.Echo
	template *template.Template
}

func main() {

	s, err := CreateServer()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on :8080")
	if err := s.router.Start(":8080"); err != nil {
		panic(err)
	}

}

// Middleware pour logger les headers en entrée et en sortie
func logHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Logger les headers en entrée
		fmt.Println("=== Headers en entrée ===")
		for name, values := range c.Request().Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		fmt.Println("========================")

		// appel du handler
		if err := next(c); err != nil {
			c.Error(err)
		}

		// Logger les headers en sortie
		fmt.Println("=== Headers en sortie ===")
		for name, values := range c.Response().Header() {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		fmt.Println("========================")

		return nil
	}
}

// Middleware pour ajouter un identifiant unique à chaque requête
func requestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := uuid.New().String()
		// stocker dans le context Echo et header de réponse
		c.Set("requestID", id)
		c.Response().Header().Set("X-Request-ID", id)
		// propagate context if needed
		ctx := context.WithValue(c.Request().Context(), "requestID", id)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
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

	e := echo.New()
	s := &server{
		router:   e,
		template: template,
	}

	// Middlewares Echo
	e.Use(echomw.Recover())
	e.Use(requestID)
	e.Use(echomw.Logger())
	e.Use(logHeaders)

	// Servir les fichiers statiques via http.FileServer en utilisant l'embed FS
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(static))))

	e.GET("/hello", s.hello)

	return s, nil
}

func (s *server) hello(c echo.Context) error {
	message := "ceci est un message de test"
	if err := s.template.ExecuteTemplate(c.Response().Writer, "hello", map[string]string{"Message": message}); err != nil {
		return err
	}
	return nil
}
