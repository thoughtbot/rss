package main

import (
	"testing"

	rss "github.com/mattn/go-pkg-rss"
)

func TestStripPodcastEpisodePrefix(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"Better Commit Messages", "Better Commit Messages"},
		{"10: Minisode 0.1.1", "Minisode 0.1.1"},
		{"", ""},
	} {
		got := stripPodcastEpisodePrefix(tt.in)

		if got != tt.want {
			t.Errorf("stripPodcastEpisodePrefix(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestGetDescription(t *testing.T) {
	itemWithoutItunesSummary := &rss.Item{
		Description: "description!",
	}
	got := getDescription(itemWithoutItunesSummary)
	want := itemWithoutItunesSummary.Description
	if got != want {
		t.Errorf("getDescription(%v) = %q; want %q", itemWithoutItunesSummary, got, want)
	}

	extensions := map[string]map[string][]rss.Extension{
		"http://www.itunes.com/dtds/podcast-1.0.dtd": map[string][]rss.Extension{
			"summary": []rss.Extension{
				{
					Value: "dude",
				},
			},
		},
	}
	itemWithItunesSummary := &rss.Item{Extensions: extensions}
	got = getDescription(itemWithItunesSummary)
	want = "dude"
	if got != want {
		t.Errorf("getDescription(%v) = %q; want %q", itemWithItunesSummary, got, want)
	}
}
