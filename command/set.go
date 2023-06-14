package command

import (
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type SetCommand struct {
	baseCommand
	key string
	val string
}

func NewSetCommand() *SetCommand {

	cmd := &SetCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.StringVar(&cmd.key, "k", "", "")
	cmd.flagSet.StringVar(&cmd.key, "key", "", "")
	cmd.flagSet.StringVar(&cmd.val, "v", "", "")
	cmd.flagSet.StringVar(&cmd.val, "value", "", "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *SetCommand) Init(args []string) error {

	err := cmd.flagSet.Parse(args)
	if err != nil {
		return err
	}

	if cmd.key == "" {
		return cmd.failf("flag needs to be provided: -key")
	}

	if cmd.val == "" {
		return cmd.failf("flag needs to be provided: -value")
	}

	return nil
}

func (cmd *SetCommand) Run(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	m.SetTag(cmd.key, cmd.val)

	err = m.SaveFiles()
	if err != nil {
		return err
	}

	return nil
}
