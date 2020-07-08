package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sukhjit/url-shortener/model"
	"github.com/sukhjit/url-shortener/repo"
	"github.com/sukhjit/url-shortener/repo/dynamodb"
	"github.com/sukhjit/url-shortener/repo/inmemory"
	"github.com/sukhjit/util"
)

var (
	errURLNotFound = fmt.Errorf("URL does not exist")
	errorLogger    = log.New(os.Stderr, "[ERROR] ", log.Llongfile)
	shortenerDB    repo.Shortener
)

// NewHandler function create routes and return mux router
func NewHandler(isLocal bool, awsRegion, dynamoDBTable string) *mux.Router {
	if isLocal {
		shortenerDB = inmemory.NewShortener()
	} else {
		shortenerDB = dynamodb.NewShortener(awsRegion, dynamoDBTable)
	}

	return buildRouter()
}

func buildRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/info/{slug}", infoHandler).Methods("GET")
	router.HandleFunc("/add", addHandler).Methods("POST")
	router.HandleFunc("/status", statusHandler).Methods("GET")
	router.HandleFunc("/{slug}", redirectHandler).Methods("GET")

	return router
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	responseJSONHandle(w, http.StatusOK, payload)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	resp, err := shortenerDB.Info(slug)
	if err != nil {
		responseErrorHandle(w, http.StatusInternalServerError, err)
		return
	}

	if resp == nil || resp.Slug == "" {
		responseErrorHandle(w, http.StatusNotFound, errURLNotFound)
		return
	}

	responseJSONHandle(w, http.StatusOK, resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	url, err := shortenerDB.Load(slug)
	if err != nil {
		errMsg := fmt.Errorf("Unable to load from database: %v", err)
		responseErrorHandle(w, http.StatusInternalServerError, errMsg)
		return
	}

	if url == "" {
		responseErrorHandle(w, http.StatusNotFound, errURLNotFound)
		return
	}

	http.Redirect(w, r, string(url), http.StatusFound)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	reqPayload := struct {
		URL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		errMsg := fmt.Errorf("Invalid request payload")
		responseErrorHandle(w, http.StatusBadRequest, errMsg)
		return
	}

	if len(reqPayload.URL) == 0 {
		errMsg := fmt.Errorf("Missing 'url' variable")
		responseErrorHandle(w, http.StatusBadRequest, errMsg)
		return
	}

	url, err := url.ParseRequestURI(reqPayload.URL)
	if err != nil {
		errMsg := fmt.Errorf("Not valid url: %s", reqPayload.URL)
		responseErrorHandle(w, http.StatusBadRequest, errMsg)
		return
	}

	obj := &model.Shortener{
		Slug: util.RandomString(8),
		URL:  url.String(),
	}

	err = shortenerDB.Add(obj)
	if err != nil {
		errMsg := fmt.Errorf("Could not save to database: %v", err)
		responseErrorHandle(w, http.StatusInternalServerError, errMsg)
		return
	}

	responseJSONHandle(w, http.StatusCreated, obj)
}

func responseJSONHandle(w http.ResponseWriter, statusCode int, payload interface{}) {
	result, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(result)
}

func responseErrorHandle(w http.ResponseWriter, code int, err error) {
	payload := map[string]string{
		"error": err.Error(),
	}

	if code > 499 {
		errID := uuid.New().String()

		// log error
		errorLogger.Printf("ErrorID: %s, %v", errID, err.Error())

		// add error id to response
		payload["code"] = errID
		payload["error"] = "Internal server error"
	}

	responseJSONHandle(w, code, payload)
}
