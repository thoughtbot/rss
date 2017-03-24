package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/feeds"
	rss "github.com/jteeuwen/go-pkg-rss"
)

func TestRSSHandler(t *testing.T) {
	var (
		recentPost = newFeedItem("Recent Post", 0)
		olderPost  = newFeedItem("Older Post", 3)
		oldPost    = newFeedItem("Old Post", 30)

		blogFeed = newFeed("Blog", recentPost, olderPost, oldPost)

		// podcast titles are prefixed with an episode number
		podcastTitle        = "We Record Things"
		podcastEpisodeTitle = "12: " + podcastTitle
		recentPodcast       = newFeedItem(podcastEpisodeTitle, 1)

		podcastFeed = newFeed("Podcast", recentPodcast)
	)

	feedServer := newFeedServer(map[string]*feeds.Feed{
		"/blog":    blogFeed,
		"/podcast": podcastFeed,
	})
	defer feedServer.Close()

	server := httptest.NewServer(rssHandler([]sourceFeed{
		{uri: feedServer.URL + "/blog", name: "Blog"},
		{uri: feedServer.URL + "/podcast", name: "Podcast"},
	}))
	defer server.Close()

	feed := rss.NewWithHandlers(0, false, nil, nil)
	if err := feed.Fetch(server.URL, nil); err != nil {
		t.Fatalf("failed to fetch feed: %s", err)
	}

	if got, want := len(feed.Channels), 1; got != want {
		t.Fatalf("len(feed.Channels) = %d, want %d", got, want)
	}

	channel := feed.Channels[0]

	if got, want := len(channel.Items), 3; got != want {
		t.Fatalf("len(channel.Items) = %d, want %d", got, want)
	}

	if got, want := channel.Title, "thoughtbot"; got != want {
		t.Errorf("channel.Title = %q, want %q", got, want)
	}

	if got, want := channel.Items[0].Title, recentPost.Title; got != want {
		t.Errorf("channel.Items[0].Title = %q, want %q", got, want)
	}

	if got, want := channel.Items[1].Title, podcastTitle; got != want {
		t.Errorf("channel.Items[1].Title = %q, want %q", got, want)
	}

	if got, want := channel.Items[2].Title, olderPost.Title; got != want {
		t.Errorf("channel.Items[2].Title = %q, want %q", got, want)
	}
}

func newFeed(title string, items ...*feeds.Item) *feeds.Feed {
	return &feeds.Feed{
		Title: title,
		Link:  &feeds.Link{Href: "http://example.com/feed"},
		Items: items,
	}
}

func newFeedItem(title string, ageInDays int) *feeds.Item {
	return &feeds.Item{
		Title:       title,
		Link:        &feeds.Link{Href: "http://example.com"},
		Description: title,
		Updated:     time.Now().AddDate(0, 0, -ageInDays),
	}
}

func newFeedServer(routes map[string]*feeds.Feed) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed, ok := routes[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err := feed.WriteRss(w)
		if err != nil {
			log.Fatal(err)
		}
	}))
}
