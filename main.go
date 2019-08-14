package main

import (
	"log"
	"os"

	"github.com/swisscom/korp/actions"
	"github.com/urfave/cli"
)

func main() {
	var scanpath, registry string
	app := cli.NewApp()
	app.Name = "korp"
	app.Usage = "push images to a corporate registry based on Kubernetes yaml files"

	app.Commands = []cli.Command{
		{
			Name:    "scan",
			Aliases: []string{"s"},
			Usage:   "collects images referenced in yaml files and creates a kustomization file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "files, f",
					Value:       ".",
					Usage:       "path of yaml files to scan",
					Destination: &scanpath,
				},
				cli.StringFlag{
					Name:        "registry, r",
					Value:       "docker.io",
					Usage:       "name of the corporate registry to use",
					Destination: &registry,
				},
			},
			Action: actions.Scan(&scanpath, &registry),
		},
		{
			Name:    "pull",
			Aliases: []string{"p"},
			Usage:   "pulls referenced images to the local docker registry",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path, p",
					Value:       ".",
					Usage:       "path of the kustomize file",
					Destination: &scanpath,
				},
			},
			Action: actions.Pull(&scanpath),
		},
		{
			Name:    "push",
			Aliases: []string{"u"},
			Usage:   "tags images and pushes them to the corporate registry",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path, p",
					Value:       ".",
					Usage:       "path of the kustomize file",
					Destination: &scanpath,
				},
				cli.StringFlag{
					Name:        "registry, r",
					Value:       "docker.io",
					Usage:       "name of the corporate registry to use",
					Destination: &registry,
				},
			},
			Action: actions.Push(&scanpath),
		},
		{
			Name:    "patch",
			Aliases: []string{"a"},
			Usage:   "patches the image references in all yaml files in the path",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path, p",
					Value:       ".",
					Usage:       "path of the kustomize file and the yaml files",
					Destination: &scanpath,
				},
			},
			Action: actions.Patch(&scanpath),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
