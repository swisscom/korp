package actions

import (
	"context"

	"github.com/swisscom/korp/docker_utils"

	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func checkDockerDaemon(cli *client.Client, ctx *context.Context) error {

	if !docker_utils.CheckDockerDaemonRunning(cli, ctx) {
		log.Warn("Docker daemon NOT RUNNING")
		daemonErr := docker_utils.StartDockerDaemon(cli, ctx)
		if daemonErr != nil {
			// log.Error(daemonErr)
			return daemonErr
		}
	}
	return nil
}

func setLogLevel(debug *bool) {
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}
