package utility

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

// ColorWriter is a custom writer that adds color to the output.
type ColorWriter struct {
	w     *os.File
	color *color.Color
}

func (cw *ColorWriter) Write(p []byte) (n int, err error) {
	return cw.w.Write([]byte(cw.color.Sprint(string(p))))
}

// LogHandler is a custom slog handler that routes log messages to different
// loggers based on their severity level. It embeds slog.Handler and defines
// loggers for different severity levels like Debug, Info, Warn, and Error.
type LogHandler struct {
	slog.Handler
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

func (log *LogHandler) Handle(_ context.Context, r slog.Record) error {

	switch r.Level {
	case slog.LevelDebug:
		log.Debug.Println(r.Message)

	case slog.LevelInfo:
		log.Info.Println(r.Message)

	case slog.LevelWarn:
		log.Warn.Println(r.Message)

	case slog.LevelError:
		log.Error.Println(r.Message)
	}

	return nil
}

// NewLogger creates a new instance of slog.Logger. The logging level can be
// specified by level. The returned logger writes info messages to STDOUT and
// warning, error, and debug messages to STDERR. The messages are color-coded
// for readability: debug messages are light gray, warning messages are yellow,
// and error messages are red.
func NewLogger(level slog.Level) *slog.Logger {

	options := &slog.HandlerOptions{
		Level: level,
	}

	cwDebug := &ColorWriter{w: os.Stderr, color: color.New(color.FgHiBlack)}
	cwInfo := &ColorWriter{w: os.Stdout, color: color.New()}
	cwWarn := &ColorWriter{w: os.Stderr, color: color.New(color.FgYellow)}
	cwError := &ColorWriter{w: os.Stderr, color: color.New(color.FgRed)}

	logger := slog.New(&LogHandler{
		Handler: slog.NewTextHandler(os.Stdout, options),

		Debug: log.New(cwDebug, "", 0),
		Info:  log.New(cwInfo, "", 0),
		Warn:  log.New(cwWarn, "", 0),
		Error: log.New(cwError, "", 0),
	})

	return logger
}
