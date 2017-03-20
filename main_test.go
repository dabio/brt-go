package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := context{
		templates: template.Must(template.ParseGlob("./views/*.tmpl")),
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.index)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Berlin Racing Team"
	body := rr.Body.String()
	if !strings.Contains(body, expected) {
		t.Errorf("wrong body: got %v want %v", body, expected)
	}

	expected = `text/html; charset=utf-8`
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("wrong content type: got %v want %v", contentType, expected)
	}
}

func TestPathMoved(t *testing.T) {
	req, err := http.NewRequest("GET", "/rennen", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := context{
		templates: template.Must(template.ParseGlob("./views/*.tmpl")),
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.redirect)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}
}

func TestPathMovedSubdirectory(t *testing.T) {
	req, err := http.NewRequest("GET", "/rennen/2012-12-12/blahblah", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := context{
		templates: template.Must(template.ParseGlob("./views/*.tmpl")),
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.redirect)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}
}

// func TestMethodNotAllowed(t *testing.T) {
// 	req, err := http.NewRequest("POST", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	c := context{
// 		templates: template.Must(template.ParseGlob("./views/*.tmpl")),
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(c.index)
//
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusMethodNotAllowed {
// 		t.Errorf("wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
// 	}
// }
