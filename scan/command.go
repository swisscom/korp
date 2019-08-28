package scan

import (
	"github.com/swisscom/korp/cli_utils"
	"github.com/urfave/cli"
)

// BuildCommand - Build CLI application command
func BuildCommand() *cli.Command {

	var flags = []cli.Flag{
		cli.StringFlag{
			Name:     "files, f",
			Usage:    "path to yaml files to scan (default: current dir)",
			EnvVar:   "KORP_SCAN_FILES",
			Value:    ".",
			Required: false,
		},
		cli.StringFlag{
			Name:     "registry, r",
			Usage:    "name of the Docker registry to use (default: 'docker.io')",
			EnvVar:   "KORP_SCAN_REGISTRY",
			Required: false,
		},
		cli.StringFlag{
			Name:     "output, o",
			Usage:    "path of the kustomization file to be written (default: current dir)",
			EnvVar:   "KORP_SCAN_OUTPUT",
			Value:    ".",
			Required: false,
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: "KORP_SCAN_CONFIG",
			Value:  "~/.korp/config.yml",
		},
	}

	return &cli.Command{
		Name:    "scan",
		Aliases: []string{"s"},
		Usage:   "collect images referenced in all yaml files in the path and create a kustomization file",
		Flags:   flags,
		Action:  scan,
		Before:  cli_utils.GetFileFlagsBeforeFunc(flags, "config"),
	}
}
