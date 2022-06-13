package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

// The HTML returned from soundcloud has a lot of script tags. If you follow
// the URL in the last tag you get some javacript back which contains a client ID
func getClientIdUrl(res *http.Response) string {
	body, _ := ioutil.ReadAll(res.Body)

	urlRegex := regexp.MustCompile(`<script crossorigin src="(.*?)"`)
	regExUrls := urlRegex.FindAllStringSubmatch(string(body), -1)

	var urls []string
	for _, arr := range regExUrls {
		urls = append(urls, arr[1])
	}

	return urls[len(urls)-1]
}

// TODO: how to handle erros
func getClientId() string {
	res, _ := http.Get("https://soundcloud.com/")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	url := getClientIdUrl(res)

	idRes, _ := http.Get(url)
	idBody, _ := ioutil.ReadAll(idRes.Body)

	// TODO: handle empty matches
	idRegex := regexp.MustCompile(`,client_id:"(.*?)"`)
	idVal := idRegex.FindStringSubmatch(string(idBody))

	return idVal[1]
}

func getUserLikes(userId string, clientId string) string {
	// TODO: parse as JSON, also follow next urls
	res, _ := http.Get("https://api-v2.soundcloud.com/users/" + userId + "/likes?client_id=" + clientId + "&limit=1")

	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
