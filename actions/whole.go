package actions

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func Whole(filesPath, registry, patch *string, debug *bool) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		setLogLevel(debug)

		// OPTIMIZE DOCKER-CLIENT OPENING/CLOSING

		// scan

		// pull

		// push

		if *patch == "true" {
			// patch
		}

		errMsg := "Method not yet implemented!"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}
