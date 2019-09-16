package mocks

import (
	"context"
	"errors"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/matryer/is"
	"github.com/swisscom/korp/docker_utils"
	dockermocks "github.com/swisscom/korp/docker_utils/mocks"
	kustimage "sigs.k8s.io/kustomize/pkg/image"
)

type readCloserMockImpl struct {
	readFunc  func(p []byte) (n int, err error)
	closeFunc func() error
}

func (r readCloserMockImpl) Read(p []byte) (n int, err error) {
	return r.readFunc(p)
}

func (r readCloserMockImpl) Close() error {
	return r.closeFunc()
}

func getReadCloserMock() io.ReadCloser {

	readFunc := func(p []byte) (n int, err error) {
		return 0, errors.New("EOF")
	}

	closeFunc := func() error {
		return nil
	}

	return readCloserMockImpl{
		readFunc, closeFunc,
	}
}

func getDockerClientMock(is is.I) docker_utils.DockerClient {

	closeFunc := func() error {
		return nil
	}

	containerListFunc := func(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
		dummyContainer := types.Container{
			ID:      "dummyid",
			Names:   []string{"dummyname1", "dummyname2"},
			Image:   "dummyimage",
			ImageID: "dummyimageid",
		}
		return []types.Container{dummyContainer}, nil
	}

	imagePullFunc := func(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
		is.Equal("docker.io/bitnami/minideb:latest", refStr) // uses correct image reference in pull call
		return getReadCloserMock(), nil
	}

	imagePushFunc := func(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error) {
		is.Equal("registry.example.com/bitnami/minideb:latest", image)
		return getReadCloserMock(), nil
	}

	imageTagFunc := func(ctx context.Context, source string, target string) error {
		is.Equal("docker.io/bitnami/minideb:latest", source)
		is.Equal("registry.example.com/bitnami/minideb:latest", target)
		return nil
	}

	result := dockermocks.DockerClientMock{
		CloseFunc:         closeFunc,
		ContainerListFunc: containerListFunc,
		ImagePullFunc:     imagePullFunc,
		ImagePushFunc:     imagePushFunc,
		ImageTagFunc:      imageTagFunc,
	}

	return &result
}

// GetIoMocks -- get mock implementation of PullPushIo
func GetIoMocks(is is.I) PullPushIoMock {
	loadKustomizationFileFunc := func(kstPath string) ([]kustimage.Image, error) {
		image := kustimage.Image{
			Name:    "bitnami/minideb",
			NewName: "registry.example.com/bitnami/minideb",
			NewTag:  "latest",
		}
		return []kustimage.Image{image}, nil
	}

	openDockerClientFunc := func() (docker_utils.DockerClient, error) {
		return getDockerClientMock(is), nil
	}

	return PullPushIoMock{
		LoadKustomizationFileFunc: loadKustomizationFileFunc,
		OpenDockerClientFunc:      openDockerClientFunc,
	}
}
