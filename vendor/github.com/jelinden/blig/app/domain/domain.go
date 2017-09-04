package domain

import (
	"html/template"
	"time"
)

type Blog struct {
	BlogName  string
	BlogPosts []BlogPost
}

type BlogItem struct {
	BlogName string
	BlogItem BlogPost
}

type BlogPost struct {
	ID        string
	Title     string
	Slug      string
	Post      template.HTML
	Markdown  string
	Date      time.Time
	Modified  time.Time
	Published bool
}

type TimeSlice []BlogPost

func (p TimeSlice) Len() int {
	return len(p)
}

func (p TimeSlice) Less(i, j int) bool {
	return p[i].Date.After(p[j].Date)
}

func (p TimeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type User struct {
	Username string
	Password string
}
