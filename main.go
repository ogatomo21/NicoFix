package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var niconicoURL = regexp.MustCompile(`https?://(?:(?:www|sp)\.)?nicovideo\.jp/watch/[A-Za-z]{2}[0-9]+(?:\?[^\s#<>"']*)?(?:#[^\s<>"']*)?`)

func fixedURLs(content string) []string {
	seen := make(map[string]struct{})
	var urls []string

	for _, url := range niconicoURL.FindAllString(content, -1) {
		url = strings.Replace(url, "nicovideo.jp", "nicovideo.gay", 1)
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		urls = append(urls, url)
	}
	return urls
}

func suppressEmbeds(flags discordgo.MessageFlags) discordgo.MessageFlags {
	return flags | discordgo.MessageFlagsSuppressEmbeds
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	urls := fixedURLs(m.Content)
	if len(urls) == 0 {
		return
	}

	_, err := s.ChannelMessageSendReply(m.ChannelID, strings.Join(urls, "\n"), &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	})
	if err != nil {
		log.Printf("reply to message %s: %v", m.ID, err)
		return
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      m.ID,
		Channel: m.ChannelID,
		Flags:   suppressEmbeds(m.Flags),
	})
	if err != nil {
		log.Printf("suppress embeds for message %s: %v", m.ID, err)
	}
}

func main() {
	token := strings.TrimSpace(os.Getenv("DISCORD_TOKEN"))
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create Discord session: %v", err)
	}
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	session.AddHandler(messageCreate)

	if err := session.Open(); err != nil {
		log.Fatalf("open Discord session: %v", err)
	}
	if err := session.UpdateGameStatus(0, ".jp → .gay"); err != nil {
		log.Printf("failed to set presence: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("close Discord session: %v", err)
		}
	}()

	log.Print("bot is running")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Print("shutting down")
}
