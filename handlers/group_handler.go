package handlers

import (
	"bigboofer/database"

	"fmt"
	"log"
	"strings"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnAddedToGroup handles what should happen when the bot is
// newly added to a group.
func OnAddedToGroup(bot *telegram.Bot, message *telegram.Message) {
	bot.Send(message.Chat, "Woof! Woof! ▽・ω・▽")
	bot.Send(
		message.Chat,
		"Admins, please promote me to admin and configure me "+
			"by running /setchannel <channel_url> <passphrase>!",
	)
}

// OnUserJoined handles what should happen when
// the bot sees a new user join a group it is a part of.
func OnUserJoined(bot *telegram.Bot, message *telegram.Message) {
	if database.GetAuthChannel(message.Chat) == "" {
		log.Printf(
			"New user %v (%v) in %v (%v), but auth channel was not set here.\n",
			message.UserJoined.Username, message.UserJoined.ID,
			message.Chat.Username, message.Chat.ID,
		)

		bot.Send(
			message.Chat,
			"Admins, please promote me to admin and configure me "+
				"by running /setchannel <channel_url> <passphrase>!",
		)
		return
	}

	log.Printf(
		"New user %v (%v) in %v (%v), issuing challenge.\n",
		message.UserJoined.Username, message.UserJoined.ID,
		message.Chat.Username, message.Chat.ID,
	)

	database.AddUser(message.UserJoined, message.Chat)
	bot.Send(message.Chat, constructVetMessage(
		message.UserJoined.Username,
		database.GetAuthChannel(message.Chat),
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
	// Check whether it matches the passphrase,
	if !strings.HasPrefix(message.Text, "/") &&
		database.CheckPassphrase(message.Chat, message.Text) {
		// Passphrase matches! Vet this user.
		database.VetUser(message.Sender, message.Chat)

		log.Printf(
			"User %v (%v) was vetted in %v (%v)",
			message.Sender.Username, message.Sender.ID,
			message.Chat.Username, message.Chat.ID,
		)

		bot.Send(
			message.Chat,
			fmt.Sprintf(
				"Woof!! Thanks, @%v! You are free to chat as you wish. ▽ - ω - ▽",
				message.Sender.Username,
			),
		)

		// Delete the message to clean up
		bot.Delete(message)

		return
	}

	// Delete the message,
	err := bot.Delete(message)

	if err != nil {
		log.Printf(
			"Could not delete message sent by unvetted %v (%v) ",
			message.Sender.Username, message.Sender.ID,
		)
		log.Printf(
			"in chat %v (%v). Do we have admin permission there?\n",
			message.Chat.Username, message.Chat.ID,
		)
	}

	// ...and PM the user.
	bot.Send(
		message.Sender, constructVetMessage(
			message.Sender.Username,
			database.GetAuthChannel(message.Chat),
		),
	)
}

// constructVetMessage returns the string to send an unvetted user
// on join or on message send before vetting.
func constructVetMessage(username string, rulesURL string) string {
	return fmt.Sprintf(
		"Hello, @%v! Welcome to the group. Please read %v "+
			"and reply with the passphrase written in the channel. "+
			"To prevent spam, you will be prevented from sending "+
			"messages until you do so. Admins, you can manually "+
			"approve this user by typing /approve @%v.",
		username,
		rulesURL,
		username,
	)
}
