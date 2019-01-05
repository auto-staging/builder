package helper

import (
	"log"
	"os"
	"strconv"

	lightning "github.com/janritter/go-lightning-log"
)

// Logger contains the Lightning Logger instance configured by the Init function, it's used for logging by calling the Log function on it.
var Logger *lightning.Lightning

// Init is used to initalize Lightning Logger with the configured LogLevel.
func Init() {
	logLevel, err := strconv.Atoi(os.Getenv("CONFIGURATION_LOG_LEVEL"))
	if err != nil {
		log.Println("Init() - Couldn't convert logLevel")
		log.Println(err)
	}
	Logger, err = lightning.Init(logLevel)
	if err != nil {
		log.Println("Init() - Couldn't init Logger")
		log.Println(err)
	}
}
