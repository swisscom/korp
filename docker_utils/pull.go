package docker_utils

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// PullDockerImage - Pull Docker image
func PullDockerImage(cli DockerClient, ctx *context.Context,
	imageName, imageTag string, options *types.ImagePullOptions, normalize bool) error {

	if normalize {
		canonicalImageName, normErr := NormalizeImageName(imageName)
		if normErr != nil {
			// log.Error(normErr)
			return normErr
		}
		imageName = canonicalImageName
	}

	imageRef := BuildCompleteDockerImage(imageName, imageTag)

	pullReader, pullErr := cli.ImagePull(*ctx, imageRef, *options)
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
