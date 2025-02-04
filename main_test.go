package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	s, err := CreateServer()
	if err != nil {
		t.Fatalf("Creating server failed: %v", err)
	}
	
	// Create a new request targeting the /hello route.
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("Creating request failed: %v", err)
	}

	// Record the response.
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	// Check for a 200 OK response.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Optionally, check if content contains expected string. Replace "Hello" with relevant content.
	if !strings.Contains(rr.Body.String(), "Bonjour") {
		t.Errorf("handler returned unexpected body: %v", rr.Body.String())
	}
}

func TestStaticFileHandler(t *testing.T) {
	s, err := CreateServer()
	if err != nil {
		t.Fatalf("Creating server failed: %v", err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Creating request failed: %v", err)
	}

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Static file handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}