package command

import (
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type Command interface {
	Parse([]string) error
	Apply(*manager.Manager) error
}

type baseCommand struct {
	flagSet *flag.FlagSet
}

func (cmd *baseCommand) Parse(args []string) error {
	return cmd.flagSet.Parse(args)
}

func (cmd *baseCommand) Apply(_ *manager.Manager) error {
	return nil
}
