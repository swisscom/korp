package docker_utils

import (
	"context"
	"io"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// PullDockerImage - Pull Docker image
func PullDockerImage(cli *client.Client, ctx *context.Context,
	dockerImage string, options *types.ImagePullOptions) error {

	canonicalDockerImage, normErr := NormalizeImageName(dockerImage)
	if normErr != nil {
		// log.Error(normErr)
		return normErr
	}

	pullReader, pullErr := cli.ImagePull(*ctx, canonicalDockerImage, *options)
	if pullErr != nil {
		// log.Error(pullErr)
		return pullErr
	}
	defer pullReader.Close()

	_, ioErr := io.Copy(log.StandardLogger().Writer(), pullReader)
	if ioErr != nil {
		// log.Error(ioErr)
		return ioErr
	}

	return nil
}

// NormalizeImageName - Normalize Docker image names to canonical
func NormalizeImageName(dockerImage string) (string, error) {

	named, err := reference.ParseNormalizedNamed(dockerImage)
	if err != nil {
		return "", err
	}

	return named.Name(), nil
}
