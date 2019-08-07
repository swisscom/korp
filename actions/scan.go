package actions

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli"
)

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
	var imageRefRegex = regexp.MustCompile(`(?m)image:\s*(?P<image>[^{\s]+)\s+`)
	for _, match := range imageRefRegex.FindAllStringSubmatch(string(data), -1) {
		if len(match) > 1 {
			imageRefs = append(imageRefs, match[1])
		}
	}
	return imageRefs, err
}

func Scan(path *string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		paths, _ := listYamlFiles(*path)
		for _, yamlPath := range paths {
			imageRefs, _ := listImageReferences(yamlPath)
			if len(imageRefs) > 0 {
				fmt.Println(listImageReferences(yamlPath))
			}
		}
		return nil
	}
}
