package scan_test

import (
	"flag"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/swisscom/korp/scan"
	"github.com/swisscom/korp/scan/mocks"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	kust "sigs.k8s.io/kustomize/pkg/image"
	"sigs.k8s.io/kustomize/pkg/types"
)

func TestBuildCommand(t *testing.T) {

	command := scan.BuildCommand()

	t.Run("is called scan", func(t *testing.T) {
		is := is.New(t)
		is.Equal("scan", command.Name)
	})

	t.Run("has three flags", func(t *testing.T) {
		is := is.New(t)
		is.Equal(3, len(command.Flags)) // must contain 3 flags
	})

	t.Run("has a string flag called files", func(t *testing.T) {
		is := is.New(t)
		stringFlag, ok := command.Flags[0].(cli.StringFlag)
		is.True(ok)
		is.Equal("files, f", stringFlag.Name)

	})

	t.Run("has a flag called registry", func(t *testing.T) {
		is := is.New(t)
		stringFlag, ok := command.Flags[1].(cli.StringFlag)
		is.True(ok)
		is.Equal("registry, r", stringFlag.Name)
	})

	t.Run("has a flag called output", func(t *testing.T) {
		is := is.New(t)
		stringFlag, ok := command.Flags[2].(cli.StringFlag)
		is.True(ok)
		is.Equal("output, o", stringFlag.Name)
	})
}

func containsImage(images []kust.Image, img kust.Image) bool {
	for _, i := range images {
		if i == img {
			return true
		}
	}
	return false
}

func TestScanAction(t *testing.T) {

	t.Run("reads yaml files with image references and produces a kustomization.yml file", func(t *testing.T) {
		is := is.New(t)

		listYamlFilesPathsMock := func(rootPath string) ([]string, error) {
			return []string{"/test/path/test.yaml"}, nil
		}

		readFileMock := func(filename string) ([]byte, error) {
			return []byte(mocks.YamlTwoImages), nil
		}

		writeFileMock := func(filename string, data []byte, perm os.FileMode) error {
			var kustomization types.Kustomization
			yamlErr := yaml.Unmarshal(data, &kustomization)
			is.Equal(nil, yamlErr)                 // must be valid yaml
			is.Equal(2, len(kustomization.Images)) // must produce two image elements
			minideb := kust.Image{Name: "bitnami/minideb", NewName: "registry.example.com/bitnami/minideb", NewTag: "latest"}
			is.True(containsImage(kustomization.Images, minideb)) // must contain patch for bitnami/minideb
			postgres := kust.Image{Name: "bitnami/postgresql", NewName: "registry.example.com/bitnami/postgresql", NewTag: "11.3.0-debian-9-r38"}
			is.True(containsImage(kustomization.Images, postgres)) // must contain patch for bitnami/postgresql
			return nil
		}

		ioMock := mocks.IoMock{
			ListYamlFilesPathsFunc: listYamlFilesPathsMock,
			ReadFileFunc:           readFileMock,
			WriteFileFunc:          writeFileMock,
		}

		action := scan.Action{Io: &ioMock}
		set := flag.NewFlagSet("test", 0)
		set.String("files", "/test/path/files", "usage files")
		set.String("registry", "registry.example.com", "usage registry")
		set.String("output", "/test/path/files", "usage output")
		context := cli.NewContext(nil, set, nil)
		action.Scan(context)
	})

}
