package pull

import (
	"github.com/urfave/cli"
)

func BuildCommand() *cli.Command {

	var kstPath string
	return &cli.Command{
		Name:    "pull",
		Aliases: []string{"p"},
		Usage:   "pull Docker images listed in the kustomization file to the local Docker registry",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "kustomization-path, k",
				Usage:       "path to the kustomization file (default: current dir)",
				EnvVar:      "KORP_PULL_KST_PATH",
				Value:       ".",
				Required:    false,
				Destination: &kstPath,
			},
		},
		Action: pull(&kstPath),
	}
}
