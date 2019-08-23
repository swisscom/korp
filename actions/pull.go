package actions

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/swisscom/korp/docker_utils"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	dockerimage "sigs.k8s.io/kustomize/pkg/image"
	ksttypes "sigs.k8s.io/kustomize/pkg/types"
)

// Pull - Pull Docker images listed in the kustomization file to the local Docker registry
func Pull(kstPath *string) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		dockerImages, loadErr := loadKustomizationFile(kstPath)
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

// loadKustomizationFile - Load Kustomize kustomization yaml file into correspoding object
func loadKustomizationFile(kstPath *string) ([]dockerimage.Image, error) {

	kstFilePath, _ := filepath.Abs(*kstPath + "/" + kustomizationFileName)
	kstYaml, readErr := ioutil.ReadFile(kstFilePath)
	if readErr != nil {
		// log.Error(readErr)
		return nil, readErr
	}

	var kustomization ksttypes.Kustomization
	yamlErr := yaml.Unmarshal(kstYaml, &kustomization)
	if yamlErr != nil {
		// log.Error(yamlErr)
		return nil, yamlErr
	}

	log.Debugf("Total Docker image to pull: %d", len(kustomization.Images))
	return kustomization.Images, nil
}

// pullDockerImages - Pull all Docker images from given list
func pullDockerImages(dockerImages []dockerimage.Image) error {

	if len(dockerImages) > 0 {

		cli, cliErr := docker_utils.OpenDockerClient()
		if cliErr != nil {
			// log.Error(cliErr)
			return cliErr
		}
		defer cli.Close()

		ctx := context.Background()
		if !docker_utils.CheckDockerDaemonRunning(cli, &ctx) {
			log.Warn("Docker daemon NOT RUNNING")
			docker_utils.StartDockerDaemon(cli, &ctx)
		}

		pullOk := 0
		pullKo := 0
		for _, img := range dockerImages {
			pullErr := docker_utils.
				PullDockerImage(
					cli, &ctx,
					docker_utils.BuildCompleteDockerImage(img.Name, img.NewTag),
					&types.ImagePullOptions{})
			if pullErr != nil {
				log.Errorf("Error pulling Docker image %s: %s", img.Name, pullErr.Error())
				pullKo++
			} else {
				log.Infof("%s image pulled", img.Name)
				pullOk++
			}
		}
		log.Infof("Total Docker images pulled: %d - Total Docker images pulls failed: %d", pullOk, pullKo)
	} else {
		log.Warn("No Docker images to pull")
	}

	return nil
}
