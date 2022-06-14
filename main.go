package main

import (
	"fmt"
	"net/http"
	"regexp"
)

var likesRoute = regexp.MustCompile(`/api/users/(\d+$)`)

func route(w http.ResponseWriter, r *http.Request) {
	switch {
	case likesRoute.MatchString(r.URL.Path):
		likesHandler(w, r)
	default:
		w.Write([]byte("Route not found\n"))
	}
}

func likesHandler(w http.ResponseWriter, r *http.Request) {
	userId := likesRoute.FindStringSubmatch(r.URL.Path)

	// TODO: handle errs
	clientId, _ := getClientId()
	likes, _ := getUserLikes(userId[1], clientId)
	fmt.Fprint(w, likes)
}

func main() {
	http.HandleFunc("/", route)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
