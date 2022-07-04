package main

import (
	"encoding/json"
    "errors"
	"net/http"
)


type APIError struct {
    status 		int
    message     string
}

func (e APIError) Error() string {
    return e.message
}

var (
    ErrNotFound  = &APIError{status: http.StatusNotFound, message: "Route not found"}
    ErrServer = &APIError{status: http.StatusInternalServerError, message: "Something went wrong"}
)

// TODO: add logging
func JSONError(w http.ResponseWriter, err error) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")

    var apiErr APIError
    if errors.As(err, &apiErr) {
        w.WriteHeader(apiErr.status)
        json.NewEncoder(w).Encode(apiErr.message)
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode("Something went wrong")

    }
}