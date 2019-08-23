package utils

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	dockerImageNameTagRegexStr = "(?m)([^:]+):?(.*)"
	defaultReturn              = ""
)

// TrimQuotes - Remove single or double quotes from a string
func TrimQuotes(str string) string {

	return strings.Trim(str, "'\"")
}

// GetDockerImageNameAndTag - Retrieve Docker image name and tag from a string
func GetDockerImageNameAndTag(imageRef string) (dockerImageName, dockerImageTag string) {

	var dockerImageNameTagRegex = regexp.MustCompile(dockerImageNameTagRegexStr)
	regexMatch := dockerImageNameTagRegex.FindStringSubmatch(imageRef)
	dockerImageName = defaultReturn
	if len(regexMatch) > 1 {
		dockerImageName = regexMatch[1]
	} else {
		log.Warn("No Docker image name found")
	}
	dockerImageTag = defaultReturn
	if len(regexMatch) > 2 {
		dockerImageTag = regexMatch[2]
	} else {
		log.Debug("No Docker image tag found")
	}
	return dockerImageName, dockerImageTag
}

// GetDockerImageName - Retrieve Docker image name from a string
func GetDockerImageName(imageRef string) string {

	var dockerImageNameTagRegex = regexp.MustCompile(dockerImageNameTagRegexStr)
	match := dockerImageNameTagRegex.FindStringSubmatch(imageRef)
	if len(match) > 1 {
		return match[1]
	}
	log.Warn("No Docker image name found")
	return defaultReturn
}

// GetDockerImageTag - Retrieve Docker image tag from a string
func GetDockerImageTag(imageRef string) string {

	var dockerImageNameTagRegex = regexp.MustCompile(dockerImageNameTagRegexStr)
	match := dockerImageNameTagRegex.FindStringSubmatch(imageRef)
	if len(match) > 2 {
		return match[2]
	}
	log.Warn("No Docker image tag found")
	return defaultReturn
}
