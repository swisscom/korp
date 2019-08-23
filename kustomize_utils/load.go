package kustomize_utils

import (
	"io/ioutil"
	"path/filepath"

	"github.com/swisscom/korp/korp_utils"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	kust "sigs.k8s.io/kustomize/pkg/image"
	ksttypes "sigs.k8s.io/kustomize/pkg/types"
)

// loadKustomizationFile - Load Kustomize kustomization yaml file into correspoding object
func LoadKustomizationFile(kstPath *string) ([]kust.Image, error) {

	kstFilePath, _ := filepath.Abs(*kstPath + "/" + korp_utils.KustomizationFileName)
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
