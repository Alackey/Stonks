package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")

	sess, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln("Error creating Discord session:", err)
		return
	}
	defer sess.Close()

	// Listen to only channel messages
	sess.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	sess.AddHandler(messageCreate)

	err = sess.Open()
	if err != nil {
		log.Fatalln("Error opening Discord connection:", err)
		os.Exit(1)
	}

	// Initialize stock client
	err = NewStocksClient()
	if err != nil {
		log.Fatalln("Error creating stocks client:", err)
		os.Exit(1)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// messageCreate handles Message Create events
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	token := "$"

	// Quote - get price
	if strings.HasPrefix(m.Content, token+"q ") {
		ticker := strings.TrimPrefix(m.Content, token+"q ")

		price, err := stocks.Quote(ticker)
		if err != nil {
			log.Fatalf("Error getting quote for ticker $%s: %v\n", ticker, err)
		}

		s.ChannelMessageSend(m.ChannelID, price)
	}

	// Market - get the market heatmap image
	if strings.TrimSpace(m.Content) == token+"market" {
		heatmap, err := stocks.Market()
		if err != nil {
			log.Fatalln("Error getting the market heatmap:", err)
		}

		s.ChannelFileSend(m.ChannelID, "marketHeatmap.png", heatmap)
	}
}
