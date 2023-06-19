package main

import (
	"context"
	"log"
	"os"

	"github.com/fatih/color"
	"golang.org/x/exp/slog"

	"github.com/StackAdapt/systags/command"
	"github.com/StackAdapt/systags/manager"
)

type LogHandler struct {
	slog.Handler
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func (log *LogHandler) Handle(_ context.Context, r slog.Record) error {

	switch r.Level {
	case slog.LevelInfo:
		log.info.Println(r.Message)

	case slog.LevelWarn:
		log.warn.Println(color.YellowString(r.Message))

	case slog.LevelError:
		log.error.Println(color.RedString(r.Message))
	}

	return nil
}

func main() {

	logger := slog.New(&LogHandler{
		Handler: slog.NewTextHandler(os.Stdout, nil),
		info:    log.New(os.Stdout, "", 0),
		warn:    log.New(os.Stderr, "", 0),
		error:   log.New(os.Stderr, "", 0),
	})

	m := manager.NewManager()

	configDir := os.Getenv("SYSTAGS_CONFIG_DIR")
	systemDir := os.Getenv("SYSTAGS_SYSTEM_DIR")

	if configDir != "" {
		m.ConfigDir = configDir
	}

	if systemDir != "" {
		m.SystemDir = systemDir
	}

	command.SetLogger(logger)
	// Perform CLI parsing, errors are logged using logger
	if err := command.ParseArgs(m, os.Args); err != nil {
		os.Exit(1)
	}
}
