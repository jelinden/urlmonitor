package util

import (
	"time"

	"github.com/jelinden/blig/app/domain"
)

func OnlyPublished(blogs []domain.BlogPost) []domain.BlogPost {
	publishedBlogs := []domain.BlogPost{}
	for _, blog := range blogs {
		if blog.Date.Before(time.Now()) && !blog.Date.IsZero() {
			publishedBlogs = append(publishedBlogs, blog)
		}
	}
	return publishedBlogs
}
