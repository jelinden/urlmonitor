package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/orcaman/concurrent-map"
)

var cache cmap.ConcurrentMap

func init() {
	cache = cmap.New()
}

type Duration struct {
	FirstLoads []string `json:"firstLoads"`
	FinalLoads []string `json:"finalLoads"`
}

func cacheGet(url string) string {
	item, found := cache.Get(url)
	if found {
		return item.(string)
	}
	return ""
}

func cacheAdd(url string, t time.Duration, wholeTime time.Duration) {
	item, found := cache.Get(url)
	var durations Duration
	var durationString string
	if found {
		json.Unmarshal([]byte(item.(string)), &durations)
		if len(durations.FirstLoads) >= 60 {
			durations.FirstLoads = durations.FirstLoads[1:]
			durations.FinalLoads = durations.FinalLoads[1:]
		}
		durations.FirstLoads = append(durations.FirstLoads, strconv.FormatFloat(t.Seconds(), 'f', 2, 64))
		durations.FinalLoads = append(durations.FinalLoads, strconv.FormatFloat(wholeTime.Seconds(), 'f', 2, 64))
		durationString = marshal(durations)
	} else {
		durationString = marshal(Duration{
			[]string{strconv.FormatFloat(t.Seconds(), 'f', 2, 64)},
			[]string{strconv.FormatFloat(wholeTime.Seconds(), 'f', 2, 64)},
		})
	}
	cache.Set(url, durationString)
}

func marshal(durations Duration) string {
	d, _ := json.Marshal(durations)
	return string(d)
}
