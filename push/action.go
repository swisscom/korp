package push

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

// push - Push Docker images listed in the kustomization file to the new Docker registry
func push(kstPath *string) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		dockerImages, loadErr := kustomize_utils.LoadKustomizationFile(kstPath)
		if loadErr != nil {
			log.Error(loadErr)
			return loadErr
		}

		tagPushErr := tagAndPushDockerImages(dockerImages)
		if tagPushErr != nil {
			log.Error(tagPushErr)
			return tagPushErr
		}

		return nil
	}
}

// tagAndPushDockerImages - Tag and push all Docker images from given list
func tagAndPushDockerImages(dockerImages []kust.Image) error {

	if len(dockerImages) > 0 {

		ctx := context.Background()

		cli, cliErr := docker_utils.OpenDockerClient()
		if cliErr != nil {
			// log.Error(cliErr)
			return cliErr
		}
		defer cli.Close()

		daemonErr := docker_utils.CheckDockerDaemon(cli, &ctx)
		if daemonErr != nil {
			// log.Error(daemonErr)
			return daemonErr
		}

		tagOk, tagKo, pushOk, pushKo := 0, 0, 0, 0
		for _, img := range dockerImages {
			tagResult, pushResult := tagDockerImage(cli, &ctx, img.Name, img.NewTag, img.NewName, img.NewTag)
			if tagResult {
				tagOk++
			} else {
				tagKo++
			}
			if pushResult {
				pushOk++
			} else {
				pushKo++
			}
		}
		log.Infof("Total Docker images tagged: %d - Total Docker images tag failed: %d", tagOk, tagKo)
		log.Infof("Total Docker images pushed: %d - Total Docker images push failed: %d", pushOk, pushKo)
	} else {
		log.Warn("No Docker images to tag & push")
	}

	return nil
}

// tagDockerImage -
func tagDockerImage(cli *client.Client, ctx *context.Context, imageName, imageTag, imageNameNew, imageTagNew string) (tagResult, pushResult bool) {

	imageRef := docker_utils.BuildCompleteDockerImage(imageName, imageTag)
	imageRefNew := docker_utils.BuildCompleteDockerImage(imageName, imageTag)
	tagErr := docker_utils.TagDockerImage(cli, ctx, imageName, imageTag, imageNameNew, imageTagNew, true)
	if tagErr != nil {
		log.Errorf("Error tagging Docker image %s to %s: %s", imageRef, imageRefNew, tagErr.Error())
		return false, false
	}

	log.Infof("Tag %s created", imageRefNew)

	if pushDockerImage(cli, ctx, imageNameNew, imageTagNew) {
		return true, true
	}

	return true, false
}

// pushDockerImage -
func pushDockerImage(cli *client.Client, ctx *context.Context, imageName, imageTag string) bool {

	// PLEASE NOTE: this is a required trick even with fake auth
	pushOpt := types.ImagePushOptions{
		All:          true,
		RegistryAuth: "123",
	}

	imageRef := docker_utils.BuildCompleteDockerImage(imageName, imageTag)
	pushErr := docker_utils.PushDockerImage(cli, ctx, imageName, imageTag, &pushOpt, false)
	if pushErr != nil {
		log.Errorf("Error pushing Docker image %s: %s", imageRef, pushErr.Error())
		return false
	}

	log.Infof("%s image pushed", imageRef)
	return true
}
