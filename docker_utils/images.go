package docker_utils

import (
	"github.com/docker/distribution/reference"
)

// NormalizeImageName - Normalize Docker image names to canonical
// WARN: Normalizing a Docker image reference in the form name:tag will output always name:latest
func NormalizeImageName(imageName string) (string, error) {

	named, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return "", err
	}

	return named.Name(), nil
}
