package command

import (
	"flag"
	"fmt"

	"github.com/StackAdapt/systags/manager"
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

func (cmd *VersionCommand) Apply(m *manager.Manager) error {

	result := fmt.Sprintf(
		"%s/%s",
		manager.AppName(),
		manager.Version(),
	)

	m.GetLogger().Info(result)

	return nil
}
