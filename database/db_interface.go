package database

import (
	"database/sql"
	"log"
	"strings"

	telegram "gopkg.in/tucnak/telebot.v2"

	// Importing solely for database/sql driver use
	_ "github.com/mattn/go-sqlite3"
)

// DBFile points to the location of the database file to write to
// (it will be created if it doesn't exist)
const DBFile = "bigboofer_data.sqlite3"

// AddUser adds a new user and their group to the challenged users list.
func AddUser(user *telegram.User, group *telegram.Chat) {
	db := GetDB()
	transaction, _ := db.Begin()

	db.Exec(
		"INSERT INTO challenge(group_id, user_id, issued_on) "+
			"VALUES (?, ?, CURRENT_TIMESTAMP)",
		group.ID, user.ID,
	)

	transaction.Commit()
}

// VetUser removes a new user and their group from the challenged users list.
// (i.e., they passed a challenge or it timed out.)
func VetUser(user *telegram.User, group *telegram.Chat) {
	db := GetDB()
	transaction, _ := db.Begin()

	db.Exec(
		"DELETE FROM challenge(group_id, user_id) "+
			"VALUES (?, ?)",
		group.ID, user.ID,
	)

	transaction.Commit()
}

// UserWasVetted returns true if the bot is not currently
// waiting for a challenge response from the user in the
// given group.
func UserWasVetted(user *telegram.User, group *telegram.Chat) bool {
	db := GetDB()

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
	return countResult == 0
}

// GetAuthChannel returns the channel username of the channel
// containing the passphrase for a given chat.
func GetAuthChannel(group *telegram.Chat) string {
	db := GetDB()

	var channelUsername string
	queryResult, err := db.Query(
		"SELECT channel_username FROM channels WHERE group_id=?",
		group.ID,
	)

	if err != nil {
		log.Printf("Error in GetAuthChannel query!! Returning empty string. %v\n", err)
		return ""
	}

	queryResult.Next()
	queryResult.Scan(&channelUsername)

	// Prepend the @ prefix to the username if necessary so the channel
	// username is clickable in message
	if channelUsername != "" && !strings.HasPrefix(channelUsername, "@") {
		channelUsername = "@" + channelUsername
	}
	return channelUsername
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
	return countResult == 1
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

	log.Println("Database ready!")
}

// GetDB returns a new SQL connection object.
// Since this is SQLite, we create a new connection for each
// transaction. This functionality should be changed if another
// database engine is used.
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
