package command

import (
	"flag"
	"os"

	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
)

type HelpCommand struct {
	baseCommand
}

func NewHelpCommand() *HelpCommand {

	cmd := &HelpCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	// NOTE: Help command should not include any additional flags

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *HelpCommand) Run(_ *manager.Manager) error {

	utility.Fprintln(
		os.Stdout,
		cmd.Usage(),
	)

	return nil
}
