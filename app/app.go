package app

import (
	"os"
	"sort"

	"github.com/swisscom/korp/actions"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	version = "0.0.1"
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
