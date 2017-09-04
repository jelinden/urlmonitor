package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func urlHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m, _ := json.Marshal(urls)
		w.Header().Add("Content-type", "application/json")
		fmt.Fprintf(w, `{"urls": `+string(m)+`}`)
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
