package main

import (
	"os"

	"github.com/StackAdapt/systags/command"
	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
)

func main() {

	m := manager.NewManager()

	configDir := os.Getenv("SYSTAGS_CONFIG_DIR")
	systemDir := os.Getenv("SYSTAGS_SYSTEM_DIR")

	if configDir != "" {
		m.ConfigDir = configDir
	}

	if systemDir != "" {
		m.SystemDir = systemDir
	}

	help := true
	// Ensure min arguments
	if len(os.Args) >= 2 {

		// Attempt to find the associated command
		cmd, found := command.Commands[os.Args[1]]

		if found {
			// Attempt to initialize the requested command
			if err := cmd.Init(os.Args[2:]); err != nil {
				// No need to print since Init does that
				os.Exit(1)
			}

			// Attempt to execute the command
			if err := cmd.Run(m); err != nil {
				utility.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			help = false
		}
	}

	if help {
		// Print out top-level usage instructions
		if err := command.Commands["help"].Run(m); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}
