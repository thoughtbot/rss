package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"time"

	"github.com/gorilla/feeds"
	rss "github.com/mattn/go-pkg-rss"
)

var podcastEpisodePrefix = regexp.MustCompile(`^\d+: `)

var sourceFeeds = []sourceFeed{
	{uri: "https://robots.thoughtbot.com/summaries.xml", name: "Giant Robots blog"},
	{uri: "https://feeds.simplecast.com/KARThxOK", name: "Giant Robots podcast"},
	{uri: "https://feeds.simplecast.com/ky3kewHN", name: "The Bike Shed podcast"},
	{uri: "https://feeds.simplecast.com/ZBfsoMJW", name: "Tentative podcast"},
	{uri: "https://thoughtbot.com/upcase/the-weekly-iteration.rss", name: "The Weekly Iteration videos"},
	{uri: "https://hub.thoughtbot.com/releases.atom", name: "Open source software releases"},
}

func main() {
	port := flag.String("port", "8080", "HTTP Port to listen on")
	flag.Parse()
	http.Handle("/", rssHandler(sourceFeeds))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func rssHandler(sourceFeeds []sourceFeed) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Forwarded-Proto") == "http" {
			destination := *req.URL
			destination.Host = req.Host
			destination.Scheme = "https"
			http.Redirect(w, req, destination.String(), http.StatusFound)
			return
		}

		master := &feeds.Feed{
			Title:       "thoughtbot",
			Link:        &feeds.Link{Href: "https://rss.thoughtbot.com"},
			Description: "All the thoughts fit to bot.",
			Author:      &feeds.Author{Name: "thoughtbot", Email: "hello@thoughtbot.com"},
			Created:     time.Now(),
		}

		for _, feed := range sourceFeeds {
			fetch(feed, master)
		}

		sort.Sort(byCreated(master.Items))

		result, err := master.ToAtom()
		if err != nil {
			log.Printf("error generating feed: %v", err)
			http.Error(w, "error generating feed", http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintln(w, result)
		if err != nil {
			log.Printf("error printing feed: %v", err)
			http.Error(w, "error printing feed", http.StatusInternalServerError)
			return
		}
	})
}

func fetch(feed sourceFeed, master *feeds.Feed) {
	fetcher := rss.New(5, true, chanHandler, makeHandler(master, feed.name))
	client := &http.Client{
		Timeout: time.Second,
	}

	err := fetcher.FetchClient(feed.uri, client, nil)
	if err != nil {
		log.Printf("error fetching feed: %v", err)
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	// no need to do anything...
}

func makeHandler(master *feeds.Feed, sourceName string) rss.ItemHandlerFunc {
	return func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
		for i := 0; i < len(items); i++ {
			published, err := items[i].ParsedPubDate()
			if err != nil {
				log.Printf("error parsing publication date: %v", err)
				continue
			}

			weekAgo := time.Now().AddDate(0, 0, -7)

			if published.After(weekAgo) {
				item := &feeds.Item{
					Title:       stripPodcastEpisodePrefix(items[i].Title),
					Link:        &feeds.Link{Href: items[i].Links[0].Href},
					Description: getDescription(items[i]),
					Author:      &feeds.Author{Name: sourceName},
					Created:     published,
				}
				master.Add(item)
			}
		}
	}
}

type byCreated []*feeds.Item

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

func getDescription(item *rss.Item) string {
	if ext, ok := item.Extensions["http://www.itunes.com/dtds/podcast-1.0.dtd"]; ok {
		return ext["subtitle"][0].Value
	}

	return item.Description
}

type sourceFeed struct {
	uri  string
	name string
}
