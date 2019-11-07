package database

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	telegram "gopkg.in/tucnak/telebot.v2"

	// Importing solely for database/sql driver use
	_ "github.com/mattn/go-sqlite3"
)

// DBFile points to the location of the database file to write to
// (it will be created if it doesn't exist)
const DBFile = "bigboofer_data.sqlite3"

// MaxChallengeTime describes the maximum time to wait
// for a user to complete a challenge before removing them
// (in SQLite duration format, see documentation for datetime()).
const MaxChallengeTime = "+5 minutes"

// AddUser adds a new user and their group to the challenged users list.
func AddUser(user *telegram.User, group *telegram.Chat) {
	db := GetDB()
	defer db.Close()
	transaction, _ := db.Begin()

	_, err := db.Exec(
		"INSERT OR REPLACE INTO challenge(group_id, user_id, issued_on) "+
			"VALUES (?, ?, CURRENT_TIMESTAMP)",
		group.ID, user.ID,
	)

	if err != nil {
		log.Printf("Error in AddUser query!! %v\n", err)
		transaction.Rollback()
		return
	}

	transaction.Commit()
}

// VetUser removes a new user and their group from the challenged users list.
// (i.e., they passed a challenge or it timed out.)
func VetUser(user *telegram.User, group *telegram.Chat) {
	db := GetDB()
	defer db.Close()
	transaction, _ := db.Begin()

	_, err := db.Exec(
		"DELETE FROM challenge(group_id, user_id) "+
			"VALUES (?, ?)",
		group.ID, user.ID,
	)

	if err != nil {
		log.Printf("Error in VetUser query!! %v\n", err)
		transaction.Rollback()
		return
	}

	transaction.Commit()
}

// UserWasVetted returns true if the bot is not currently
// waiting for a challenge response from the user in the
// given group.
func UserWasVetted(user *telegram.User, group *telegram.Chat) bool {
	db := GetDB()
	defer db.Close()

	// If user is currently being vetted, COUNT(*) should return 1. (Else 0.)
	// So, if the user was already vetted here, we should expect a 0.
	var countResult int
	queryResult, err := db.Query(
		"SELECT COUNT(*) FROM challenge WHERE group_id=? AND user_id=?",
		group.ID, user.ID,
	)

	if err != nil {
		log.Printf("Error in UserWasVetted query!! Returning false. %v\n", err)
		return false
	}

	queryResult.Next()
	queryResult.Scan(&countResult)
	queryResult.Close()
	return countResult == 0
}

// SetAuthChannel sets the passphrase and channel username of the channel
// containing the passphrase for a given chat.
func SetAuthChannel(group *telegram.Chat, channelURL string, passphrase string) {
	db := GetDB()
	defer db.Close()
	transaction, _ := db.Begin()

	_, err := db.Exec(
		"INSERT OR REPLACE INTO channels (group_id, channel_url, passphrase) VALUES (?, ?, ?)",
		group.ID, channelURL, passphrase,
	)

	if err != nil {
		log.Printf("Error in SetAuthChannel query!! %v\n", err)
		transaction.Rollback()
		return
	}

	transaction.Commit()
}

// GetAuthChannel returns the channel username of the channel
// containing the passphrase for a given chat.
func GetAuthChannel(group *telegram.Chat) string {
	db := GetDB()

	var channelURL string
	queryResult, err := db.Query(
		"SELECT channel_url FROM channels WHERE group_id=?",
		group.ID,
	)

	if err != nil {
		log.Printf("Error in GetAuthChannel query!! Returning empty string. %v\n", err)
		return ""
	}

	queryResult.Next()
	queryResult.Scan(&channelURL)
	queryResult.Close()
	db.Close()

	// Prepend the t.me prefix to the username if necessary so the channel
	// username is clickable in message
	if channelURL != "" && !strings.HasPrefix(channelURL, "https://") {
		if !strings.HasPrefix(channelURL, "t.me/") {
			channelURL = "https://t.me/" + channelURL
		} else {
			channelURL = "https://" + channelURL
		}
	}
	return channelURL
}

