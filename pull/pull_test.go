package pull_test

import (
	"flag"
	"testing"

	"github.com/matryer/is"
	"github.com/swisscom/korp/pull"
	"github.com/urfave/cli"

	"github.com/swisscom/korp/io/mocks"
)

func TestPullCommand(t *testing.T) {
	command := pull.BuildCommand()

	t.Run("is called pull", func(t *testing.T) {
		is := is.New(t)
		is.Equal("pull", command.Name)
	})

	t.Run("has one flag", func(t *testing.T) {
		is := is.New(t)
		is.Equal(1, len(command.Flags)) // must contain 1 flag
	})

	t.Run("has a string flag called kustomization-path", func(t *testing.T) {
		is := is.New(t)
		stringFlag, ok := command.Flags[0].(cli.StringFlag)
		is.True(ok)
		is.Equal("kustomization-path, k", stringFlag.Name)

	})
}

func TestPullAction(t *testing.T) {

	t.Run("reads a kustomization yaml file and pulls all images", func(t *testing.T) {
		is := is.New(t)
		ioMocks := mocks.GetIoMocks(*is)
		action := pull.Action{
			Io: &ioMocks,
		}
		set := flag.NewFlagSet("test", 0)
		set.String("kustomization-path", "/test/path", "usage kustomization-path")
		context := cli.NewContext(nil, set, nil)
		action.Pull(context)
	})
}
