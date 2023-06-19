package command

import (
	"github.com/StackAdapt/systags/manager"
)

var Commands = map[string]Command{
	"help":    NewHelpCommand(),
	"init":    NewInitCommand(),
	"dump":    NewDumpCommand(),
	"update":  NewUpdateCommand(),
	"ls":      NewLsCommand(),
	"get":     NewGetCommand(),
	"set":     NewSetCommand(),
	"rm":      NewRmCommand(),
	"version": NewVersionCommand(),
}

func ParseArgs(m *manager.Manager, args []string) error {

	// Minimum arguments
	if len(args) >= 2 {

		cmd, found := Commands[args[1]]

		if found {
			// Attempt to initialize requested command
			if err := cmd.Parse(args[2:]); err != nil {
				logger.Error(err.Error())
				return err
			}

			// Attempt to apply requested command
			if err := cmd.Apply(m); err != nil {
				logger.Error(err.Error())
				return err
			}

			return nil
		}
	}

	// Print out all the top-level usage instructions
	if err := Commands["help"].Apply(m); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
