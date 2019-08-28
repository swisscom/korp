package patch

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// TODO to be implemented
// patch - Patch all yaml files in the given path with the new Docker image reference
func patch(c *cli.Context) error {

	patchPath := c.String("files")
	kstPath := c.String("kustomization-path")

	log.Debugf("Path to yaml files to be patched: %s", patchPath)
	log.Debugf("Kustomization file path: %s", kstPath)

	// load kustomization.yaml
	// patch all yaml in the path

	errMsg := "Method not yet implemented!"
	log.Error(errMsg)
	return errors.New(errMsg)
}
