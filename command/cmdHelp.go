package command

import (
	"flag"

	"github.com/StackAdapt/systags/manager"
)

type HelpCommand struct {
	baseCommand
}

func NewHelpCommand() *HelpCommand {

	cmd := &HelpCommand{
		baseCommand: baseCommand{
			flagSet: flag.NewFlagSet("", flag.ContinueOnError),
		},
	}

	// NOTE: Help command should not include any additional flags

	// Don't print unneeded usage
	cmd.flagSet.Usage = func() {}

	return cmd
}

func (cmd *HelpCommand) Apply(m *manager.Manager) error {

	m.GetLogger().Info("TODO: USAGE")

	/*
		Commands:
			help
				<none>

			init
				-r|--reset (bool) [optional = false]

			dump
				-k|--kind (string) [required]
					config, remote, system

			update
				-t|--timeout (duration) [optional = 5*time.Second]

			ls
				-r|--regex (bool) [optional = false]
				-f|--format (string) [optional = json]
					json, consul, env, telegraf
				-p|--pick (string) [optional = ""]
				-o|--omit (string) [optional = ""]

			get
				-k|--key     (string) [required]
				-d|--default (string) [optional = ""]

			set
				-k|--key   (string) [required]
				-v|--value (string) [required]

			rm
				-k|--key (string) [required]

			version
				<none>

		Env:
			SYSTAGS_CONFIG_DIR
				=/etc/systags.d

			SYSTAGS_SYSTEM_DIR
				=/var/lib/systags

			SYSTAGS_DEBUG
				=
	*/

	return nil
}
