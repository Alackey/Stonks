package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
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
		ticker = strings.ToUpper(ticker)

		quote, err := stocks.Quote(ticker)
		if err != nil {
			log.Fatalf("Error getting quote for ticker $%s: %v\n", ticker, err)
			return
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s (%s)\n%.2f", quote.Name, ticker, quote.Price),
			Color: 3447003,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Change",
					Value:  strconv.FormatFloat(quote.Change, 'f', 2, 64),
					Inline: true,
				},
				{
					Name:   "% Change",
					Value:  fmt.Sprintf("%.2f%%", quote.ChangesPercentage),
					Inline: true,
				},
				{
					Name:   "Volume",
					Value:  humanize.Comma(quote.Volume),
					Inline: true,
				},
				{
					Name:   "Open",
					Value:  strconv.FormatFloat(quote.Open, 'f', 2, 64),
					Inline: true,
				},
				{
					Name:   "High",
					Value:  strconv.FormatFloat(quote.DayHigh, 'f', 2, 64),
					Inline: true,
				},
				{
					Name:   "Low",
					Value:  strconv.FormatFloat(quote.DayLow, 'f', 2, 64),
					Inline: true,
				},
			},
		})
	}

	// Market - get the market heatmap image
	if strings.TrimSpace(m.Content) == token+"market" {
		heatmap, err := stocks.Market()
		if err != nil {
			log.Fatalln("Error getting the market heatmap:", err)
			return
		}

		s.ChannelFileSend(m.ChannelID, "marketHeatmap.png", heatmap)
	}

	// Help - get list of commands
	if strings.TrimSpace(m.Content) == token+"help" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Stonk Bot Help",
			Description: "`$q <symbol>` - Get the price information about the stock symbol\n" + 
						 "`$market` - Get a heatmap of the market and its sectors\n" + 
						 "`$help` - Get this help message",
			Color: 3447003,
		})
	}
}
