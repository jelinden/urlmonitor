package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func urlHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m, _ := json.Marshal(getUrls())
		n, _ := json.Marshal(getNames())
		w.Header().Add("Content-type", "application/json")
		fmt.Fprintf(w, `{"urls": `+string(m)+`, "names": `+string(n)+`}`)
	})
}

func jsonHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		url := strings.TrimPrefix(r.URL.Path, "/json/")
		if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
			url = strings.Replace(url, "http:/", "http://", 1)
			url = strings.Replace(url, "https:/", "https://", 1)
		}
		fmt.Fprintf(w, `{"times": `+cacheGet(url)+`}`)
	})
}

func getUrls() []string {
	var urls []string
	for _, domain := range domains {
		urls = append(urls, domain.Url)
	}
	return urls
}

func getNames() []string {
	var names []string
	for _, domain := range domains {
		names = append(names, domain.Name)
	}
	return names
}
