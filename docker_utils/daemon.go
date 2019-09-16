package docker_utils

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/swisscom/korp/cli_utils"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// CheckDockerDaemon - Check wheather the Docker daemon is running, if the method starts it
func CheckDockerDaemon(cli DockerClient, ctx *context.Context) error {

	if !CheckDockerDaemonRunning(cli, ctx) {
		log.Warn("Docker daemon NOT RUNNING")
		daemonErr := StartDockerDaemon(cli, ctx)
		if daemonErr != nil {
			// log.Error(daemonErr)
			return daemonErr
		}
	}
	return nil
}

// CheckDockerDaemonRunning - Check wheather the Docker daemon is running (making a simple request)
func CheckDockerDaemonRunning(cli DockerClient, ctx *context.Context) bool {

	_, listErr := cli.ContainerList(*ctx, types.ContainerListOptions{})
	return listErr == nil
}

// StartDockerDaemon - Start up the Docker deamon based on the current OS
func StartDockerDaemon(cli DockerClient, ctx *context.Context) error {

	log.Info("Starting Docker daemon...")

	switch runtime.GOOS {
	case "darwin": // macos
		// open --background -a Docker
		_, startErr := cli_utils.ExecCommand("open", "--background", "-a", "Docker")
		if startErr != nil {
			// log.Error(startErr)
			return startErr
		}
	case "linux":
		// TODO: which is the right command?
		return errors.New("Start Docker daemon on Linux not yet supported, please run it manually and retry!")
	case "windows":
		// TODO: which is the right command?
		return errors.New("Start Docker daemon on Windows not yet supported, please run it manually and retry!")
	}

	log.Info("Waiting for Docker daemon to be up&running...")
	waitDockerDaemon(cli, ctx)

	return nil
}

// waitDockerDaemon - Wait for Docker daemon to be up&running
func waitDockerDaemon(cli DockerClient, ctx *context.Context) {

	for {
		_, err := cli.ContainerList(*ctx, types.ContainerListOptions{})
		if err == nil {
			break
		}
		log.Info("Waiting for Docker daemon to be up&running for other 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	log.Info("Docker daemon be up&running")
}
