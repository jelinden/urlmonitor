package chrome

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gobs/args"
	"github.com/gobs/pretty"
	"github.com/gobs/simplejson"
	"github.com/jelinden/urlmonitor/app/domain"
	"github.com/orcaman/concurrent-map"
	"github.com/raff/godet"
)

var chromeapp string
var remote *godet.RemoteDebugger

const port = "localhost:9222"

type Resource struct {
	RequestID string
	URL       string
	Took      time.Duration
	Ready     bool
	Sent      time.Time
}

func init() {
	switch runtime.GOOS {
	case "darwin":
		for _, c := range []string{
			"/Applications/Google Chrome Canary.app",
		} {
			if info, err := os.Stat(c); err == nil && info.IsDir() {
				chromeapp = fmt.Sprintf("open %q --args", c)
				break
			}
		}

	case "linux":
		for _, c := range []string{
			"headless_shell",
			"chromium",
			"chromium-browser",
			"google-chrome-beta",
			"google-chrome-unstable",
			"google-chrome-stable"} {
			if _, err := exec.LookPath(c); err == nil {
				chromeapp = c
				break
			}
		}
	}
	if chromeapp != "" {
		if chromeapp == "headless_shell" {
			chromeapp += " --no-sandbox"
		} else {
			chromeapp += " --headless"
		}
		chromeapp += " --remote-debugging-port=9222 --hide-scrollbars --disable-extensions --disable-gpu about:blank"
	}

	var errConn error
	if remote, errConn = godet.Connect(port, false); errConn == nil {

	} else {
		if errRun := runCommand(chromeapp); errRun != nil {
			log.Println("cannot start browser", errRun)
		}
		var err error
		for i := 0; i < 20; i++ {
			if i > 0 {
				time.Sleep(500 * time.Millisecond)
			}

			remote, err = godet.Connect(port, false)
			if err == nil {
				break
			}
			log.Println("connect", err)
		}

		if err != nil {
			log.Println("cannot connect to browser")
		}
	}
}

func Render(d domain.Domain, c chan domain.Times) {
	t := time.Now()
	var indexTime time.Duration
	remote.AllEvents(true)

	var wg sync.WaitGroup
	var count = 0
	wg.Add(1)
	count++

	remote.CallbackEvent("Page.loadEventFired", func(params godet.Params) {
		wg.Done()
		count--
	})

	resources := cmap.New()
	remote.CallbackEvent("Network.requestWillBeSent", func(params godet.Params) {
		resourceURL := params["request"].(map[string]interface{})["url"].(string)
		requestID := params["requestId"].(string)
		if strings.Contains(resourceURL, "http") {
			if _, ok := resources.Get(requestID); !ok {
				count++
				r := Resource{RequestID: requestID, URL: resourceURL, Sent: time.Now()}
				resources.Set(requestID, r)
				wg.Add(1)
			}
		}
	})

	remote.CallbackEvent("Network.responseReceived", func(params godet.Params) {
		requestID := params["requestId"].(string)
		resource, ok := resources.Get(requestID)
		if ok {
			r := resource.(Resource)
			if r.URL == d.Url {
				indexTime = time.Now().Sub(t)
			}
		}
	})

	remote.CallbackEvent("Network.loadingFinished", func(params godet.Params) {
		requestID := params["requestId"].(string)
		if _, ok := resources.Get(requestID); ok {
			resource, _ := resources.Get(requestID)
			r := resource.(Resource)
			if !r.Ready {
				r.Ready = true
				resources.Set(requestID, r)
				wg.Done()
				count--
			}
		}
	})

	remote.CallbackEvent("Network.loadingFailed", func(params godet.Params) {
		requestID := params["requestId"].(string)
		if _, ok := resources.Get(requestID); ok {
			resource, _ := resources.Get(requestID)
			r := resource.(Resource)
			if !r.Ready {
				r.Ready = true
				resources.Set(requestID, r)
				wg.Done()
				count--
				//log.Println("COUNT DECREASED", count, allReady(resources, &wg, &count))
			}
		}
	})

	remote.Navigate(d.Url)
	for i := 0; i < 100; i++ {
		time.Sleep(300 * time.Millisecond)
		ok := allReady(resources, &wg, &count)
		if ok {
			break
		}
	}

	wg.Wait()
	waitForNode(d.WaitVisibleNode)
	log.Println(d.Url, "ALL DONE")
	remote.SaveScreenshot("assets/img/"+d.Name+".png", 0644, 0, true)
	c <- domain.Times{IndexTime: indexTime, RenderTime: time.Now().Sub(t)}
}

func allReady(resources cmap.ConcurrentMap, wg *sync.WaitGroup, count *int) bool {
	allReady := true
	for key, resource := range resources.Items() {
		r := resource.(Resource)
		if r.Ready == false {
			if time.Now().After(r.Sent.Add(2 * time.Second)) {
				r.Ready = true
				resources.Set(key, r)
				log.Println(r.URL, "TIMED OUT")
				wg.Done()
				*count--
			} else {
				allReady = false
				/*
					if len(r.URL) > 70 {
						fmt.Println(r.URL[:70], key, " not ready ")
					} else {
						fmt.Println(r.URL, key, " not ready ")
					}
				*/
			}
		}
	}
	return allReady
}

func runCommand(commandString string) error {
	parts := args.GetArgs(commandString)
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Start()
}

func waitForNode(query string) bool {
	id := documentNode(false)

	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		res, err := remote.QuerySelector(id, query)
		if err != nil {
			log.Fatal("error in querySelector: ", err)
			return false
		}
		if res != nil {
			id = int(res["nodeId"].(float64))
			res, err = remote.ResolveNode(id)
			if err != nil {
				log.Fatal("error in resolveNode: ", err)
				return false
			}
			return true
		}
	}
	return false
}

func documentNode(verbose bool) int {
	res, err := remote.GetDocument()

	if err != nil {
		log.Fatal("error getting document: ", err)
	}

	if verbose {
		pretty.PrettyPrint(res)
	}

	doc := simplejson.AsJson(res)
	return doc.GetPath("root", "nodeId").MustInt(-1)
}

func AtExit() {
	remote.Close()
	time.Sleep(5 * time.Second)
}
