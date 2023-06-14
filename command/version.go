package command

import (
	"flag"
	"os"

	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
)

type VersionCommand struct {
	baseCommand
}

func NewVersionCommand() *VersionCommand {

	cmd := &VersionCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *VersionCommand) Run(_ *manager.Manager) error {

	utility.Fprintf(
		os.Stdout,
		"%s/%s\n",
		manager.AppName(),
		manager.Version(),
	)

	return nil
}
