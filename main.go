package main

import (
	"net/http"
)

func likesHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/api/users/", likesHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
