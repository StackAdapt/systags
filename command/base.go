package command

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
)

type Command interface {
	Init([]string) error
	Run(m *manager.Manager) error
	Usage() string
}

type baseCommand struct {
	flagSet *flag.FlagSet
}

func (cmd *baseCommand) Init(args []string) error {
	return cmd.flagSet.Parse(args)
}

func (cmd *baseCommand) Run(_ *manager.Manager) error {
	return nil
}

func (cmd *baseCommand) Usage() string {

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
	*/

	return "TODO: USAGE"
}

func (cmd *baseCommand) failf(format string, a ...any) error {

	msg := fmt.Sprintf(format, a...)
	utility.Fprintln(os.Stderr, msg)
	cmd.flagSet.Usage()
	return errors.New(msg)
}
