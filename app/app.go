package app

import (
	"os"
	"sort"

	"github.com/swisscom/korp/all"
	"github.com/swisscom/korp/patch"
	"github.com/swisscom/korp/pull"
	"github.com/swisscom/korp/push"
	"github.com/swisscom/korp/scan"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	version = "0.0.4"
)

type CliApp struct {
	app *cli.App
}

func Create() *CliApp {

	app := cli.NewApp()
	addGlobaConfig(app)
	// addGlobalFlags(app)
	// addBefore(app)
	addCommands(app)
	lastConfig(app)
	return &CliApp{
		app: app,
	}
}

// globaConfig - Add global configurations to CLI application
func addGlobaConfig(app *cli.App) {

	app.Name = "korp"
	app.Usage = "push images to a corporate registry based on Kubernetes yaml files"
	app.Version = version
	// app.UseShortOptionHandling = true // ?!
}

// TODO to be completed
// addGlobalFlags - Add global flag to CLI application
func addGlobalFlags(app *cli.App) {

	app.Flags = []cli.Flag{

		// TODO add debug flag

		cli.StringFlag{
			Name:     "config, c",
			Usage:    "Load configuration from `FILE`",
			FilePath: "~/.korp/config",
		},
	}
}

// TODO to be implemented
// addBefore - Add before-action to CLI application
func addBefore(app *cli.App) {

	// TODO add debug flag action

	// TODO add config-file action
}

// addCommands - Add commands to CLI application
func addCommands(app *cli.App) {

	app.Commands = []cli.Command{
		*scan.BuildCommand(),
		*pull.BuildCommand(),
		*push.BuildCommand(),
		*patch.BuildCommand(),
		*all.BuildCommand(),
	}
}

func lastConfig(app *cli.App) {

	// sorting flags in help section
	sort.Sort(cli.FlagsByName(app.Flags))
	// sorting commands in help section
	sort.Sort(cli.CommandsByName(app.Commands))
}

// Start - Execute CLI application
func (a *CliApp) Start() {

	err := a.app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
