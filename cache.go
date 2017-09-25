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
	Service     string   `json:"service"`
	FirstLoads  []string `json:"firstLoads"`
	FinalLoads  []string `json:"finalLoads"`
	RenderLoads []string `json:"renderLoads"`
	Latest      string   `json:"latest"`
}

func cacheGet(url string) string {
	item, found := cache.Get(url)
	if found {
		return item.(string)
	}
	return ""
}

func cacheAdd(url string, service string, t time.Duration, wholeTime time.Duration, renderTime time.Duration, now *string) {
	item, found := cache.Get(url)
	var durations Duration
	var durationString string

	if found {
		json.Unmarshal([]byte(item.(string)), &durations)
		n := durations.Latest
		if now != nil {
			n = *now
		}
		if len(durations.FirstLoads) >= 60 {
			durations.FirstLoads = durations.FirstLoads[1:]
			durations.FinalLoads = durations.FinalLoads[1:]
			durations.RenderLoads = durations.RenderLoads[1:]
		}
		durations.FirstLoads = append(durations.FirstLoads, strconv.FormatFloat(t.Seconds(), 'f', 2, 64))
		durations.FinalLoads = append(durations.FinalLoads, strconv.FormatFloat(wholeTime.Seconds(), 'f', 2, 64))
		durations.RenderLoads = append(durations.RenderLoads, strconv.FormatFloat(renderTime.Seconds(), 'f', 2, 64))
		durations.Latest = n
		durationString = marshal(durations)
	} else {
		durationString = marshal(Duration{
			Service:     service,
			FirstLoads:  []string{strconv.FormatFloat(t.Seconds(), 'f', 2, 64)},
			FinalLoads:  []string{strconv.FormatFloat(wholeTime.Seconds(), 'f', 2, 64)},
			RenderLoads: []string{strconv.FormatFloat(renderTime.Seconds(), 'f', 2, 64)},
			Latest:      *now,
		})
	}
	cache.Set(url, durationString)
}

func marshal(durations Duration) string {
	d, _ := json.Marshal(durations)
	return string(d)
}
