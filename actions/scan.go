package actions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/swisscom/korp/docker_utils"
	"github.com/swisscom/korp/string_utils"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	kustomize "sigs.k8s.io/kustomize/pkg/image"
	"sigs.k8s.io/kustomize/pkg/types"
)

// Scan - Collect images referenced in all yaml files in the path and create a kustomization file
func Scan(scanPath, registry, output *string) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		dockerImages, dockerErr := retrieveDockerImages(scanPath, registry)
		if dockerErr != nil {
			log.Error(dockerErr)
			return dockerErr
		}

		kustomization := buildKustomization(dockerImages)

		writeErr := writeKustomizationFile(kustomization, output)
		if writeErr != nil {
			log.Error(writeErr)
			return writeErr
		}

		return nil
	}
}

// retrieveDockerImages - Retrieve all Docker images in all yaml files in the given path
func retrieveDockerImages(scanPath, registry *string) ([]kustomize.Image, error) {

	var dockerImages []kustomize.Image
	filesPaths, yamlErr := listYamlFilesPaths(*scanPath)
	if yamlErr != nil {
		// log.Error(yamlErr)
		return nil, yamlErr
	}
	for _, yamlPath := range filesPaths {
		dockerImageRefs, dockerErr := listDockerImageReferences(yamlPath)
		if dockerErr != nil {
			// log.Error(dockerErr)
			return nil, dockerErr
		}
		if len(dockerImageRefs) > 0 {
			for _, dockerImageRef := range dockerImageRefs {
				dockerImageName, dockerImageTag := docker_utils.ParseDockerImageNameAndTag(dockerImageRef)
				dockerImages = append(dockerImages, buildNewDockerImage(dockerImageName, dockerImageTag, *registry))
			}
		}
	}
	dockerImages = removeDockerImageDuplicates(dockerImages)
	log.Infof("Total Docker images found in %s: %d", *scanPath, len(dockerImages))
	return dockerImages, nil
}

// listYamlFilesPaths - List all yaml files paths in the given root path
func listYamlFilesPaths(rootPath string) ([]string, error) {

	var filesPaths []string
	yamlFileRegex := regexp.MustCompile(yamlFileRegexStr)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && yamlFileRegex.MatchString(info.Name()) {
			log.Debugf("Found yaml file: %s", path)
			filesPaths = append(filesPaths, path)
		}
		return nil
	})
	log.Debugf("Total yaml files found in %s: %d", rootPath, len(filesPaths))
	return filesPaths, err
}

// listDockerImageReferences - List all Docker image reference in the given file path
func listDockerImageReferences(filePath string) ([]string, error) {

	fileContent, err := ioutil.ReadFile(filePath)
	var dockerImagesRefs []string
	var dockerImageRefRegex = regexp.MustCompile(dockerImageRefRegexStr)
	for _, match := range dockerImageRefRegex.FindAllStringSubmatch(string(fileContent), -1) {
		if len(match) > 1 {
			ref := match[1]
			log.Debugf("Found Docker image reference: %s", ref)
			dockerImagesRefs = append(dockerImagesRefs, ref)
		}
	}
	log.Debugf("Total Docker image references found in %s: %d", filePath, len(dockerImagesRefs))
	return dockerImagesRefs, err
}

// buildNewDockerImage - Build Docker image object
func buildNewDockerImage(dockerImageName, dockerImageTag, registry string) kustomize.Image {

	trimmedDockerImageName := string_utils.TrimQuotes(dockerImageName)
	image := kustomize.Image{
		Name:    trimmedDockerImageName,
		NewName: registry + "/" + trimmedDockerImageName,
	}
	if dockerImageTag != "" {
		image.NewTag = string_utils.TrimQuotes(dockerImageTag)
	}
	return image
}

// removeDockerImageDuplicates - Remove duplicated Docker image references
func removeDockerImageDuplicates(dockerImages []kustomize.Image) []kustomize.Image {

	encountered := map[kustomize.Image]bool{}
	results := []kustomize.Image{}

	totalDuplicates := 0
	for img := range dockerImages {
		if !encountered[dockerImages[img]] {
			// Append
			results = append(results, dockerImages[img])
			// Record as encountered
			encountered[dockerImages[img]] = true
		} else {
			totalDuplicates++
		}
	}
	log.Debugf("Total Docker image duplicated references: %d", totalDuplicates)
	return results
}

// buildKustomization - Build Kustomize kustomization yaml definition object
func buildKustomization(dockerImages []kustomize.Image) *types.Kustomization {

	return &types.Kustomization{
		TypeMeta: types.TypeMeta{
			Kind:       types.KustomizationKind,
			APIVersion: types.KustomizationVersion,
		},
		Images: dockerImages,
	}
}

// writeKustomizationFile - Write Kustomize kustomization yaml object to yaml file
func writeKustomizationFile(kustomization *types.Kustomization, output *string) error {

	outputFileName, fileErr := filepath.Abs(*output + "/" + kustomizationFileName)
	if fileErr != nil {
		// log.Error(fileErr)
		return fileErr
	}

	yamlContent, yamlErr := yaml.Marshal(kustomization)
	if yamlErr != nil {
		// log.Error(yamlErr)
		return yamlErr
	}

	writeErr := ioutil.WriteFile(outputFileName, yamlContent, 0644)
	if writeErr != nil {
		// log.Error(writeErr)
		return writeErr
	}

	return nil
}
