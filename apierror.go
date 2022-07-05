package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type APIError struct {
	status  int
	message string
}

func (e APIError) Error() string {
	return e.message
}

var (
	ErrNotFound = &APIError{status: http.StatusNotFound, message: "Route not found"}
	ErrServer   = &APIError{status: http.StatusInternalServerError, message: "Something went wrong"}
)

func JSONError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		log.Printf("%d %s %s %s Message: %s\n", apiErr.status, r.RemoteAddr, r.Method, r.URL, apiErr.message)

		w.WriteHeader(apiErr.status)
		json.NewEncoder(w).Encode(apiErr.message)
	} else {
		log.Printf("500 %s %s %s Error: %s\n", r.RemoteAddr, r.Method, r.URL, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")

	}
}
