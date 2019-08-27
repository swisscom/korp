package main

import (
	"github.com/swisscom/korp/app"

	log "github.com/sirupsen/logrus"
)

const (
	logLevel = "info"
)

// main -
func main() {

	setLogLevel(logLevel)

	app.Create().Start()
}

// setLogLevel - Setup logrus logging level [error, warn, info, debug]
func setLogLevel(levelStr string) {

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(level)
}
