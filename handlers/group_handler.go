package handlers

import (
	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnAddedToGroupHandler handles what should happen when
// the bot is newly added to a group.
func OnAddedToGroupHandler(bot *telegram.Bot, message *telegram.Message) {
	bot.Send(message.Chat, "Hello! It's a warp gate!")

	// TODO: Add all users in the chat to whitelist
}
