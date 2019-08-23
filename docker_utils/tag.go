package docker_utils

import (
	"context"

	"github.com/docker/docker/client"
)

// TagDockerImage - Tag Docker image
func TagDockerImage(cli *client.Client, ctx *context.Context,
	imageName, imageTag, imageNameNew, imageTagNew string, normalize bool) error {

	if normalize {
		canonicalImageName, normErr := NormalizeImageName(imageName)
		if normErr != nil {
			// log.Error(normErr)
			return normErr
		}
		imageName = canonicalImageName

		canonicalImageNameNew, normNewErr := NormalizeImageName(imageNameNew)
		if normNewErr != nil {
			// log.Error(normNewErr)
			return normNewErr
		}
		imageNameNew = canonicalImageNameNew
	}

	imageRef := BuildCompleteDockerImage(imageName, imageTag)
	imageRefNew := BuildCompleteDockerImage(imageNameNew, imageTagNew)

	tagErr := cli.ImageTag(*ctx, imageRef, imageRefNew)
	if tagErr != nil {
		// log.Error(tagErr)
		return tagErr
	}

	return nil
}
