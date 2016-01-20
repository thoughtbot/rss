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

var podcastEpisodePrefix = regexp.MustCompile(`^\d+: `)

var sourceFeeds = []sourceFeed{
	{uri: "https://robots.thoughtbot.com/summaries.xml", name: "Giant Robots blog"},
	{uri: "http://simplecast.fm/podcasts/271/rss", name: "Giant Robots podcast"},
	{uri: "http://simplecast.fm/podcasts/272/rss", name: "Build Phase podcast"},
	{uri: "http://simplecast.fm/podcasts/282/rss", name: "The Bike Shed podcast"},
	{uri: "http://simplecast.fm/podcasts/1088/rss", name: "Tentative podcast"},
	{uri: "https://upcase.com/the-weekly-iteration.rss", name: "The Weekly Iteration videos"},
}

func main() {
	port := flag.String("port", "8080", "HTTP Port to listen on")
	flag.Parse()
	http.HandleFunc("/", rssHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func rssHandler(rw http.ResponseWriter, r *http.Request) {
	jobs := fetch()
	atomFeed := merge(jobs)
	fmt.Fprintln(rw, atomFeed)
}

func fetch() chan *feeds.Item {
	var wg sync.WaitGroup
	jobs := make(chan *feeds.Item, len(sourceFeeds))
	client := &http.Client{
		Timeout: time.Second,
	}
	for _, feed := range sourceFeeds {
		wg.Add(1)
		go func(feed sourceFeed) {
			defer wg.Done()
			fetcher := rss.New(1, true, chanHandler, makeHandler(feed.name, jobs))
			fetcher.FetchClient(feed.uri, client, nil)
		}(feed)
	}
	wg.Wait()
	return jobs
}

func merge(jobs chan *feeds.Item) string {
	master := &feeds.Feed{
		Title:       "thoughtbot",
		Link:        &feeds.Link{Href: "http://rss.thoughtbot.com"},
		Description: "All the thoughts fit to bot.",
		Author:      &feeds.Author{Name: "thoughtbot", Email: "hello@thoughtbot.com"},
		Created:     time.Now(),
	}
	deadline := time.After(1 * time.Second)
	for i := 0; i < len(sourceFeeds); i++ {
		select {
		case rssItem := <-jobs:
			master.Add(rssItem)
		case <-deadline:
			// abort
		}
	}
	sort.Sort(byCreated(master.Items))
	atomFeed, _ := master.ToAtom()
	return atomFeed
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	// no need to do anything...
}

func makeHandler(sourceName string, jobs chan *feeds.Item) rss.ItemHandlerFunc {
	return func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
		for i := 0; i < len(items); i++ {
			published, _ := items[i].ParsedPubDate()
			weekAgo := time.Now().AddDate(0, 0, -7)

			if published.After(weekAgo) {
				jobs <- &feeds.Item{
					Title:       stripPodcastEpisodePrefix(items[i].Title),
					Link:        &feeds.Link{Href: items[i].Links[0].Href},
					Description: items[i].Description,
					Author:      &feeds.Author{Name: sourceName},
					Created:     published,
				}
			}
		}
	}
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

func stripPodcastEpisodePrefix(s string) string {
	return podcastEpisodePrefix.ReplaceAllString(s, "")
}

type byCreated []*feeds.Item

type sourceFeed struct {
	uri  string
	name string
}
