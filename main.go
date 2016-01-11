package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/feeds"
	rss "github.com/jteeuwen/go-pkg-rss"
)

func main() {
	port := flag.String("port", "8080", "HTTP Port to listen on")
	flag.Parse()
	http.HandleFunc("/", rssHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func rssHandler(rw http.ResponseWriter, r *http.Request) {
	master := &feeds.Feed{
		Title:       "thoughtbot",
		Link:        &feeds.Link{Href: "https://rss.thoughtbot.com"},
		Description: "All the thoughts fit to bot.",
		Author:      &feeds.Author{Name: "thoughtbot", Email: "hello@thoughtbot.com"},
		Created:     time.Now(),
	}

	fetchFeeds(master)

	sort.Sort(byCreated(master.Items))

	result, _ := master.ToAtom()
	fmt.Fprintln(rw, result)
}

func fetchFeeds(master *feeds.Feed) {
	var wg sync.WaitGroup
	for _, feed := range sourceFeeds {
		go func(feed sourceFeed) {
			wg.Add(1)
			defer wg.Done()
			fetch(feed, master)
		}(feed)
	}
	wg.Wait()
}

func fetch(feed sourceFeed, master *feeds.Feed) {
	fetcher := rss.New(5, true, chanHandler, makeHandler(master, feed.name))
	client := &http.Client{
		Timeout: time.Second,
	}

	fetcher.FetchClient(feed.uri, client, nil)
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	// no need to do anything...
}

func makeHandler(master *feeds.Feed, sourceName string) rss.ItemHandlerFunc {
	return func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
		for i := 0; i < len(items); i++ {
			published, _ := items[i].ParsedPubDate()
			weekAgo := time.Now().AddDate(0, 0, -7)

			if published.After(weekAgo) {
				item := &feeds.Item{
					Title:       stripPodcastEpisodePrefix(items[i].Title),
					Link:        &feeds.Link{Href: items[i].Links[0].Href},
					Description: items[i].Description,
					Author:      &feeds.Author{Name: sourceName},
					Created:     published,
				}
				master.Add(item)
			}
		}
	}
}

type byCreated []*feeds.Item

type sourceFeed struct {
	uri  string
	name string
}

// sourceFeeds defines the list of thoughtbot RSS feeds
// to be combined into a master feed.
var sourceFeeds = []sourceFeed{
	{uri: "https://robots.thoughtbot.com/summaries.xml", name: "Giant Robots blog"},
	{uri: "http://simplecast.fm/podcasts/271/rss", name: "Giant Robots podcast"},
	{uri: "http://simplecast.fm/podcasts/272/rss", name: "Build Phase podcast"},
	{uri: "http://simplecast.fm/podcasts/282/rss", name: "The Bike Shed podcast"},
	{uri: "http://simplecast.fm/podcasts/1088/rss", name: "Tentative podcast"},
	{uri: "https://upcase.com/the-weekly-iteration.rss", name: "The Weekly Iteration videos"},
}

func (s byCreated) Len() int {
	return len(s)
}

func (s byCreated) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byCreated) Less(i, j int) bool {
	return s[j].Created.Before(s[i].Created)
}

var podcastEpisodePrefix = regexp.MustCompile(`^\d+: `)

func stripPodcastEpisodePrefix(s string) string {
	return podcastEpisodePrefix.ReplaceAllString(s, "")
}
