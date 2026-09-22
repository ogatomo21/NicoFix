package main

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestFixedURLs(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "all supported hosts and schemes",
			content: "http://nicovideo.jp/watch/sm1 https://nicovideo.jp/watch/sm2 http://www.nicovideo.jp/watch/sm3 https://www.nicovideo.jp/watch/sm4",
			want: []string{
				"http://nicovideo.gay/watch/sm1",
				"https://nicovideo.gay/watch/sm2",
				"http://www.nicovideo.gay/watch/sm3",
				"https://www.nicovideo.gay/watch/sm4",
			},
		},
		{
			name:    "query and fragment are preserved",
			content: "https://www.nicovideo.jp/watch/sm12345678?ref=test#part",
			want:    []string{"https://www.nicovideo.gay/watch/sm12345678?ref=test#part"},
		},
		{
			name:    "duplicates are removed in first appearance order",
			content: "https://nicovideo.jp/watch/sm1 https://nicovideo.jp/watch/sm1 https://nicovideo.jp/watch/sm2",
			want:    []string{"https://nicovideo.gay/watch/sm1", "https://nicovideo.gay/watch/sm2"},
		},
		{
			name:    "unsupported URLs are ignored",
			content: "https://nicovideo.gay/watch/sm1 https://nicovideo.jp/watch/sm https://example.com/watch/sm1",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixedURLs(tt.content); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("fixedURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuppressEmbedsPreservesExistingFlags(t *testing.T) {
	existing := discordgo.MessageFlagsCrossPosted | discordgo.MessageFlagsUrgent
	want := existing | discordgo.MessageFlagsSuppressEmbeds
	if got := suppressEmbeds(existing); got != want {
		t.Fatalf("suppressEmbeds() = %b, want %b", got, want)
	}
}
