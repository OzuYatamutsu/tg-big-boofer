package handlers

import (
	"bigboofer/database"

	"fmt"
	"log"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnAddedToGroup handles what should happen when the bot is
// newly added to a group.
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
	bot.Send(message.Chat, constructVetMessage(
		message.UserJoined.Username,
		"#RULES_CHANNEL_URL",
	))
}

// OnMessage encomposes the following events: OnText, OnPhoto, OnAudio,
// OnDocument, OnSticker, OnVideo, OnVoice, OnVideoNote, OnContact,
// OnLocation, OnVenue. If a non-vetted non-admin in a group chat
// attempts to send a message, it will be automatically deleted
// and a PM will be sent restating instructions on how to be vetted.
func OnMessage(bot *telegram.Bot, message *telegram.Message) {
	if !message.FromGroup() ||
		database.UserWasVetted(message.Sender, message.Chat) {
		// Either it was a PM, or it was a group message from someone already vetted.
		return
	}

	// Message was sent in group by non-vetted user!

	// Delete the message,
	err := bot.Delete(message)

	if err != nil {
		log.Printf(
			"Could not delete message sent by unvetted %v (%v) ",
			message.Sender.Username, message.Sender.ID,
		)
		log.Printf(
			"in chat %v (%v). Do we have admin permission there?",
			message.Chat.Username, message.Chat.ID,
		)
	}

	// ...and PM the user.
	bot.Send(
		message.Sender, constructVetMessage(
			message.Sender.Username, "#RULES_CHANNEL_URL",
		),
	)
}

// constructVetMessage returns the string to send an unvetted user
// on join or on message send before vetting.
func constructVetMessage(username string, rulesURL string) string {
	return fmt.Sprintf(
		"Hello, @%v! Welcome to the group. Please read %v "+
			"and reply here with the passphrase in the channel. To prevent spam, "+
			"you will be prevented from sending messages until you do so.",
		username,
		rulesURL,
	)
}
