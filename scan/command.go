package scan

import (
	"github.com/urfave/cli"
)

func BuildCommand() *cli.Command {

	var filesPath, registry, output string
	return &cli.Command{
		Name:    "scan",
		Aliases: []string{"s"},
		Usage:   "collect images referenced in all yaml files in the path and create a kustomization file",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "files, f",
				Usage:       "path to yaml files to scan (default: current dir)",
				EnvVar:      "KORP_SCAN_FILES",
				Value:       ".",
				Required:    false,
				Destination: &filesPath,
			},
			cli.StringFlag{
				Name:        "registry, r",
				Usage:       "name of the Docker registry to use (default: 'docker.io')",
				EnvVar:      "KORP_SCAN_REGISTRY",
				Value:       "docker.io",
				Required:    false,
				Destination: &registry,
			},
			cli.StringFlag{
				Name:        "output, o",
				Usage:       "path of the kustomization file to be written (default: current dir)",
				EnvVar:      "KORP_SCAN_OUTPUT",
				Value:       ".",
				Required:    false,
				Destination: &output,
			},
		},
		Action: scan(&filesPath, &registry, &output),
	}
}
