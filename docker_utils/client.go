package docker_utils

import (
	"github.com/docker/docker/client"
)

// OpenDockerClient - Create new Docker client
func OpenDockerClient() (*client.Client, error) {

	cli, cliErr := client.NewEnvClient()
	if cliErr != nil {
		// log.Error(cliErr)
		return nil, cliErr
	}
	return cli, nil
}
