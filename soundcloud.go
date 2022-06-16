package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Likes struct {
	Collection []struct {
		Track struct {
			URL string `json:"permalink_url"`
		} `json:"track"`
	} `json:"collection"`
	NextHref string `json:"next_href"`
}

type User struct {
	ID int `json:"id"`
}

// TODO: handle errors like 401, 404, etc, do you have to check response code? or is non-200 auto error

var client = &http.Client{Timeout: 15 * time.Second}

// The HTML returned from soundcloud has a lot of script tags. If you follow
// the URL in the last tag you get some javacript back which contains a client ID
func getClientIdUrl(res *http.Response) (string, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	urlRegex := regexp.MustCompile(`<script crossorigin src="(.*?)"`)
	regExUrls := urlRegex.FindAllStringSubmatch(string(body), -1)

	var urls []string
	for _, arr := range regExUrls {
		urls = append(urls, arr[1])
	}

	return urls[len(urls)-1], nil
}

func getClientId() (string, error) {
	res, err := client.Get("https://soundcloud.com/")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	url, err := getClientIdUrl(res)
	if err != nil {
		return "", err
	}

	idRes, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer idRes.Body.Close()

	idBody, err := ioutil.ReadAll(idRes.Body)
	if err != nil {
		return "", err
	}

	// TODO: handle empty matches
	idRegex := regexp.MustCompile(`,client_id:"(.*?)"`)
	idVal := idRegex.FindStringSubmatch(string(idBody))

	return idVal[1], nil
}

func getUserId(profileUrl string, clientId string) (string, error) {
	params := url.Values{
		"client_id": []string{clientId},
		"url":       []string{profileUrl},
	}

	res, err := client.Get("https://api-v2.soundcloud.com/resolve?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var user = new(User)
	err = json.Unmarshal(body, &user)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(user.ID), nil
}

func getUserLikes(userId string, clientId string) ([]string, error) {
	params := url.Values{
		"limit":     []string{"1"},
		"client_id": []string{clientId},
	}

	// TODO: also follow next urls
	res, err := client.Get("https://api-v2.soundcloud.com/users/" + userId + "/likes?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var likes = new(Likes)
	err = json.Unmarshal(body, &likes)

	urls := []string{}
	for _, track := range likes.Collection {
		urls = append(urls, track.Track.URL)
	}

	shuffled := make([]string, len(urls))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(urls))

	for i, v := range perm {
		shuffled[v] = urls[i]
	}

	return shuffled, nil
}
