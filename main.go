package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/swisscom/korp/actions"
	"github.com/urfave/cli"
)

const (
	version  = "0.0.1"
	logLevel = "info"
)

// main -
func main() {

	setLogLevel(logLevel)

	app := createApp()
	addFlags(app)
	addBefore(app)
	addCommands(app)
	execApp(app)
}

// setLogLevel - Setup logrus logging level [error, warn, info, debug]
func setLogLevel(levelStr string) {

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(level)
}

// createApp - Create CLI application
func createApp() *cli.App {

	app := cli.NewApp()
	app.Name = "korp"
	app.Usage = "push images to a corporate registry based on Kubernetes yaml files"
	app.Version = version
	return app
}

// addBefore - Add a before handler to the app, called before any command
func addBefore(app *cli.App) {
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
}

// addFlags - Add global flags to CLI application
func addFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "switch on debug log output",
			EnvVar: "KORP_LOG_DEBUG",
		},
	}
}

// addCommands - Add commands to CLI application
func addCommands(app *cli.App) {

	var filesPath, registry, output, kstPath, patch string
	app.Commands = []cli.Command{
		{
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
			Action: actions.Scan(&filesPath, &registry, &output),
		},
		{
			Name:    "pull",
			Aliases: []string{"p"},
			Usage:   "pull Docker images listed in the kustomization file to the local Docker registry",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "kustomization-path, k",
					Usage:       "path to the kustomization file (default: current dir)",
					EnvVar:      "KORP_PULL_KST_PATH",
					Value:       ".",
					Required:    false,
					Destination: &kstPath,
				},
			},
			Action: actions.Pull(&kstPath),
		},
		{
			Name:    "push",
			Aliases: []string{"u"},
			Usage:   "re-tag original Docker images and push them to the new Docker registry",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "kustomization-path, k",
					Usage:       "path to the kustomization file (default: current dir)",
					EnvVar:      "KORP_PUSH_KST_PATH",
					Value:       ".",
					Required:    false,
					Destination: &kstPath,
				},
			},
			Action: actions.Push(&kstPath),
		},
		{
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
			Action: actions.Patch(&filesPath, &kstPath),
		},
		{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "scan >> pull >> push [>> patch]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "files, f",
					Usage:       "path to yaml files to scan (default: current dir)",
					EnvVar:      "KORP_ALL_FILES",
					Value:       ".",
					Required:    false,
					Destination: &filesPath,
				},
				cli.StringFlag{
					Name:        "registry, r",
					Usage:       "name of the Docker registry to use (default: 'docker.io')",
					EnvVar:      "KORP_ALL_REGISTRY",
					Value:       "docker.io",
					Required:    false,
					Destination: &registry,
				},
				cli.StringFlag{
					Name:        "patch, p",
					Usage:       "execute patch phase",
					EnvVar:      "KORP_ALL_PATCH",
					Value:       "false",
					Required:    false,
					Destination: &patch,
				},
			},
			Action: actions.All(&filesPath, &registry, &patch),
		},
	}
}

// execApp - Execute CLI application
func execApp(app *cli.App) {

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
