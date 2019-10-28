package helpers

import (
	"os"
	"path/filepath"
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
