package docker_utils

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

const (
	dockerImageNameTagRegexStr = "(?m)([^:]+):?(.*)"
	defaultReturn              = ""
)

// ParseDockerImageNameAndTag - Retrieve Docker image name and tag from a string
func ParseDockerImageNameAndTag(imageRef string) (dockerImageName, dockerImageTag string) {

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

// ParseDockerImageName - Retrieve Docker image name from a string
func ParseDockerImageName(imageRef string) string {

	var dockerImageNameTagRegex = regexp.MustCompile(dockerImageNameTagRegexStr)
	match := dockerImageNameTagRegex.FindStringSubmatch(imageRef)
	if len(match) > 1 {
		return match[1]
	}
	log.Warn("No Docker image name found")
	return defaultReturn
}

// ParseDockerImageTag - Retrieve Docker image tag from a string
func ParseDockerImageTag(imageRef string) string {

	var dockerImageNameTagRegex = regexp.MustCompile(dockerImageNameTagRegexStr)
	match := dockerImageNameTagRegex.FindStringSubmatch(imageRef)
	if len(match) > 2 {
		return match[2]
	}
	log.Warn("No Docker image tag found")
	return defaultReturn
}

// BuildCompleteDockerImage - Build the complete Docker image name from name and tag
func BuildCompleteDockerImage(imageName, imageTag string) string {

	return imageName + ":" + imageTag
}
