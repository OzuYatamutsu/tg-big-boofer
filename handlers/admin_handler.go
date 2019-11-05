package handlers

import (
	"log"
	"strings"

	"bigboofer/database"
	"bigboofer/helpers"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// OnApproveCommand manually approves the provided user (bypassing the passphrase check).
// Checks that the user who sent the command is an admin of the group they sent it in.
func OnApproveCommand(bot *telegram.Bot, message *telegram.Message) {
	return // TODO
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
	channelName, passphrase := parseSetChannelArgs(message.Text, bot.Me.Username)

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
		bot.Reply(
			message,
			"It looks like you aren't an admin of this group, "+
				"so you can't configure the channel for it ▽ - ω - ▽.",
		)
		return false
	}
	// Validate that channel was sent
	channelName, passphrase := parseSetChannelArgs(message.Text, bot.Me.Username)

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

// parseSetChannelArgs returns the channel name and passphrase (in that order)
// for a raw message relating to a /setchannel command. If one of these arguments
// was missing from the original message, returns an empty string (in the same order).
func parseSetChannelArgs(rawMessage string, botUsername string) (string, string) {
	stringArgs := strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(rawMessage, "/setchannel", ""),
			"@"+botUsername,
			"",
		),
	)

	args := strings.Split(stringArgs, " ")

	if len(args) != 2 {
		return "", ""
	}

	return args[0], args[1]
}
