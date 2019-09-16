package push_test

import (
	"flag"
	"testing"

	"github.com/matryer/is"
	"github.com/urfave/cli"

	"github.com/swisscom/korp/io/mocks"
	"github.com/swisscom/korp/push"
)

func TestPushCommand(t *testing.T) {
	command := push.BuildCommand()

	t.Run("is called push", func(t *testing.T) {
		is := is.New(t)
		is.Equal("push", command.Name) // is called push
	})

	t.Run("has one flag", func(t *testing.T) {
		is := is.New(t)
		is.Equal(1, len(command.Flags)) // must contain 1 flag
	})

	t.Run("has a string flag called kustomization-path", func(t *testing.T) {
		is := is.New(t)
		stringFlag, ok := command.Flags[0].(cli.StringFlag)
		is.True(ok)
		is.Equal("kustomization-path, k", stringFlag.Name) // the flag must be called kustomization-path

	})
}

func TestPushAction(t *testing.T) {

	t.Run("reads a kustomization yaml file and pushes all images", func(t *testing.T) {
		is := is.New(t)
		ioMocks := mocks.GetIoMocks(*is)
		action := push.Action{
			Io: &ioMocks,
		}
		set := flag.NewFlagSet("test", 0)
		set.String("kustomization-path", "/test/path", "usage kustomization-path")
		context := cli.NewContext(nil, set, nil)
		action.Push(context)
	})
}
