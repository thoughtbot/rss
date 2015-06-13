package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func main() {
	http.HandleFunc("/", rssHandler)
	http.ListenAndServe(":8080", nil)
}

func rssHandler(rw http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title:       "thoughtbot",
		Link:        &feeds.Link{Href: "https://rss.thoughtbot.com"},
		Description: "All the thoughts fit to bot.",
		Author:      &feeds.Author{"thoughtbot", "hello@thoughtbot.com"},
		Created:     time.Now(),
	}

	item := &feeds.Item{
		Title:       "HTTP Safety Doesn't Happen by Accident",
		Link:        &feeds.Link{Href: "https://robots.thoughtbot.com/http-safety-doesnt-happen-by-accident"},
		Description: "What are safe and unsafe HTTP methods, and why does it matter?",
		Author:      &feeds.Author{Name: "George Brocklehurst"},
	}
	feed.Add(item)

	result, _ := feed.ToAtom()
	fmt.Fprintln(rw, result)
}
