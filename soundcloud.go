package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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

var client = &http.Client{Timeout: 15 * time.Second}

func getRequest(url string) (resp *http.Response, err error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Non-200 in getRequest")
	}

	return res, nil
}

func likesRequest(url string) ([]string, string, error) {
	res, err := getRequest(url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	var likes = new(Likes)
	err = json.Unmarshal(body, &likes)

	urls := []string{}
	for _, track := range likes.Collection {
		urls = append(urls, track.Track.URL)
	}

	return urls, likes.NextHref, nil
}

// The HTML returned from soundcloud has a lot of script tags. If you follow
// the URL in the last tag you get some javacript back which contains a client ID
func parseForUrl(res *http.Response) (string, error) {
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

func validateClientId(clientId string) (bool, error) {
	if clientId == "" {
		return false, nil
	}

	params := url.Values{
		"client_id": []string{clientId},
		// this is my profile so I know it is a valid URL
		"url": []string{"https://soundcloud.com/brian-brenner-4"},
	}

	// there might be a request that is faster to check, but this works
	res, err := client.Get("https://api-v2.soundcloud.com/resolve?" + params.Encode())
	if err != nil {
		return false, err
	}

	if res.StatusCode == 200 {
		return true, nil
	}

	// 401 means clientId is invalid
	if res.StatusCode == 401 {
		return false, nil
	}

	return false, errors.New("Recieved " + strconv.Itoa(res.StatusCode) + " in validateClientId")
}

func getClientId() (string, error) {
	log.Println("Fetching client-id")
	res, err := getRequest("https://soundcloud.com/")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	url, err := parseForUrl(res)
	if err != nil {
		return "", err
	}

	idRes, err := getRequest(url)
	if err != nil {
		return "", err
	}
	defer idRes.Body.Close()

	idBody, err := ioutil.ReadAll(idRes.Body)
	if err != nil {
		return "", err
	}

	idRegex := regexp.MustCompile(`,client_id:"(.*?)"`)
	idVal := idRegex.FindStringSubmatch(string(idBody))

	if idVal == nil || len(idVal) < 2 {
		return "", errors.New("Unable to find client_id in regex")
	}

	return idVal[1], nil
}

func getUserId(profileUrl string, clientId string) (string, error) {
	params := url.Values{
		"client_id": []string{clientId},
		"url":       []string{profileUrl},
	}

	res, err := getRequest("https://api-v2.soundcloud.com/resolve?" + params.Encode())
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
		// this seems to have no limit, but don't want soundcloud to block me or anything
		// 200 is the limit on apiv1
		"limit":     []string{"10000"},
		"client_id": []string{clientId},
	}

	// The soundcloud API returns paginated requests by adding a next_href property to the response. We follow
	// this till it becomes nil, meaning we've reached the end of the pagination
	finalUrls, nextUrl, err := likesRequest("https://api-v2.soundcloud.com/users/" + userId + "/likes?" + params.Encode())
	if err != nil {
		return nil, err
	}

	for nextUrl != "" {
		urls, localNextUrl, err := likesRequest(nextUrl + "&client_id=" + clientId)
		if err != nil {
			return nil, err
		}

		finalUrls = append(finalUrls, urls...)
		nextUrl = localNextUrl
	}

	shuffled := make([]string, len(finalUrls))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(finalUrls))

	for i, v := range perm {
		shuffled[v] = finalUrls[i]
	}

	return shuffled, nil
}
