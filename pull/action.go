package pull

import (
	"context"

	"github.com/swisscom/korp/docker_utils"
	"github.com/swisscom/korp/kustomize_utils"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	kust "sigs.k8s.io/kustomize/pkg/image"
)

type Action struct {
	Io Io
}

type Io interface {
	LoadKustomizationFile(kstPath string) ([]kust.Image, error)
	OpenDockerClient() (docker_utils.DockerClient, error)
}

// pull - Pull Docker images listed in the kustomization file from remote to the local Docker registry
func (p *Action) pull(c *cli.Context) error {

	kstPath := c.String("kustomization-path")

	dockerImages, loadErr := kustomize_utils.LoadKustomizationFile(kstPath)
	if loadErr != nil {
		log.Error(loadErr)
		return loadErr
	}

	pullErr := p.pullDockerImages(dockerImages)
	if pullErr != nil {
		log.Error(pullErr)
		return pullErr
	}

	return nil
}

// pullDockerImages - Pull all Docker images from given list
func (p *Action) pullDockerImages(dockerImages []kust.Image) error {

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

		pullOk, pullKo := 0, 0
		for _, img := range dockerImages {
			if p.pullDockerImage(cli, &ctx, img.Name, img.NewTag) {
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

// pullDockerImage -
func (p *Action) pullDockerImage(cli docker_utils.DockerClient, ctx *context.Context, imageName, imageTag string) bool {

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
