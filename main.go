package main

import (
	"os"

	"github.com/bdlm/log"
	stdLogger "github.com/bdlm/std/logger"

	"github.com/FreifunkBremen/yanic/cmd"
)

type Hook struct{}

func (hook *Hook) Fire(entry *log.Entry) error {
	switch entry.Level {
	case log.PanicLevel:
		entry.Logger.Out = os.Stderr
	case log.FatalLevel:
		entry.Logger.Out = os.Stderr
	case log.ErrorLevel:
		entry.Logger.Out = os.Stderr
	case log.WarnLevel:
		entry.Logger.Out = os.Stdout
	case log.InfoLevel:
		entry.Logger.Out = os.Stdout
	case log.DebugLevel:
		entry.Logger.Out = os.Stdout
	default:
	}

	return nil
}

func (hook *Hook) Levels() []stdLogger.Level {
	return log.AllLevels
}

func main() {
	log.AddHook(&Hook{})

	cmd.Execute()
}
