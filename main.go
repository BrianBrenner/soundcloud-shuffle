package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

// TODO: store clientId, before doing anything check if valid, if not then fetch one

var likesRoute = regexp.MustCompile(`/api/likes$`)

func router(w http.ResponseWriter, r *http.Request) {
	switch {
	case likesRoute.MatchString(r.URL.Path):
		likesHandler(w, r)
	default:
		JSONError(w, r, ErrNotFound)
	}
}

func validateParams(r *http.Request) (string, error) {
	profileUrl, ok := r.URL.Query()["url"]

	if !ok || len(profileUrl[0]) < 1 {
		return "", &APIError{status: http.StatusBadRequest, message: "Invalid params passed. A soundcloud profile url must be passed as the url query param\n"}
	}

	u, err := url.ParseRequestURI(profileUrl[0])
	if err != nil || u.Host != "soundcloud.com" {
		return "", &APIError{status: http.StatusBadRequest, message: "Invalid params passed. A valid soundcloud profile url must be passed as the url query param\n"}
	}

	return profileUrl[0], nil
}

func likesHandler(w http.ResponseWriter, r *http.Request) {
	profileUrl, err := validateParams(r)
	if err != nil {
		JSONError(w, r, err)
		return
	}

	clientId, err := getClientId()
	if err != nil {
		JSONError(w, r, err)
		return
	}

	userId, err := getUserId(profileUrl, clientId)
	if err != nil {
		JSONError(w, r, err)
		return
	}

	likes, err := getUserLikes(userId, clientId)
	if err != nil {
		JSONError(w, r, err)
		return
	}

	jsonLikes, err := json.Marshal(likes)
	if err != nil {
		JSONError(w, r, err)
		return
	}

	// this isn't scalable, with more routes use logging middleware, but fine here
	log.Printf("200 %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonLikes)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	fs := http.FileServer(http.Dir("./web/dist"))

	http.Handle("/api/", http.HandlerFunc(router))
	http.Handle("/", fs)

	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		panic(err)
	}
}
