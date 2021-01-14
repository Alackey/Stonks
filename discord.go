package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/spacecodewor/fmpcloud-go/objects"
)

// onMessageCreate handles Message Create events
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	token := "$"

	// Quote - get price
	if strings.HasPrefix(m.Content, token+"q ") {
		symbol := strings.TrimPrefix(m.Content, token+"q ")
		symbol = strings.ToUpper(symbol)

		quote, err := stocks.Quote(symbol)
		if err != nil {
			log.Fatalf("Error getting quote for symbol $%s: %v\n", symbol, err)
			return
		}

		s.ChannelMessageSendEmbed(m.ChannelID, createQuoteMessage(symbol, quote))
	}

	// Futures - get the futures data
	if strings.TrimSpace(m.Content) == token+"futures" {
		quote, err := stocks.Futures()
		if err != nil {
			log.Fatalln("Error getting futures:", err)
			return
		}

		esChange := addPlus(quote.ES.ChangeInDouble)
		esPercentChange := addPlus(quote.ES.FuturePercentChange * 100)

		nqChange := addPlus(quote.NQ.ChangeInDouble)
		nqPercentChange := addPlus(quote.NQ.FuturePercentChange * 100)

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Futures",
			Color: 3447003,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "S&P 500",
					Value: fmt.Sprintf("%s (%s%%)", esChange, esPercentChange),
				},
				{
					Name:  "NASDAQ",
					Value: fmt.Sprintf("%s (%s%%)", nqChange, nqPercentChange),
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

// createQuoteMessage creates the quote message response for the quote command
func createQuoteMessage(symbol string, quote objects.StockQuote) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (%s)\n%.2f", quote.Name, symbol, quote.Price),
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
	}
}

// addPlus converts a float64 to string and adds a "+" if it is a positive number
func addPlus(num float64) string {
	result := strconv.FormatFloat(num, 'f', 2, 64)
	if num > 0 {
		return "+" + result
	}
	return result
}
