package all

import (
	"github.com/urfave/cli"
)

func BuildCommand() *cli.Command {

	var filesPath, registry, patch string
	return &cli.Command{
		Name:    "all",
		Aliases: []string{"a"},
		Usage:   "scan >> pull >> push [>> patch]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "files, f",
				Usage:       "path to yaml files to scan (default: current dir)",
				EnvVar:      "KORP_ALL_FILES",
				Value:       ".",
				Required:    false,
				Destination: &filesPath,
			},
			cli.StringFlag{
				Name:        "registry, r",
				Usage:       "name of the Docker registry to use (default: 'docker.io')",
				EnvVar:      "KORP_ALL_REGISTRY",
				Value:       "docker.io",
				Required:    false,
				Destination: &registry,
			},
			cli.StringFlag{
				Name:        "patch, p",
				Usage:       "execute patch phase",
				EnvVar:      "KORP_ALL_PATCH",
				Value:       "false",
				Required:    false,
				Destination: &patch,
			},
		},
		Action: all(&filesPath, &registry, &patch),
	}
}
