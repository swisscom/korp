package scan

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/swisscom/korp/docker_utils"
	"github.com/swisscom/korp/korp_utils"
	"github.com/swisscom/korp/string_utils"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	kustomize "sigs.k8s.io/kustomize/pkg/image"
	"sigs.k8s.io/kustomize/pkg/types"
)

const (
	dockerImageRefRegexStr = `(?m)image:\s*(?P<image>[^[{\s]+)\s+`
)

// Action - struct for scan action
type Action struct {
	Io Io
}

//go:generate moq -out mocks/io.go -pkg mocks . Io

// Io - interface for all io functions used by scan
type Io interface {
	ListYamlFilesPaths(rootPath string) ([]string, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// IoImpl - real io implementation using ioutil and file_utils
type IoImpl struct {
	listYamlFilesPaths func(rootPath string) ([]string, error)
	readFile           func(filename string) ([]byte, error)
	writeFile          func(filename string, data []byte, perm os.FileMode) error
}

// ListYamlFilesPaths - real implementation of ListYamlFilesPaths backed by file_utils
func (s IoImpl) ListYamlFilesPaths(rootPath string) ([]string, error) {
	return s.listYamlFilesPaths(rootPath)
}

// ReadFile - real implementation of ReadFile backed by ioutil
func (s IoImpl) ReadFile(filename string) ([]byte, error) {
	return s.readFile(filename)
}

// WriteFile - real implementation of WriteFile backed by ioutil
func (s IoImpl) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return s.writeFile(filename, data, perm)
}

// Scan - Collect images referenced in all yaml files in the path and create a kustomization file
func (s *Action) Scan(c *cli.Context) error {

	scanPath := c.String("files")
	registry := c.String("registry")
	output := c.String("output")

	dockerImages, dockerErr := s.retrieveDockerImages(scanPath, registry)
	if dockerErr != nil {
		log.Error(dockerErr)
		return dockerErr
	}

	kustomization := s.buildKustomization(dockerImages)

	writeErr := s.writeKustomizationFile(kustomization, output)
	if writeErr != nil {
		log.Error(writeErr)
		return writeErr
	}

	return nil
}

// retrieveDockerImages - Retrieve all Docker images in all yaml files in the given path
func (s *Action) retrieveDockerImages(scanPath, registry string) ([]kustomize.Image, error) {

	var dockerImages []kustomize.Image
	filesPaths, yamlErr := s.Io.ListYamlFilesPaths(scanPath)
	if yamlErr != nil {
		// log.Error(yamlErr)
		return nil, yamlErr
	}
	for _, yamlPath := range filesPaths {
		dockerImageRefs, dockerErr := s.listDockerImageReferences(yamlPath)
		if dockerErr != nil {
			// log.Error(dockerErr)
			return nil, dockerErr
		}
		if len(dockerImageRefs) > 0 {
			for _, dockerImageRef := range dockerImageRefs {
				dockerImageName, dockerImageTag := docker_utils.ParseDockerImageNameAndTag(dockerImageRef)
				dockerImages = append(dockerImages, s.buildNewDockerImage(dockerImageName, dockerImageTag, registry))
			}
		}
	}
	dockerImages = s.removeDockerImageDuplicates(dockerImages)
	log.Infof("Total Docker images found in %s: %d", scanPath, len(dockerImages))
	return dockerImages, nil
}

// listDockerImageReferences - List all Docker image reference in the given file path
func (s *Action) listDockerImageReferences(filePath string) ([]string, error) {

	fileContent, err := s.Io.ReadFile(filePath)
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
func (s *Action) buildNewDockerImage(dockerImageName, dockerImageTag, registry string) kustomize.Image {

	trimmedDockerImageName := string_utils.TrimQuotes(dockerImageName)
	normalizedImageName, _ := docker_utils.NormalizeImageName(trimmedDockerImageName)
	image := kustomize.Image{
		Name:    trimmedDockerImageName,
		NewName: registry + "/" + normalizedImageName,
	}
	if dockerImageTag != "" {
		image.NewTag = string_utils.TrimQuotes(dockerImageTag)
	}
	return image
}

// removeDockerImageDuplicates - Remove duplicated Docker image references
func (s *Action) removeDockerImageDuplicates(dockerImages []kustomize.Image) []kustomize.Image {

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
func (s *Action) buildKustomization(dockerImages []kustomize.Image) *types.Kustomization {

	return &types.Kustomization{
		TypeMeta: types.TypeMeta{
			Kind:       types.KustomizationKind,
			APIVersion: types.KustomizationVersion,
		},
		Images: dockerImages,
	}
}

// writeKustomizationFile - Write Kustomize kustomization yaml object to yaml file
func (s *Action) writeKustomizationFile(kustomization *types.Kustomization, output string) error {

	outputFileName, fileErr := filepath.Abs(output + "/" + korp_utils.KustomizationFileName)
	if fileErr != nil {
		// log.Error(fileErr)
		return fileErr
	}

	yamlContent, yamlErr := yaml.Marshal(kustomization)
	if yamlErr != nil {
		// log.Error(yamlErr)
		return yamlErr
	}

	writeErr := s.Io.WriteFile(outputFileName, yamlContent, 0644)
	if writeErr != nil {
		// log.Error(writeErr)
		return writeErr
	}

	return nil
}
