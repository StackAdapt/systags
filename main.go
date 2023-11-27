package main

import (
	"log/slog"
	"os"

	"github.com/StackAdapt/systags/command"
	"github.com/StackAdapt/systags/manager"
	"github.com/StackAdapt/systags/utility"
)

func main() {

	m := manager.NewManager()

	var level slog.Level
	if os.Getenv("SYSTAGS_DEBUG") == "" {
		level = slog.LevelInfo
	} else {
		level = slog.LevelDebug
	}

	logger := utility.NewLogger(level)
	m.SetLogger(logger)

	configDir := os.Getenv("SYSTAGS_CONFIG_DIR")
	systemDir := os.Getenv("SYSTAGS_SYSTEM_DIR")

	if configDir != "" {
		m.ConfigDir = configDir
	}

	if systemDir != "" {
		m.SystemDir = systemDir
	}

	// Perform CLI parsing, errors are logged using logger
	if err := command.ParseArgs(m, os.Args); err != nil {
		os.Exit(1)
	}
}
