package command

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
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

func (cmd *DumpCommand) Init(args []string) error {

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
		return cmd.failf("flag needs to be provided: -kind")

	default:
		return cmd.failf("flag has unsupported value: -kind")
	}

	return nil
}

func (cmd *DumpCommand) Run(m *manager.Manager) error {

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

	utility.Fprintln(
		os.Stdout,
		string(out),
	)

	return nil
}
