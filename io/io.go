package io

import (
	"github.com/swisscom/korp/docker_utils"
	kust "sigs.k8s.io/kustomize/pkg/image"
)

//go:generate moq -out mocks/io.go -pkg mocks . PullPushIo

// Io - interface for all io functions used by pull
type PullPushIo interface {
	LoadKustomizationFile(kstPath string) ([]kust.Image, error)
	OpenDockerClient() (docker_utils.DockerClient, error)
}

// IoImpl - real io implementation using kustomize_utils and docker_utils
type PullPushIoImpl struct {
	LoadKustomizationFileFunc func(kstPath string) ([]kust.Image, error)
	OpenDockerClientFunc      func() (docker_utils.DockerClient, error)
}

func (i PullPushIoImpl) LoadKustomizationFile(kstPath string) ([]kust.Image, error) {
	return i.LoadKustomizationFileFunc(kstPath)
}

func (i PullPushIoImpl) OpenDockerClient() (docker_utils.DockerClient, error) {
	return i.OpenDockerClientFunc()
}
