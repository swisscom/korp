package actions

import (
	"fmt"

	"github.com/urfave/cli"
)

func Patch(path *string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		paths, _ := listYamlFiles(*path)
		for _, yamlPath := range paths {
			imageRefs, _ := listImageReferences(yamlPath)
			if len(imageRefs) > 0 {
				fmt.Println(listImageReferences(yamlPath))
			}
		}
		return nil
	}
}
