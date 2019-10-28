package main

import (
	"io/ioutil"
	"log"
	"time"

	"bigboofer/database"
	"bigboofer/handlers"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// ConfigFile is the location where the API key is loaded from.
const ConfigFile = "API_KEY.config"

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
		Token:  loadAPIKey(),
		Poller: &telegram.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Printf("Could not connect to Telegram. Make sure you are ")
		log.Printf("connected to the internet and have set the API key ")
		log.Println("in API_KEY.config. Error details follow:")
		log.Panicln(err)
	}

	return bot
}

// loadAPIKey returns the API key set in API_KEY.config
func loadAPIKey() string {
	apiKey, err := ioutil.ReadFile(ConfigFile)

	if err != nil {
		log.Printf("Error reading config file. Does API_KEY.config exist ")
		log.Printf("in the project root, with the contents set to an ")
		log.Println("API key from @BotFather? Error details follow:")
		log.Panic(err)
	}

	return string(apiKey)
}
