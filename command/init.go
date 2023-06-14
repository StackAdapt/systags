package command

import (
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type InitCommand struct {
	baseCommand
	reset bool
}

func NewInitCommand() *InitCommand {

	cmd := &InitCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.BoolVar(&cmd.reset, "r", false, "")
	cmd.flagSet.BoolVar(&cmd.reset, "reset", false, "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *InitCommand) Run(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	if cmd.reset {
		m.Reset()
	}

	err = m.SaveFiles()
	if err != nil {
		return err
	}

	return nil
}
