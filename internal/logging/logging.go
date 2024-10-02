package logging

import (
	"log"
	"os"
)

func SetupLogging() {
	log.SetOutput(os.Stdout)
	log.Println("Logging started")
}