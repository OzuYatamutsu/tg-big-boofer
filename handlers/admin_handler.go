package handlers

import (
	"fmt"
	"log"
	"strings"

	"bigboofer/database"
	"bigboofer/helpers"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnApproveCommand manually approves the provided user (bypassing the passphrase check).
// Checks that the user who sent the command is an admin of the group they sent it in.
func OnApproveCommand(bot *telegram.Bot, message *telegram.Message) {
	log.Printf(
		"%v (%v) is attempting to manually approve in %v (%v)",
		message.Sender.Username, message.Sender.ID,
		message.Chat.Username, message.Chat.ID,
	)

	user := parseApproveArgs(message)

	// Validate metadata and contents
	if !validateApproveCommand(bot, message) {
		log.Printf(
			"%v (%v) failed validation for %v (%v)",
			message.Sender.Username, message.Sender.ID,
			message.Chat.Username, message.Chat.ID,
		)
		return
	}

	database.VetUser(user, message.Chat)
	log.Printf(
		"%v (%v) manually approved %v (%v) in %v (%v)",
		message.Sender.Username, message.Sender.ID,
		user.Username, user.ID,
		message.Chat.Username, message.Chat.ID,
	)
	bot.Reply(
		message, fmt.Sprintf(
			"OK!! @%v was manually approved! ▽・ω・▽",
			user.Username,
		),
	)
}

// OnSetChannelCommand sets the channel containing the passphrase (and possibly rules)
// in the current group. Checks that the user who sent the command is an admin
// of the group they sent it in.
func OnSetChannelCommand(bot *telegram.Bot, message *telegram.Message) {
	log.Printf(
		"%v (%v) is attempting to set auth channel for %v (%v)",
		message.Sender.Username, message.Sender.ID,
		message.Chat.Username, message.Chat.ID,
	)
	channelName, passphrase := parseSetChannelArgs(message)

	// Validate metadata and contents
	if !validateSetChannelCommand(bot, message) {
		log.Printf(
			"%v (%v) failed validation for %v (%v)",
			message.Sender.Username, message.Sender.ID,
			message.Chat.Username, message.Chat.ID,
		)
		return
	}

	database.SetAuthChannel(message.Chat, channelName, passphrase)
	log.Printf(
		"%v (%v) set auth channel for %v (%v): %v",
		message.Sender.Username, message.Sender.ID,
		message.Chat.Username, message.Chat.ID,
		channelName,
	)

	bot.Reply(message, "You got it, dood! Channel updated! ▽・ω・▽")
}

// validateSetChannelCommand returns true if all args are valid, returns false
// and replies with a message explaining why if not
func validateSetChannelCommand(bot *telegram.Bot, message *telegram.Message) bool {
	// Validate message contents and metadata
	if !message.FromGroup() {
		bot.Reply(message, "Please send this command from the group you wish to configure.")
		return false
	}
	// Validate that the sender is an admin of this chat
	admins, _ := bot.AdminsOf(message.Chat)
	if !helpers.ChatMemberContains(&admins, message.Sender) {
		// This person is not an admin.
		bot.Delete(message)
		return false
	}
	// Validate that channel was sent
	channelName, passphrase := parseSetChannelArgs(message)

	if channelName == "" {
		bot.Reply(
			message,
			"Please send a channel name along with your command! "+
				"(/setchannel <channel_url> <passphrase>)",
		)
		return false
	}

	if passphrase == "" {
		bot.Reply(
			message,
			"Please send a passphrase along with your command!"+
				"(/setchannel <channel_url> <passphrase>)",
		)
		return false
	}

	return true
}

// validateApproveCommand returns true if all args are valid, returns false
// and replies with a message explaining why if not
func validateApproveCommand(bot *telegram.Bot, message *telegram.Message) bool {
	// Validate message contents and metadata
	if !message.FromGroup() {
		bot.Reply(message, "Please send this command from the group you wish to configure.")
		return false
	}
	// Validate that the sender is an admin of this chat
	admins, _ := bot.AdminsOf(message.Chat)
	if !helpers.ChatMemberContains(&admins, message.Sender) {
		// This person is not an admin.
		bot.Delete(message)
		return false
	}
	// Validate that the message returned a username
	user := parseApproveArgs(message)

	if user == nil {
		bot.Reply(
			message,
			"Please send a username along with your command! "+
				"(/approve @<username>)",
		)
		return false
	}

	return true
}

// parseSetChannelArgs returns the channel name and passphrase (in that order)
// for a message relating to a /setchannel command. If one of these arguments
// was missing from the original message, returns an empty string (in the same order).
func parseSetChannelArgs(message *telegram.Message) (string, string) {
	args := strings.Split(message.Payload, " ")

	if len(args) != 2 {
		return "", ""
	}

	return args[0], args[1]
}

// parseApproveArgs returns the user that an admin is trying to approve
// for a message relating to a /setchannel command. If the command
// is missing this parameter, returns nil.
func parseApproveArgs(message *telegram.Message) *telegram.User {
	username := message.Payload

	if username == "" {
		return nil
	}

	if strings.HasPrefix(username, "@") {
		username = strings.ReplaceAll(username, "@", "")
	}

	return &telegram.User{
		ID:       database.GetIDForChallengedUsername(message.Chat, username),
		Username: username,
	}
}
