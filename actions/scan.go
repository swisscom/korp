package actions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	dockerimage "sigs.k8s.io/kustomize/pkg/image"
	"sigs.k8s.io/kustomize/pkg/types"
)

const imageNameTagRegex = "(?m)([^:]+):?(.*)"

func Scan(path, registry *string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		kustomization := types.Kustomization{}
		typeMeta := types.TypeMeta{
			Kind:       types.KustomizationKind,
			APIVersion: types.KustomizationVersion,
		}
		kustomization.TypeMeta = typeMeta
		var images []dockerimage.Image
		paths, _ := listYamlFiles(*path)
		for _, yamlPath := range paths {
			imageRefs, _ := listImageReferences(yamlPath)
			if len(imageRefs) > 0 {
				for _, imageRef := range imageRefs {
					image := dockerimage.Image{
						Name:    trimQuotes(getImageName(imageRef)),
						NewName: *registry + "/" + trimQuotes(getImageName(imageRef)),
					}
					imageTag := trimQuotes(getImageTag(imageRef))
					if imageTag != "" {
						image.NewTag = imageTag
					}
					images = append(images, image)
				}
			}
		}
		images = removeDuplicateImages(images)
		kustomization.Images = images
		outFileName, _ := filepath.Abs("./kustomization.yaml")
		yamlOutFile, _ := yaml.Marshal(kustomization)
		ioutil.WriteFile(outFileName, yamlOutFile, 0644)
		return nil
	}
}

func trimQuotes(str string) string {
	return strings.Trim(str, "'\"")
}

func removeDuplicateImages(images []dockerimage.Image) []dockerimage.Image {
	encountered := map[dockerimage.Image]bool{}
	result := []dockerimage.Image{}

	for v := range images {
		if encountered[images[v]] == true {
			// Duplicate. Do nothing.
		} else {
			// Record as duplicate and append
			encountered[images[v]] = true
			result = append(result, images[v])
		}
	}
	return result
}

func getImageName(imageRef string) string {
	var regex = regexp.MustCompile(imageNameTagRegex)
	match := regex.FindStringSubmatch(imageRef)
	return match[1]
}

func getImageTag(imageRef string) string {
	var regex = regexp.MustCompile(imageNameTagRegex)
	match := regex.FindStringSubmatch(imageRef)
	if len(match) > 2 {
		return match[2]
	}
	return ""

}

func listYamlFiles(rootpath string) ([]string, error) {
	var files []string
	var yamlRegex = regexp.MustCompile(`(?mi).*\.(yaml|yml)`)
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && yamlRegex.MatchString(info.Name()) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func listImageReferences(filepath string) ([]string, error) {
	data, err := ioutil.ReadFile(filepath)
	var imageRefs []string
	var imageRefRegex = regexp.MustCompile(`(?m)image:\s*(?P<image>[^[{\s]+)\s+`)
	for _, match := range imageRefRegex.FindAllStringSubmatch(string(data), -1) {
		if len(match) > 1 {
			imageRefs = append(imageRefs, match[1])
		}
	}
	return imageRefs, err
}
