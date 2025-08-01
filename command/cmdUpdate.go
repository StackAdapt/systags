package command

import (
	"flag"
	"strings"
	"time"

	"github.com/StackAdapt/systags/manager"
)

type UpdateCommand struct {
	baseCommand
	timeout time.Duration
	retry   time.Duration
	keys    string
}

func NewUpdateCommand() *UpdateCommand {

	cmd := &UpdateCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.DurationVar(&cmd.timeout, "t", 5*time.Second, "")
	cmd.flagSet.DurationVar(&cmd.timeout, "timeout", 5*time.Second, "")
	cmd.flagSet.DurationVar(&cmd.retry, "r", 0*time.Second, "") // TODO: Does this get overwritten by `--retry` default?
	cmd.flagSet.DurationVar(&cmd.retry, "retry", 0*time.Second, "")
	cmd.flagSet.StringVar(&cmd.keys, "k", "", "")
	cmd.flagSet.StringVar(&cmd.keys, "keys", "", "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *UpdateCommand) Apply(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	var keys []string
	if cmd.keys != "" {
		keys = strings.Split(cmd.keys, ",")
	}

	err = m.UpdateRemote(cmd.timeout, cmd.retry, keys)
	if err != nil {
		return err
	}

	err = m.SaveFiles()
	if err != nil {
		return err
	}

	return nil
}
