package docker_utils

import (
	"context"
)

// TagDockerImage - Tag Docker image
func TagDockerImage(cli DockerClient, ctx *context.Context,
	imageName, imageTag, imageNameNew, imageTagNew string) error {

	imageRef := BuildCompleteDockerImage(imageName, imageTag)
	imageRefNew := BuildCompleteDockerImage(imageNameNew, imageTagNew)

	tagErr := cli.ImageTag(*ctx, imageRef, imageRefNew)
	if tagErr != nil {
		// log.Error(tagErr)
		return tagErr
	}

	return nil
}
