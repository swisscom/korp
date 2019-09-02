package app

import (
	"os"
	"sort"

	"github.com/swisscom/korp/all"
	"github.com/swisscom/korp/autocompletion"
	"github.com/swisscom/korp/pull"
	"github.com/swisscom/korp/push"
	"github.com/swisscom/korp/scan"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	version = "0.4.0"
)

type CliApp struct {
	app *cli.App
}

func Create() *CliApp {

	app := cli.NewApp()
	addGlobaConfig(app)
	addGlobalFlags(app)
	addBefore(app)
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
	// app.UseShortOptionHandling = true // flag not found in this version?!
	app.EnableBashCompletion = true
}

// addGlobalFlags - Add global flag to CLI application
func addGlobalFlags(app *cli.App) {

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "switch on debug log output",
			EnvVar: "KORP_GLOBAL_DEBUG",
		},
		cli.StringFlag{
			Name:     "config, c",
			Usage:    "Load configuration from `FILE`",
			EnvVar:   "KORP_GLOBAL_CONFIG",
			FilePath: "./config",
		},
	}
}

// TODO to be completed
// addBefore - Add before-action to CLI application
func addBefore(app *cli.App) {

	app.Before = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		// TODO add config-file action
		// try using https://github.com/kelseyhightower/envconfig

		return nil
	}
}

// addCommands - Add commands to CLI application
func addCommands(app *cli.App) {

	app.Commands = []cli.Command{
		*scan.BuildCommand(),
		*pull.BuildCommand(),
		*push.BuildCommand(),
		*all.BuildCommand(),
		*autocompletion.BuildCommand(),
	}
}

// lastConfig - Add last configuration right before start CLI application
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
