package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type hook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *hook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, writer := range hook.Writer {
		writer.Write([]byte(line))
	}

	return err
}

func (hook *hook) Levels() []logrus.Level {
	return hook.LogLevels
}

var entry *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func NewLogger() Logger {
	return Logger{entry}
}

func (l *Logger) NewLoggerWithField(key string, value interface{}) Logger {
	return Logger{l.WithField(key, value)}
}

func init() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)

			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", fileName, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	log.SetOutput(io.Discard)

	log.AddHook(&hook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	log.SetLevel(logrus.TraceLevel)

	entry = logrus.NewEntry(log)
}
