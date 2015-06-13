package main

import (
	"fmt"
	"net/http"
	"sort"
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

	fetch("https://robots.thoughtbot.com/summaries.xml", master)
	fetch("http://simplecast.fm/podcasts/271/rss", master)
	fetch("http://simplecast.fm/podcasts/272/rss", master)
	fetch("http://simplecast.fm/podcasts/282/rss", master)
	fetch("http://simplecast.fm/podcasts/1088/rss", master)

	sort.Sort(ByCreated(master.Items))

	result, _ := master.ToAtom()
	fmt.Fprintln(rw, result)
}

func fetch(uri string, master *feeds.Feed) {
	fetcher := rss.New(5, true, chanHandler, makeHandler(master))
	fetcher.Fetch(uri, nil)
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	// no need to do anything...
}

func makeHandler(master *feeds.Feed) rss.ItemHandlerFunc {
	return func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
		for i := 0; i < len(items) && i < 10; i++ {
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

type ByCreated []*feeds.Item

func (s ByCreated) Len() int {
	return len(s)
}

func (s ByCreated) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCreated) Less(i, j int) bool {
	return s[j].Created.Before(s[i].Created)
}
