package command

import (
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type RmCommand struct {
	baseCommand
	key string
}

func NewRmCommand() *RmCommand {

	cmd := &RmCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.StringVar(&cmd.key, "k", "", "")
	cmd.flagSet.StringVar(&cmd.key, "key", "", "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *RmCommand) Init(args []string) error {

	err := cmd.flagSet.Parse(args)
	if err != nil {
		return err
	}

	if cmd.key == "" {
		return cmd.failf("flag needs to be provided: -key")
	}

	return nil
}

func (cmd *RmCommand) Run(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	m.RemoveTag(cmd.key)

	err = m.SaveFiles()
	if err != nil {
		return err
	}

	return nil
}
