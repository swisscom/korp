package actions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/swisscom/korp/utils"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	dockerimage "sigs.k8s.io/kustomize/pkg/image"
	"sigs.k8s.io/kustomize/pkg/types"
)

const (
	yamlFileRegexStr       = `(?mi).*\.(yaml|yml)`
	dockerImageRefRegexStr = `(?m)image:\s*(?P<image>[^[{\s]+)\s+`

	kustomizationFileName = "kustomization.yaml"
)

// Scan - Collect images referenced in all yaml files in the path and create a kustomization file
func Scan(scanPath, registry, output *string) func(c *cli.Context) error {

	return func(c *cli.Context) error {

		dockerImages := getDockerImages(scanPath, registry)
		kustomization := buildKustomization(dockerImages)
		return writeKustomizationFile(kustomization, output)
	}
}

// getDockerImages - Retrieve all Docker images in all yaml files in the given path
func getDockerImages(scanPath, registry *string) []dockerimage.Image {

	var dockerImages []dockerimage.Image
	filesPaths, _ := listYamlFilesPaths(*scanPath)
	for _, yamlPath := range filesPaths {
		dockerImageRefs, _ := listDockerImageReferences(yamlPath)
		if len(dockerImageRefs) > 0 {
			for _, dockerImageRef := range dockerImageRefs {
				dockerImageName, dockerImageTag := utils.GetDockerImageNameAndTag(dockerImageRef)
				dockerImages = append(dockerImages, buildNewDockerImage(dockerImageName, dockerImageTag, *registry))
			}
		}
	}
	dockerImages = removeDockerImageDuplicates(dockerImages)
	log.Debugf("Total Docker images found in %s: %d", *scanPath, len(dockerImages))
	return dockerImages
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
func buildNewDockerImage(dockerImageName, dockerImageTag, registry string) dockerimage.Image {

	trimmedDockerImageName := utils.TrimQuotes(dockerImageName)
	image := dockerimage.Image{
		Name:    trimmedDockerImageName,
		NewName: registry + "/" + trimmedDockerImageName,
	}
	if dockerImageTag != "" {
		image.NewTag = utils.TrimQuotes(dockerImageTag)
	}
	return image
}

// removeDockerImageDuplicates - Remove duplicated Docker image references
func removeDockerImageDuplicates(dockerImages []dockerimage.Image) []dockerimage.Image {

	encountered := map[dockerimage.Image]bool{}
	results := []dockerimage.Image{}

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
func buildKustomization(dockerImages []dockerimage.Image) *types.Kustomization {

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
		log.Error(fileErr)
		return fileErr
	}

	yamlContent, yamlErr := yaml.Marshal(kustomization)
	if yamlErr != nil {
		log.Error(yamlErr)
		return yamlErr
	}

	writeErr := ioutil.WriteFile(outputFileName, yamlContent, 0644)
	if writeErr != nil {
		log.Error(writeErr)
		return writeErr
	}

	return nil
}
