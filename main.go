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
	out *log.Logger
	err *log.Logger
}

func (log *LogHandler) Handle(_ context.Context, r slog.Record) error {

	switch r.Level {
	case slog.LevelDebug:
		log.out.Println(color.HiBlackString(r.Message))

	case slog.LevelInfo:
		log.out.Println(r.Message)

	case slog.LevelWarn:
		log.err.Println(color.YellowString(r.Message))

	case slog.LevelError:
		log.err.Println(color.RedString(r.Message))
	}

	return nil
}

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

	debug := os.Getenv("SYSTAGS_DEBUG")

	var loggerOpts *slog.HandlerOptions

	if debug == "" {
		loggerOpts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	} else {
		loggerOpts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	logger := slog.New(&LogHandler{
		Handler: slog.NewTextHandler(os.Stdout, loggerOpts),
		out:     log.New(os.Stdout, "", 0),
		err:     log.New(os.Stderr, "", 0),
	})

	m.SetLogger(logger)

	// Perform CLI parsing, errors are logged using logger
	if err := command.ParseArgs(m, os.Args); err != nil {
		os.Exit(1)
	}
}
