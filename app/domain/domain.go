package domain

import "time"

type Domain struct {
	Url             string
	Name            string
	ScreenshotNode  string
	WaitVisibleNode string
	ScreenshotTime  string
}

type Times struct {
	IndexTime  time.Duration
	RenderTime time.Duration
}
