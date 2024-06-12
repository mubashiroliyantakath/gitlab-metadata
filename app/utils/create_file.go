package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func CreateFile(filename string) {

	_, err := os.Create(filename)
	if err != nil {
		log.Fatal("Failed to create output file", err)
	}
}
