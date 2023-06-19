package command

import (
	"encoding/json"
	"errors"
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type DumpCommand struct {
	baseCommand
	kind string
}

func NewDumpCommand() *DumpCommand {

	cmd := &DumpCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	cmd.flagSet.StringVar(&cmd.kind, "k", "", "")
	cmd.flagSet.StringVar(&cmd.kind, "kind", "", "")

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *DumpCommand) Parse(args []string) error {

	err := cmd.flagSet.Parse(args)
	if err != nil {
		return err
	}

	switch cmd.kind {
	case "config":
	case "remote":
	case "system":
		break

	case "":
		return errors.New("flag needs to be provided: -kind")

	default:
		return errors.New("flag has unsupported value: -kind")
	}

	return nil
}

func (cmd *DumpCommand) Apply(m *manager.Manager) error {

	err := m.LoadFiles()
	if err != nil {
		return err
	}

	var tags manager.Tags

	switch cmd.kind {
	case "config":
		tags = m.ConfigTags()

	case "remote":
		tags = m.RemoteTags()

	case "system":
		tags = m.SystemTags()
	}

	// Attempt to convert the remote data to JSON
	out, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return err
	}

	logger.Info(string(out))

	return nil
}
