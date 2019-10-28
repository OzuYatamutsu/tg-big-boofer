package main

import (
	"bigboofer/database"
	"bigboofer/handlers"

	"log"
	"time"

	telegram "gopkg.in/tucnak/telebot.v2"
)

func main() {
	log.Println("Connecting to Telegram...")

	// Connect bot to Telegram
	bot := connectBot()

	// Set up database
	database.OnboardDB()

	// Register event handlers
	bot.Handle(telegram.OnAddedToGroup, func(message *telegram.Message) {
		handlers.OnAddedToGroupHandler(bot, message)
	})

	log.Printf("Bot %v is connected!\n", bot.Me.Username)
	bot.Start()
}

func connectBot() *telegram.Bot {
	bot, err := telegram.NewBot(telegram.Settings{
		Token:  APIKey,
		Poller: &telegram.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Printf("Could not connect to Telegram. Make sure you are ")
		log.Printf("connected to the internet and have set the API key ")
		log.Println("in config.go. Error details follow:")
		log.Panicln(err)
	}

	return bot
}
