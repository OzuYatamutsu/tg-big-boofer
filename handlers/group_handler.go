package handlers

import (
	"bigboofer/database"

	"fmt"
	"log"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnAddedToGroupHandler handles what should happen when
// the bot is newly added to a group.
func OnAddedToGroup(bot *telegram.Bot, message *telegram.Message) {
	bot.Send(message.Chat, "Woof! Woof! ▽・ω・▽")
}

// OnUserJoined handles what should happen when
// the bot sees a new user join a group it is a part of.
func OnUserJoined(bot *telegram.Bot, message *telegram.Message) {
	log.Printf(
		"New user %v (%v) in %v (%v), issuing challenge.\n",
		message.UserJoined.Username, message.UserJoined.ID,
		message.Chat.Username, message.Chat.ID,
	)

	database.AddUser(message.UserJoined, message.Chat)

	welcomeMessage := fmt.Sprintf(
		"Hello, @%v! Welcome to the group. Please read %v "+
			"and reply here with the passphrase in the channel. To prevent spam, "+
			"you will be prevented from sending messages until you do so.",
		message.UserJoined.Username,
		"#RULES_CHANNEL_URL",
	)
	bot.Send(message.Chat, welcomeMessage)
}
