package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	rss "github.com/jteeuwen/go-pkg-rss"
)

func main() {
	http.HandleFunc("/", rssHandler)
	http.ListenAndServe(":8080", nil)
}

func rssHandler(rw http.ResponseWriter, r *http.Request) {
	master := &feeds.Feed{
		Title:       "thoughtbot",
		Link:        &feeds.Link{Href: "https://rss.thoughtbot.com"},
		Description: "All the thoughts fit to bot.",
		Author:      &feeds.Author{"thoughtbot", "hello@thoughtbot.com"},
		Created:     time.Now(),
	}

	blog := rss.New(5, true, chanHandler, makeHandler(master))
	blog.Fetch("https://robots.thoughtbot.com/summaries.xml", nil)
	podcast := rss.New(5, true, chanHandler, makeHandler(master))
	podcast.Fetch("http://simplecast.fm/podcasts/271/rss", nil)

	result, _ := master.ToAtom()
	fmt.Fprintln(rw, result)
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	// no need to do anything...
}

func makeHandler(master *feeds.Feed) rss.ItemHandlerFunc {
	return func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
		for i := 0; i < len(items); i++ {
			published, _ := items[i].ParsedPubDate()

			item := &feeds.Item{
				Title:       items[i].Title,
				Link:        &feeds.Link{Href: items[i].Links[0].Href},
				Description: items[i].Description,
				Author:      &feeds.Author{Name: items[i].Author.Name},
				Created:     published,
			}
			master.Add(item)
		}
	}
}
