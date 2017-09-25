package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jelinden/blig/app/util"
	"github.com/jelinden/urlmonitor/app/chrome"
	"github.com/jelinden/urlmonitor/app/domain"
	"github.com/julienschmidt/httprouter"
)

const schedule = 60 * time.Second

var domains = []domain.Domain{
	domain.Domain{Name: "www_kl", Url: "https://www.kauppalehti.fi/", ScreenshotNode: ".main-navigation__container", WaitVisibleNode: ".main-navigation__container"},
	domain.Domain{Name: "m_kl", Url: "https://m.kauppalehti.fi/", ScreenshotNode: ".navigation-main", WaitVisibleNode: ".navigation-main"},
	domain.Domain{Name: "av", Url: "https://www.affarsvarlden.se/", ScreenshotNode: "#nav", WaitVisibleNode: "#nav"},
	domain.Domain{Name: "hs", Url: "http://www.hs.fi/", ScreenshotNode: ".main-footer", WaitVisibleNode: ".main-footer"},
	domain.Domain{Name: "ts", Url: "http://www.is.fi/taloussanomat/", ScreenshotNode: ".main-footer", WaitVisibleNode: ".main-footer"},
	domain.Domain{Name: "mtv", Url: "https://www.mtv.fi/", ScreenshotNode: "#main", WaitVisibleNode: "#main"},
	domain.Domain{Name: "il", Url: "http://www.iltalehti.fi/", ScreenshotNode: "#ylanavi", WaitVisibleNode: "#ylanavi"},
	domain.Domain{Name: "ap", Url: "https://www.arvopaperi.fi/", ScreenshotNode: ".header__nav", WaitVisibleNode: ".header__nav"},
	domain.Domain{Name: "te", Url: "https://demo.talouselama.media/", ScreenshotNode: ".alma-footer", WaitVisibleNode: ".alma-footer"},
	domain.Domain{Name: "jel", Url: "https://jelinden.fi/", ScreenshotNode: ".footer", WaitVisibleNode: ".footer"},
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
	http.ListenAndServe(":8800", router)
}

func checkURLs() {
	for _, d := range domains {
		go fetch(d)
		time.Sleep(6 * time.Second)
	}
}

func fetch(d domain.Domain) {
	c := make(chan domain.Times)
	go chrome.Render(d, c)
	times := <-c
	now := time.Now().Local().Format("02.01.2006 15:04:05")
	cacheAdd(d.Url, d.Name, times.IndexTime, times.IndexTime, times.RenderTime, &now)
	log.Println(d.Url, d.Name, times.IndexTime, times.IndexTime, times.RenderTime)
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}
