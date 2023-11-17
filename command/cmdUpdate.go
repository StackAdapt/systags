package command

import (
	"flag"
	"time"

	"github.com/StackAdapt/systags/manager"
)

type UpdateCommand struct {
	baseCommand
	timeout time.Duration
	retry   time.Duration
}

func NewUpdateCommand() *UpdateCommand {

	cmd := &UpdateCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.DurationVar(&cmd.timeout, "t", 5*time.Second, "")
	cmd.flagSet.DurationVar(&cmd.timeout, "timeout", 5*time.Second, "")
	cmd.flagSet.DurationVar(&cmd.retry, "r", 0*time.Second, "")
	cmd.flagSet.DurationVar(&cmd.retry, "retry", 0*time.Second, "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *UpdateCommand) Apply(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	err = m.UpdateRemote(cmd.timeout, cmd.retry)
	if err != nil {
		return err
	}

	err = m.SaveFiles()
	if err != nil {
		return err
	}

	return nil
}
