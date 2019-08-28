package cli_utils

import (
	"os/exec"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

// ExecCommand - Execute a command on the current OS
func ExecCommand(command string, args ...string) ([]byte, error) {

	return exec.Command(command, args...).Output()
}

// GetFileFlagsBeforeFunc - get flags from a file, being as tolerant as possible
func GetFileFlagsBeforeFunc(flags []cli.Flag, flagName string) cli.BeforeFunc {
	return func(context *cli.Context) error {
		var filepath = context.String(flagName)
		filepath, err := homedir.Expand(filepath)
		if err != nil {
			log.Debugf("Home directory could not be expanded: %s", filepath)
			return nil
		}
		log.Debugf("Home directory expanded: %s", filepath)
		sourceContext, err := altsrc.NewYamlSourceFromFile(filepath)
		if err != nil {
			log.Debugf("Config file not found at %s", filepath)
			return nil
		}
		log.Debugf("Config file found at %s", filepath)
		log.Debugf("Source context %s", sourceContext)
		log.Debugf("Flags %s", flags)
		return altsrc.InitInputSourceWithContext(flags, func(context *cli.Context) (altsrc.InputSourceContext, error) { return sourceContext, nil })(context)
	}
}
