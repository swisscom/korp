package scan_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/swisscom/korp/scan"
	"github.com/urfave/cli"
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
