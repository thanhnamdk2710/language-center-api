package utils

import (
	"log"
)

func Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func Error(msg string, err error) {
	log.Printf("[ERROR] %s: %v\n", msg, err)
}
