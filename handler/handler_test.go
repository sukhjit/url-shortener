package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func setup() {
	_ = NewHandler(true, "", "")
}

func TestStatusHandler(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	statusHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status: %d, got: %d", http.StatusOK, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	expected := `{"status":"ok"}`
	if expected != string(b) {
		t.Errorf("Expected body: %s, got: %s", expected, b)
	}
}

func TestAddHandlerWithInvalidPayload(t *testing.T) {
	setup()

	req, err := http.NewRequest("POST", "/add", strings.NewReader(""))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusBadRequest
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedStr := `{"error":"Invalid request payload"}`
	if expectedStr != string(b) {
		t.Errorf("Expected body: %s, got: %s", expectedStr, b)
	}
}

func TestAddHandlerWithoutUrlParam(t *testing.T) {
	setup()

	jsonStr := []byte(`{"urlx":"t"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusBadRequest
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedStr := `{"error":"Missing 'url' variable"}`
	if expectedStr != string(b) {
		t.Errorf("Expected body: %s, got: %s", expectedStr, b)
	}
}

func TestAddHandlerWithInvalidUrl(t *testing.T) {
	setup()

	jsonStr := []byte(`{"url":"http//www.itestit.com"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusBadRequest
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedStr := `{"error":"Not valid url: http//www.itestit.com"}`
	if expectedStr != string(b) {
		t.Errorf("Expected body: %s, got: %s", expectedStr, b)
	}
}

func TestAddHandler(t *testing.T) {
	setup()

	jsonStr := []byte(`{"url":"http://www.itestit.com"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status: %d, got: %d", http.StatusCreated, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedStr := `"url":"http://www.itestit.com","visits":0}`
	if !strings.Contains(string(b), expectedStr) {
		t.Errorf("Expected body contains: %s, got: %s", expectedStr, b)
	}
}

func TestGetHandlerReturn404(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/fake-val", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	redirectHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusNotFound
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}
}

func TestGetHandler(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "who-cares", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"slug": "go",
	})

	rec := httptest.NewRecorder()

	redirectHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusFound
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedStr := `<a href="https://www.google.com">Found</a>`
	if !strings.Contains(string(b), expectedStr) {
		t.Errorf("Expected body contains: %s, got: %s", expectedStr, b)
	}
}

func TestInfoHandlerReturn404(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/info/123", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	infoHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusNotFound
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}
}

func TestInfoHandler(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "who-cares", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"slug": "go",
	})

	rec := httptest.NewRecorder()

	infoHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	expected := http.StatusOK
	if res.StatusCode != expected {
		t.Errorf("Expected status: %d, got: %d", expected, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	expectedStr := `{"slug":"go","url":"https://www.google.com","visits":0}`
	if !strings.Contains(string(b), expectedStr) {
		t.Errorf("Expected body contains: %s, got: %s", expectedStr, b)
	}
}
