package pull

import (
	"github.com/swisscom/korp/docker_utils"
	korpio "github.com/swisscom/korp/io"
	"github.com/swisscom/korp/kustomize_utils"
	"github.com/urfave/cli"
)

// BuildCommand - Build CLI application command
func BuildCommand() *cli.Command {

	io := korpio.PullPushIoImpl{
		OpenDockerClientFunc:      docker_utils.OpenDockerClient,
		LoadKustomizationFileFunc: kustomize_utils.LoadKustomizationFile,
	}

	action := Action{
		Io: io,
	}

	return &cli.Command{
		Name:    "pull",
		Aliases: []string{"p"},
		Usage:   "pull Docker images listed in the kustomization file to the local Docker registry",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "kustomization-path, k",
				Usage:    "path to the kustomization file (default: current dir)",
				EnvVar:   "KORP_PULL_KST_PATH",
				Value:    ".",
				Required: false,
			},
			cli.StringFlag{
				Name:     "username, u",
				Usage:    "User name for accessing a registry which requires authentication",
				EnvVar:   "KORP_PULL_USERNAME",
				Value:    "",
				Required: false,
			},
			cli.StringFlag{
				Name:     "password, p",
				Usage:    "Password for accessing a registry which requires authentication",
				EnvVar:   "KORP_PULL_PASSWORD",
				Value:    "",
				Required: false,
			},
		},
		Action: action.Pull,
	}
}
