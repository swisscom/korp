package all

import (
	"github.com/urfave/cli"
)

// BuildCommand - Build CLI application command
func BuildCommand() *cli.Command {

	return &cli.Command{
		Name:    "all",
		Aliases: []string{"a"},
		Usage:   "scan >> pull >> push",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "files, f",
				Usage:    "path to yaml files to scan (default: current dir)",
				EnvVar:   "KORP_ALL_FILES",
				Value:    ".",
				Required: false,
			},
			cli.StringFlag{
				Name:     "registry, r",
				Usage:    "name of the Docker registry to use (default: 'docker.io')",
				EnvVar:   "KORP_ALL_REGISTRY",
				Value:    "docker.io",
				Required: false,
			},
		},
		Action: all,
	}
}
