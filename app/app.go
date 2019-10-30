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
	"github.com/spf13/viper"
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
			Name:     "config-string, c",
			Usage:    "Load configuration from `STRING`",
			EnvVar:   "KORP_GLOBAL_CONFIG_STRING",
			FilePath: "./config",
		},
		cli.StringFlag{
			Name:   "config-path, f",
			Usage:  "`PATH` of the .env file containing configuration for korp",
			EnvVar: "KORP_GLOBAL_CONFIG_PATH",
			Value:  "./",
		},
	}
}

// addBefore - Add an action to be execute before the CLI application even starts
func addBefore(app *cli.App) {

	app.Before = func(c *cli.Context) error {

		setDebug(c.Bool("debug"))

		// loadStringConfig(c.String("config-string"))
		// _, strOk := os.LookupEnv("KORP_SCAN_FILES")
		// if strOk {
		// 	log.Info("KORP_SCAN_FILES loaded")
		// }

		// loadFileConfig(c.String("config-path"))
		// _, fileOk := os.LookupEnv("KORP_PULL_KST_PATH")
		// if fileOk {
		// 	log.Info("KORP_PULL_KST_PATH loaded")
		// }

		return nil
	}
}

// setDebug -
func setDebug(debug bool) {

	if debug {
		korp_utils.SetLogLevel("debug")
	}
}

// loadStringConfig -
func loadStringConfig(wholeConfig string) {

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

// TODO to be fixed, it does not work
// loadFileConfig -
func loadFileConfig(configPath string) {

	viper.SetConfigName(".env")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warnf(".env config file not found in path %s", configPath)
		} else {
			log.Errorf("Error loading .env config file from path %s: %s", configPath, err)
		}
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
