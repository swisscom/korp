package file_utils

import (
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

const (
	yamlFileRegexStr = `(?mi).*\.(yaml|yml)`
)

// ListYamlFilesPaths - List all yaml files paths in the given root path
func ListYamlFilesPaths(rootPath string) ([]string, error) {

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
