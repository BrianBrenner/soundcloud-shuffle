package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// TODO: store clientId, before doing anything check if valid, if not then fetch one

var likesRoute = regexp.MustCompile(`/api/likes$`)

func route(w http.ResponseWriter, r *http.Request) {
	switch {
	case likesRoute.MatchString(r.URL.Path):
		likesHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Route not found\n"))
	}
}

func likesHandler(w http.ResponseWriter, r *http.Request) {
	profileUrl, ok := r.URL.Query()["url"]

	if !ok || len(profileUrl[0]) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Invalid params passed. A soundcloud profile url must bed passed as the url query param\n"))
		return
	}

	// TODO: handle errs
	clientId, err := getClientId()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong\n"))
	}

	userId, err := getUserId(profileUrl[0], clientId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong\n"))
	}

	likes, err := getUserLikes(userId, clientId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong\n"))
	}

	jsonLikes, err := json.Marshal(likes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong\n"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonLikes)
}

func main() {
	http.HandleFunc("/", route)

	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		panic(err)
	}
}
