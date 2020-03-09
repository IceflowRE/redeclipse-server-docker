package structs

import (
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirectoryExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}
