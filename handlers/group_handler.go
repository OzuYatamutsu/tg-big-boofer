package handlers

import (
	"fmt"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnAddedToGroupHandler handles what should happen when
// the bot is newly added to a group.
func OnAddedToGroupHandler(bot *telegram.Bot, message *telegram.Message) {
	bot.Send(message.Chat, "Hello! It's a warp gate!")

	// TODO: Add all users in the chat to whitelist
}

// OnUserJoined handles what should happen when
// the bot sees a new user join a group it is a part of.
func OnUserJoined(bot *telegram.Bot, message *telegram.Message) {
	welcomeMessage := fmt.Sprintf("Hello, @%v!", message.UserJoined.Username)
	bot.Send(message.Chat, welcomeMessage)

	// TODO: Add user to whitelist
}
