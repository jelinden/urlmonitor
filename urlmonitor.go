package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/headzoo/surf/browser"
	"github.com/jelinden/blig/app/util"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/headzoo/surf.v1"
)

const schedule = 60 * time.Second

var urls = []string{
	"https://www.kauppalehti.fi/",
	"https://m.kauppalehti.fi/",
	"https://demo.talouselama.media/",
	"https://www.arvopaperi.fi/",
	"http://www.hs.fi/",
	"http://www.is.fi/taloussanomat/",
	"https://www.mtv.fi/",
	"http://www.iltalehti.fi/",
	"https://www.affarsvarlden.se/",
	"http://www.etuovi.com/",
	//"https://www.uutispuro.fi/fi",
	//"https://jelinden.fi/",
	//"http://incomewithdividends.com/",
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "assets/html/index.html")
}

func main() {
	go doEvery(schedule, checkURLs)
	checkURLs()
	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	router.GET("/", Index)
	fsStatic := util.JustFilesFilesystem{Fs: http.Dir("assets/")}
	router.Handler("GET", "/assets/*filepath", http.StripPrefix("/assets", util.GH(http.FileServer(fsStatic))))
	router.Handler("GET", "/urls", util.GH(urlHandler()))
	router.Handler("GET", "/json/*filepath", util.GH(jsonHandler()))
	http.ListenAndServe(":8080", router)
}

func checkURLs() {
	for _, url := range urls {
		go fetch(url)
		time.Sleep(5 * time.Second)
	}
}

func fetch(url string) {
	t := time.Now()
	bow := surf.NewBrowser()
	err := bow.Open(url)
	if err != nil {
		log.Println(err.Error())
		cacheAdd(url, -100*time.Millisecond, -100*time.Millisecond)
	} else {
		took := time.Now().Sub(t)
		if bow.StatusCode() != 200 {
			took = 0
		}

		t := surfTo(bow, url, t)
		cacheAdd(url, took, t)
		log.Println(bow.StatusCode(), url, took, t)
	}
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func surfTo(bow *browser.Browser, url string, t time.Time) time.Duration {

	ch := make(browser.AsyncDownloadChannel, 1)
	queue := 0

	for _, script := range bow.Scripts() {
		buf := bytes.NewBuffer([]byte{})
		script.DownloadAsync(buf, ch)
		queue++
	}

	for _, styles := range bow.Stylesheets() {
		buf := bytes.NewBuffer([]byte{})
		styles.DownloadAsync(buf, ch)
		queue++
	}

	for _, image := range bow.Images() {
		if image.Url().Scheme != "data" {
			buf := bytes.NewBuffer([]byte{})
			image.DownloadAsync(buf, ch)
			queue++
		}
	}

	for {
		select {
		case result := <-ch:
			if result.Error != nil {
				log.Printf("Error download '%s'. %s\n", result.Asset.Url(), result.Error)
			} else {
				//log.Printf("Downloaded '%s'.\n", result.Asset.Url())
			}

			queue--
			if queue == 0 {
				goto FINISHED
			}
		}
	}
FINISHED:
	close(ch)
	return time.Now().Sub(t)
}
