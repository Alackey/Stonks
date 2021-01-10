package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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

	sess.AddHandler(onMessageCreate)

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
