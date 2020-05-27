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
	want := "from description"
	item := &rss.Item{
		Description: want,
	}
	got := getDescription(item)
	if got != want {
		t.Errorf("getDescription(%v) = %q; want %q", item, got, want)
	}

	want = "from itunes"
	item = &rss.Item{Extensions: map[string]map[string][]rss.Extension{
		"http://www.itunes.com/dtds/podcast-1.0.dtd": map[string][]rss.Extension{
			"subtitle": []rss.Extension{
				{
					Value: want,
				},
			},
		},
	},
	}
	got = getDescription(item)
	if got != want {
		t.Errorf("getDescription(%v) = %q; want %q", item, got, want)
	}
}
