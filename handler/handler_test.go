package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/sukhjit/url-shortener/model"
)

type MockDB struct{}

func (m *MockDB) Load(slug string) (string, error) {
	return "", fmt.Errorf("failed to load")
}

func (m *MockDB) Add(item *model.Shortener) error {
	return fmt.Errorf("failed to add")
}

func (m *MockDB) Info(slug string) (*model.Shortener, error) {
	return nil, fmt.Errorf("failed to info")
}

func setup() {
	_ = NewHandler(true, "", "")
}

func TestStatusHandler(t *testing.T) {
	setup()
	a := assert.New(t)

	req, err := http.NewRequest("GET", "/status", http.NoBody)
	a.Nil(err)

	rec := httptest.NewRecorder()

	statusHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	a.Equal(res.StatusCode, http.StatusOK)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)

	a.Equal(`{"status":"ok"}`, string(b))
}

func TestAddHandlerWithInvalidPayload(t *testing.T) {
	setup()
	a := assert.New(t)

	req, err := http.NewRequest("POST", "/add", strings.NewReader(""))
	a.Nil(err)

	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusBadRequest, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)

	a.Equal(`{"error":"invalid request payload"}`, string(b))
}

func TestAddHandlerWithoutUrlParam(t *testing.T) {
	setup()
	a := assert.New(t)

	jsonStr := []byte(`{"urlx":"t"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	a.Nil(err)

	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusBadRequest, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Equal(`{"error":"missing 'url' variable"}`, string(b))
}

func TestAddHandlerWithInvalidUrl(t *testing.T) {
	setup()
	a := assert.New(t)

	jsonStr := []byte(`{"url":"http//www.itestit.com"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	a.Nil(err)

	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusBadRequest, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Equal(`{"error":"not valid url: http//www.itestit.com"}`, string(b))
}

func TestAddHandlerReturn500(t *testing.T) {
	shortenerDB = &MockDB{}
	a := assert.New(t)

	jsonStr := []byte(`{"url":"http://www.itestit.com"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	a.Nil(err)

	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusInternalServerError, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)

	a.Contains(string(b), `"error":"Internal server error"`)
}

func TestAddHandler(t *testing.T) {
	setup()
	a := assert.New(t)

	jsonStr := []byte(`{"url":"http://www.itestit.com"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonStr))
	a.Nil(err)

	rec := httptest.NewRecorder()

	addHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusCreated, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)

	a.Contains(string(b), `"url":"http://www.itestit.com"`)
	a.Contains(string(b), `"visits":0`)
}

func TestGetHandlerReturn404(t *testing.T) {
	setup()
	a := assert.New(t)

	req, err := http.NewRequest("GET", "/fake-val", http.NoBody)
	a.Nil(err)

	rec := httptest.NewRecorder()

	redirectHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusNotFound, res.StatusCode)
}

func TestGetHandlerReturn500(t *testing.T) {
	shortenerDB = &MockDB{}
	a := assert.New(t)

	req, err := http.NewRequest("GET", "/fake-val", http.NoBody)
	a.Nil(err)

	rec := httptest.NewRecorder()

	redirectHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusInternalServerError, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Contains(string(b), `"error":"Internal server error"`)
}

func TestGetHandler(t *testing.T) {
	setup()
	a := assert.New(t)

	req, err := http.NewRequest("GET", "who-cares", http.NoBody)
	a.Nil(err)

	req = mux.SetURLVars(req, map[string]string{
		"slug": "go",
	})

	rec := httptest.NewRecorder()

	redirectHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusFound, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Contains(string(b), `<a href="https://www.google.com">Found</a>`)
}

func TestInfoHandlerReturn404(t *testing.T) {
	setup()
	a := assert.New(t)

	req, err := http.NewRequest("GET", "/info/123", http.NoBody)
	a.Nil(err)

	rec := httptest.NewRecorder()

	infoHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusNotFound, res.StatusCode)
}

func TestInfoHandlerReturn500(t *testing.T) {
	shortenerDB = &MockDB{}
	a := assert.New(t)

	req, err := http.NewRequest("GET", "/info/123", http.NoBody)
	a.Nil(err)

	rec := httptest.NewRecorder()

	infoHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusInternalServerError, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Contains(string(b), `"error":"Internal server error"`)
}

func TestInfoHandler(t *testing.T) {
	setup()
	a := assert.New(t)

	req, err := http.NewRequest("GET", "who-cares", http.NoBody)
	a.Nil(err)

	req = mux.SetURLVars(req, map[string]string{
		"slug": "go",
	})

	rec := httptest.NewRecorder()

	infoHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	a.Equal(http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Contains(string(b), `"slug":"go"`)
	a.Contains(string(b), `"url":"https://www.google.com"`)
	a.Contains(string(b), `"visits":0`)
}
