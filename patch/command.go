package patch

import (
	"github.com/urfave/cli"
)

func BuildCommand() *cli.Command {

	var filesPath, kstPath string
	return &cli.Command{
		Name:    "patch",
		Aliases: []string{"a"},
		Usage:   "patch all yaml files in the path with Docker images tags to the new Docker registry",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "files, f",
				Usage:       "path to yaml files to patch (default: current dir)",
				EnvVar:      "KORP_PATCH_FILES",
				Value:       ".",
				Required:    false,
				Destination: &filesPath,
			},
			cli.StringFlag{
				Name:        "kustomization-path, k",
				Usage:       "path to the kustomization file (default: current dir)",
				EnvVar:      "KORP_PATCH_KST_PATH",
				Value:       ".",
				Required:    false,
				Destination: &kstPath,
			},
		},
		Action: patch(&filesPath, &kstPath),
	}
}
