package utils

import (
	"log"
)

func PrintError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
