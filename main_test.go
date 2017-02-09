package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
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

	// expected := "Hello World\n"
	// if body := rr.Body.String(); body != expected {
	// 	t.Errorf("wrong body: got %v want %v", body, expected)
	// }

	expected := `text/html; charset=utf-8`
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("wrong content type: got %v want %v", contentType, expected)
	}

}
