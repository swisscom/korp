package korp_utils

import (
	log "github.com/sirupsen/logrus"
)

const (
	logLevelDefaultKey = "default"
	logLevel           = "info"
)

// SetLogLevel - Setup logrus logging level [error, warn, info, debug]
func SetLogLevel(levelStr string) {

	if levelStr == "" || levelStr == logLevelDefaultKey {
		levelStr = logLevel
	}

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(level)
}
