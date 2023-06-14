package command

import (
	"flag"
	"os"

	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
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

func (cmd *GetCommand) Init(args []string) error {

	err := cmd.flagSet.Parse(args)
	if err != nil {
		return err
	}

	if cmd.key == "" {
		return cmd.failf("flag needs to be provided: -key")
	}

	return nil
}

func (cmd *GetCommand) Run(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	utility.Fprintln(
		os.Stdout,
		m.GetTag(cmd.key, cmd.def),
	)

	return nil
}
