package docker_utils

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// PushDockerImage - Push Docker image to Docker registry
func PushDockerImage(cli DockerClient, ctx *context.Context,
	imageName, imageTag string, options *types.ImagePushOptions, normalize bool) error {

	if normalize {
		canonicalImageName, normErr := NormalizeImageName(imageName)
		if normErr != nil {
			// log.Error(normErr)
			return normErr
		}
		imageName = canonicalImageName
	}

	imageRef := BuildCompleteDockerImage(imageName, imageTag)

	pushReader, pushErr := cli.ImagePush(*ctx, imageRef, *options)
	if pushErr != nil {
		// log.Error(pushErr)
		return pushErr
	}
	defer pushReader.Close()

	_, ioErr := io.Copy(log.StandardLogger().Writer(), pushReader)
	if ioErr != nil {
		// log.Error(ioErr)
		return ioErr
	}

	return nil
}
