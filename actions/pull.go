package actions

import (
	"context"

	"github.com/swisscom/korp/docker_utils"
	"github.com/swisscom/korp/kustomize_utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	kust "sigs.k8s.io/kustomize/pkg/image"
)

// Pull - Pull Docker images listed in the kustomization file from remote to the local Docker registry
func Pull(kstPath *string) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		dockerImages, loadErr := kustomize_utils.LoadKustomizationFile(kstPath)
		if loadErr != nil {
			log.Error(loadErr)
			return loadErr
		}

		pullErr := pullDockerImages(dockerImages)
		if pullErr != nil {
			log.Error(pullErr)
			return pullErr
		}

		return nil
	}
}

// pullDockerImages - Pull all Docker images from given list
func pullDockerImages(dockerImages []kust.Image) error {

	if len(dockerImages) > 0 {

		ctx := context.Background()

		cli, cliErr := docker_utils.OpenDockerClient()
		if cliErr != nil {
			// log.Error(cliErr)
			return cliErr
		}
		defer cli.Close()

		daemonErr := checkDockerDaemon(cli, &ctx)
		if daemonErr != nil {
			// log.Error(daemonErr)
			return daemonErr
		}

		pullOk, pullKo := 0, 0
		for _, img := range dockerImages {
			if pull(cli, &ctx, img.Name, img.NewTag) {
				pullOk++
			} else {
				pullKo++
			}
		}
		log.Infof("Total Docker images pulled: %d - Total Docker images pulls failed: %d", pullOk, pullKo)
	} else {
		log.Warn("No Docker images to pull")
	}

	return nil
}

// pull -
func pull(cli *client.Client, ctx *context.Context, imageName, imageTag string) bool {

	imageRef := docker_utils.BuildCompleteDockerImage(imageName, imageTag)
	pullErr := docker_utils.PullDockerImage(cli, ctx, imageName, imageTag, &types.ImagePullOptions{}, true)
	if pullErr != nil {
		log.Errorf("Error pulling Docker image %s: %s", imageRef, pullErr.Error())
		return false
	} else {
		log.Infof("%s image pulled", imageRef)
		return true
	}
}
