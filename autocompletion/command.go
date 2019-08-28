package autocompletion

import (
	"github.com/urfave/cli"
)

func BuildCommand() *cli.Command {

	return &cli.Command{
		Name:    "autocompletion",
		Aliases: []string{"a"},
		Usage:   "Generate shell autocompletion script",
		Subcommands: cli.Commands{
			{
				Name:    "bash",
				Aliases: []string{"b"},
				Usage:   "Generate bash autocompletion script",
				Action:  bashAutocomplete,
			},
			{
				Name:    "zsh",
				Aliases: []string{"z"},
				Usage:   "Generate zsh autocompletion script",
				Action:  zshAutocomplete,
			},
		},
	}
}