// CheckPassphrase returns true if the passphrase given is valid
// for the given chat.
func CheckPassphrase(group *telegram.Chat, passphrase string) bool {
	db := GetDB()

	// If the passphrase matches, COUNT(*) should return 1. (Else 0.)
	var countResult int
	queryResult, err := db.Query(
		"SELECT COUNT(*) FROM channels WHERE group_id=? AND passphrase=?",
		group.ID, passphrase,
	)

	if err != nil {
		log.Printf("Error in CheckPassphrase query!! Returning false. %v\n", err)
		return false
	}

	queryResult.Next()
	queryResult.Scan(&countResult)
	queryResult.Close()
	db.Close()

	return countResult == 1
}

// PurgeOldChallengesForAllChats runs PurgeOldChallengesForChat for all chats
// we know of in the database.
func PurgeOldChallengesForAllChats(bot *telegram.Bot) {
	db := GetDB()

	var groupID int64
	var chatTargets []telegram.Chat

	queryResult, err := db.Query(
		"SELECT group_id FROM challenge",
	)

	if err != nil {
		log.Printf("Error in PurgeOldChallengesForAllChats query!! %v\n", err)
		db.Close()
		return
	}

	for queryResult.Next() {
		queryResult.Scan(&groupID)
		chatTargets = append(chatTargets, telegram.Chat{
			ID: groupID,
		})
	}

	queryResult.Close()
	db.Close()

	for _, chatTarget := range chatTargets {
		PurgeOldChallengesForChat(bot, &chatTarget)
	}
}

// PurgeOldChallengesForChat expires any challenge greater than the threshold
// (set in db_interface.go) for the given chat, and removes the users
// in the chat if they are still there. (Presumably, they haven't completed
// the challenge in time.)
func PurgeOldChallengesForChat(bot *telegram.Bot, group *telegram.Chat) {
	db := GetDB()

	var userID int
	var userTargets []telegram.ChatMember

	queryResult, err := db.Query(
		"SELECT user_id FROM challenge WHERE group_id=? "+
			"AND datetime(issued_on, ?, 'localtime') < datetime('now')",
		group.ID, MaxChallengeTime,
	)

	if err != nil {
		log.Printf("Error in PurgeOldChallengesForChat query!! %v\n", err)
		db.Close()
		return
	}

	for queryResult.Next() {
		queryResult.Scan(&userID)
		userTargets = append(userTargets, telegram.ChatMember{
			User: &telegram.User{
				ID: userID,
			},
		})
	}

	queryResult.Close()
	db.Close()

	for _, userTarget := range userTargets {
		userChat, _ := bot.ChatByID(fmt.Sprintf("%v", userTarget.User.ID))

		log.Printf(
			"Removing %v (%v) from %v (%v), expired challenge.\n",
			userChat.Username, userTarget.User.ID, group.Username, group.ID,
		)

		bot.Send(
			group,
			fmt.Sprintf("@%v didn't respond to challenge in time, removing!", userChat.Username),
		)

		// Expiring is the same as removing and vetting
		bot.Ban(group, &userTarget)
		VetUser(userTarget.User, group)
	}
}

// OnboardDB creates the sqlite3 database file if
// if doesn't already exist.
func OnboardDB() {
	log.Println("Preparing database...")
	db := GetDB()
	for _, statement := range readDDL() {
		_, err := db.Exec(statement)

		if err != nil {
			log.Println("Error onboarding database! Error details follow:")
			log.Panicln(err)
		}
	}

	db.Close()
	log.Println("Database ready!")
}

// GetDB returns a new SQL connection object.
// Since this is SQLite, we create a new connection for each
// transaction. This functionality should be changed if another
// database engine is used. Make sure to close the DB when you
// are done using it.
func GetDB() *sql.DB {
	DB, err := sql.Open("sqlite3", DBFile)

	if err != nil {
		log.Printf("Could not create a connection object. Do we have ")
		log.Printf("permission to create or read a file in the DBFile directory? ")
		log.Println("Error details follow:")
		log.Panicln(err)
	}

	return DB
}

// readDDL returns DDL in schema.sql as a list of DDL strings
func readDDL() []string {
	return strings.Split(Schema, ";")
}
