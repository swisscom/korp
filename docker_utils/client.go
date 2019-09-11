package docker_utils

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerClient interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImagePush(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error)
	ImageTag(ctx context.Context, source, target string) error
	Close() error
}

// OpenDockerClient - Create new Docker client
func OpenDockerClient() (DockerClient, error) {

	cli, cliErr := client.NewEnvClient()
	if cliErr != nil {
		// log.Error(cliErr)
		return nil, cliErr
	}
	return cli, nil
}
