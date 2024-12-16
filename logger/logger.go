package logger

import (
	"flag"
	"log/slog"
	"os"
	"testing"
)

type Logger = *slog.Logger

const (
	LevelTrace     = slog.Level(-8)
	LevelDebug     = slog.LevelDebug
	LevelInfo      = slog.LevelInfo
	LevelNotice    = slog.Level(2)
	LevelWarning   = slog.LevelWarn
	LevelError     = slog.LevelError
	LevelEmergency = slog.Level(12)
)

const (
	Trace     = "TRACE"
	Debug     = "DEBUG"
	Info      = "INFO"
	Notice    = "NOTICE"
	Warning   = "WARNING"
	Error     = "ERROR"
	Emergency = "EMERGENCY"
)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	// Remove time from the output for predictable test output.
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}

	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Rename the level key from "level" to "sev".
		a.Key = "sev"

		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		// This could also look up the name from a map or other structure, but
		// this demonstrates using a switch statement to rename levels. For
		// maximum performance, the string values should be constants, but this
		// example uses the raw strings for readability.
		switch {
		case level < LevelDebug:
			a.Value = slog.StringValue(Trace)
		case level < LevelInfo:
			a.Value = slog.StringValue(Debug)
		case level < LevelNotice:
			a.Value = slog.StringValue(Info)
		case level < LevelWarning:
			a.Value = slog.StringValue(Notice)
		case level < LevelError:
			a.Value = slog.StringValue(Warning)
		case level < LevelEmergency:
			a.Value = slog.StringValue(Error)
		default:
			a.Value = slog.StringValue(Emergency)
		}
	}

	return a
}

func New(level slog.Level) Logger {
	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceAttr,
	})
	return slog.New(th)
}

func Default() Logger {
	return slog.Default()
}

func init() {

	defaultTraceLog := Info
	if flag.Lookup("verbose") != nil {
		defaultTraceLog = Trace
	}

	traceLog := flag.String("trace-log", defaultTraceLog, "enable trace logging")

	testing.Init()
	flag.Parse()
	var traceLevel slog.Level

	if traceLog != nil {
		switch *traceLog {
		case Trace:
			traceLevel = LevelTrace
		case Debug:
			traceLevel = LevelDebug
		case Info:
			traceLevel = LevelInfo
		case Notice:
			traceLevel = LevelNotice
		case Warning:
			traceLevel = LevelWarning
		case Error:
			traceLevel = LevelError
		case Emergency:
			traceLevel = LevelEmergency
		}
	}

	logger := New(traceLevel)
	slog.SetDefault(logger)
}
