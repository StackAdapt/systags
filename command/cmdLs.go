package command

import (
	"errors"
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type LsCommand struct {
	baseCommand
	regex  bool
	pick   string
	omit   string
	format string
}

func NewLsCommand() *LsCommand {

	cmd := &LsCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.BoolVar(&cmd.regex, "r", false, "")
	cmd.flagSet.BoolVar(&cmd.regex, "regex", false, "")
	cmd.flagSet.StringVar(&cmd.pick, "p", "", "")
	cmd.flagSet.StringVar(&cmd.pick, "pick", "", "")
	cmd.flagSet.StringVar(&cmd.omit, "o", "", "")
	cmd.flagSet.StringVar(&cmd.omit, "omit", "", "")
	cmd.flagSet.StringVar(&cmd.format, "f", "json", "")
	cmd.flagSet.StringVar(&cmd.format, "format", "json", "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *LsCommand) Parse(args []string) error {

	err := cmd.flagSet.Parse(args)
	if err != nil {
		return err
	}

	if cmd.format == "" {
		return errors.New("flag needs to be provided: -format")
	}

	// If the specified format is supported
	_, found := manager.Formats[cmd.format]
	if !found {
		return errors.New("flag has unsupported value: -format")
	}

	return nil
}

func (cmd *LsCommand) Apply(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	// Retrieve the specified format method
	format, _ := manager.Formats[cmd.format]

	// Attempt to format tags
	out, err := m.FormatTags(
		cmd.regex, cmd.pick, cmd.omit, format,
	)

	if err != nil {
		return err
	}

	m.GetLogger().Info(out)

	return nil
}
