package pull_test

import (
	"context"
	"errors"
	"flag"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/matryer/is"
	"github.com/swisscom/korp/docker_utils"
	dockermocks "github.com/swisscom/korp/docker_utils/mocks"
	"github.com/swisscom/korp/pull"
	pullmocks "github.com/swisscom/korp/pull/mocks"
	"github.com/urfave/cli"
	kustimage "sigs.k8s.io/kustomize/pkg/image"
)

func TestPullCommand(t *testing.T) {
	command := pull.BuildCommand()

	t.Run("is called pull", func(t *testing.T) {
		is := is.New(t)
		is.Equal("pull", command.Name)
	})

	t.Run("has one flag", func(t *testing.T) {
		is := is.New(t)
		is.Equal(1, len(command.Flags)) // must contain 1 flag
	})

	t.Run("has a string flag called kustomization-path", func(t *testing.T) {
		is := is.New(t)
		stringFlag, ok := command.Flags[0].(cli.StringFlag)
		is.True(ok)
		is.Equal("kustomization-path, k", stringFlag.Name)

	})
}

type ReadCloserMockImpl struct {
	readFunc  func(p []byte) (n int, err error)
	closeFunc func() error
}

func (r ReadCloserMockImpl) Read(p []byte) (n int, err error) {
	return r.readFunc(p)
}

func (r ReadCloserMockImpl) Close() error {
	return r.closeFunc()
}

func getReadCloserMock() io.ReadCloser {

	readFunc := func(p []byte) (n int, err error) {
		return 0, errors.New("EOF")
	}

	closeFunc := func() error {
		return nil
	}

	return ReadCloserMockImpl{
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
		return getReadCloserMock(), nil
	}

	imagePushFunc := func(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error) {
		return getReadCloserMock(), nil
	}

	imageTagFunc := func(ctx context.Context, source string, target string) error {
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

func getIoMocks(is is.I) pullmocks.IoMock {
	loadKustomizationFileFunc := func(kstPath string) ([]kustimage.Image, error) {
		return []kustimage.Image{}, nil
	}

	openDockerClientFunc := func() (docker_utils.DockerClient, error) {
		return getDockerClientMock(is), nil
	}

	return pullmocks.IoMock{
		LoadKustomizationFileFunc: loadKustomizationFileFunc,
		OpenDockerClientFunc:      openDockerClientFunc,
	}
}

func TestPullAction(t *testing.T) {

	t.Run("reads a kustomization yaml file and pulls all images", func(t *testing.T) {
		is := is.New(t)
		ioMocks := getIoMocks(*is)
		action := pull.Action{
			Io: &ioMocks,
		}
		set := flag.NewFlagSet("test", 0)
		set.String("files", "/test/path/files", "usage files")
		set.String("registry", "registry.example.com", "usage registry")
		set.String("output", "/test/path/files", "usage output")
		context := cli.NewContext(nil, set, nil)
		action.Pull(context)
	})
}
