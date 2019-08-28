package push

import (
	"github.com/urfave/cli"
)

// BuildCommand - Build CLI application command
func BuildCommand() *cli.Command {

	return &cli.Command{
		Name:    "push",
		Aliases: []string{"u"},
		Usage:   "re-tag original Docker images and push them to the new Docker registry",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "kustomization-path, k",
				Usage:    "path to the kustomization file (default: current dir)",
				EnvVar:   "KORP_PUSH_KST_PATH",
				Value:    ".",
				Required: false,
			},
		},
		Action: push,
	}
}
