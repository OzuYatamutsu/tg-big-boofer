package helpers

import (
	"os"
	"path/filepath"

	telegram "gopkg.in/tucnak/telebot.v2"
)

// GetRelativeProjPath will return the path to the file
// relative to the path where run.go was invoked.
func GetRelativeProjPath(pathComponents ...string) string {
	path, err := os.Executable()

	if err != nil {
		panic(err)
	}

	return filepath.Join(append(
		[]string{path}, pathComponents...,
	)...)
}

// ChatMemberContains returns true if the ChatMember list contains the User.
func ChatMemberContains(list *[]telegram.ChatMember, user *telegram.User) bool {
	for _, listUser := range *list {
		if *user == *listUser.User {
			return true
		}
	}

	return false
}
