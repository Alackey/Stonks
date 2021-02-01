package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/spacecodewor/fmpcloud-go/objects"
)

// onMessageCreate handles Message Create events
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "$") {
		return
	}

	args := strings.Fields(m.Content[1:])
	args[0] = strings.ToLower(args[0])

	// Quote - get price
	if args[0] == "q" {
		if len(args) < 2 {
			return
		}

		symbol := strings.ToUpper(args[1])

		quote, err := stocks.Quote(symbol)
		if err != nil {
			log.Fatalf("Error getting quote for symbol $%s: %v\n", symbol, err)
			return
		}

		s.ChannelMessageSendEmbed(m.ChannelID, createQuoteMessage(symbol, quote))
	}

	// Futures - get the futures data
	if args[0] == "futures" {
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
	if args[0] == "market" {
		key := "Stock Market"
		if len(args) > 1 && strings.ToLower(args[1]) == "crypto" {
			key = "Crypto"
		}

		image, err := stocks.Market(key)
		if err != nil {
			log.Fatalln("Error getting the market heatmap:", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, image)
	}

	// News - get the latest newest for a stock
	if args[0] == "news" {
		if len(args) < 2 {
			return
		}
		symbol := strings.ToUpper(args[1])

		news, err := stocks.News(symbol)
		if err != nil {
			log.Fatalf("Error getting quote for symbol $%s: %v\n", symbol, err)
			return
		}

		if len(news) > 0 {
			s.ChannelMessageSendEmbed(m.ChannelID, createNewsMessage(symbol, news))
		}
	}

	// Help - get list of commands
	if args[0] == "help" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Stonk Bot Help",
			Description: "`$q <symbol>` - Get the price information about the stock symbol\n" +
				"`$futures` - Get the price information for some futures\n" +
				"`$market` - Get a heatmap of the market and its sectors\n" +
				"`$market crypto` - Get a heatmap of the crypto market\n" +
				"`$news <symbol>` - Get the most recent news about a stock\n" +
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
				Value:  addPlus(quote.Change),
				Inline: true,
			},
			{
				Name:   "% Change",
				Value:  addPlus(quote.ChangesPercentage) + "%",
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

// createNewsMessage creates the news message response for the news command
func createNewsMessage(symbol string, news []objects.StockNews) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for _, v := range news {
		var publishedDate string
		dateTime, err := time.Parse("2006-01-02 15:04:05", v.PublishedDate)
		if err != nil {
			log.Printf("Could not parse published date: %s\n", err.Error())
		} else {
			publishedDate = dateTime.Format("Mon Jan 2 2006, 3:04 PM")
		}

		field := &discordgo.MessageEmbedField{
			Name:  v.Title,
			Value: fmt.Sprintf("%s\n%s\n[Read More](%s)", publishedDate, v.Text, v.URL),
			
		}
		fields = append(fields, field)
	}

	return &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("%s News", symbol),
		Color:  3447003,
		Fields: fields,
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
