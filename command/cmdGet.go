package command

import (
	"errors"
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type GetCommand struct {
	baseCommand
	key string
	def string
}

func NewGetCommand() *GetCommand {

	cmd := &GetCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.StringVar(&cmd.key, "k", "", "")
	cmd.flagSet.StringVar(&cmd.key, "key", "", "")
	cmd.flagSet.StringVar(&cmd.def, "d", "", "")
	cmd.flagSet.StringVar(&cmd.def, "default", "", "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *GetCommand) Parse(args []string) error {

	err := cmd.flagSet.Parse(args)
	if err != nil {
		return err
	}

	if cmd.key == "" {
		return errors.New("flag needs to be provided: -key")
	}

	return nil
}

func (cmd *GetCommand) Apply(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	m.GetLogger().Info(m.GetTag(cmd.key, cmd.def))

	return nil
}
