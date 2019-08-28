package all

import (
	"github.com/urfave/cli"
)

// BuildCommand - Build CLI application command
func BuildCommand() *cli.Command {

	return &cli.Command{
		Name:    "all",
		Aliases: []string{"a"},
		Usage:   "scan >> pull >> push [>> patch]",
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
			cli.BoolFlag{
				Name:     "patch, p",
				Usage:    "execute patch phase",
				EnvVar:   "KORP_ALL_PATCH",
				Required: false,
			},
			cli.StringFlag{
				Name:     "kustomization-path, k",
				Usage:    "path to the kustomization file (default: current dir)",
				EnvVar:   "KORP_PATCH_KST_PATH",
				Value:    ".",
				Required: false,
			},
		},
		Action: all,
	}
}
