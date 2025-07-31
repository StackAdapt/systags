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

type StringArray []string

// String is the method to format the flag's value, part of the flag.Value interface.
func (s *StringArray) String() string {
	return strings.Join(*s, ",")
}

// Set is the method to set the flag value, part of the flag.Value interface.
func (s *StringArray) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
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

	err = m.UpdateRemote(cmd.timeout, cmd.retry, strings.Split(cmd.keys, ","))
	if err != nil {
		return err
	}

	err = m.SaveFiles()
	if err != nil {
		return err
	}

	return nil
}
