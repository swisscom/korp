package actions

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func Patch(patchPath, kstPath *string, debug *bool) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		setLogLevel(debug)

		// loac kustomization.yaml

		// TBD

		errMsg := "Method not yet implemented!"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}
