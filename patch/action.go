package patch

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func patch(patchPath, kstPath *string) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		// loac kustomization.yaml

		// TBD

		errMsg := "Method not yet implemented!"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}
