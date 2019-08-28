package app

import (
	"os"
	"sort"
	"strings"

	// "github.com/swisscom/korp/all"
	"github.com/swisscom/korp/autocompletion"
	"github.com/swisscom/korp/korp_utils"

	// "github.com/swisscom/korp/patch"
	"github.com/swisscom/korp/pull"
	"github.com/swisscom/korp/push"
	"github.com/swisscom/korp/scan"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	version = "0.3.2"
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
			Name:     "string-config, c",
			Usage:    "Load configuration from `STRING`",
			EnvVar:   "KORP_GLOBAL_STRING_CONFIG",
			FilePath: "./config",
		},
		cli.StringFlag{
			Name:     "config, f",
			Usage:    "Load configuration from `FILE`",
			EnvVar:   "KORP_GLOBAL_CONFIG",
			FilePath: "./config",
		},
	}
}

// addBefore - Add an action to be execute before the CLI application even starts
func addBefore(app *cli.App) {

	app.Before = func(c *cli.Context) error {

		setDebug(c.Bool("debug"))
		loadConfigAsString(c.String("string-config"))

		return nil
	}
}

// setDebug -
func setDebug(debug bool) {

	if debug {
		korp_utils.SetLogLevel("debug")
	}
}

// loadConfigAsString -
func loadConfigAsString(wholeConfig string) {

	if wholeConfig != "" {
		log.Debugf("Whole config as string: %s\n", wholeConfig)

		wholeConfigSplit := strings.Fields(wholeConfig)
		log.Debugf("Whole config as array: %s\n", wholeConfigSplit)

		for _, config := range wholeConfigSplit {
			configSplit := strings.Split(config, "=")
			log.Debugf("single config: %s\n", configSplit)
			if len(configSplit) == 2 && configSplit[0] != "" && configSplit[1] != "" {
				log.Debugf("env-var set: %s=%s", configSplit[0], configSplit[1])
				os.Setenv(configSplit[0], configSplit[1])
			} else {
				log.Warnf("Config not valid: %s - skipping and taking default...", configSplit)
			}
		}
	} else {
		log.Debug("Config file not set or empty")
	}
}

// addCommands - Add commands to CLI application
func addCommands(app *cli.App) {

	app.Commands = []cli.Command{
		*scan.BuildCommand(),
		*pull.BuildCommand(),
		*push.BuildCommand(),
		// *patch.BuildCommand(),
		// *all.BuildCommand(),
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
