package util

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// LogrusLogger logrus implementation of a logger
type LogrusLogger struct {
	internalLog *logrus.Logger
}

// LogrusEntry logrus implementation of an entry
type LogrusEntry struct {
	internalEntry *logrus.Entry
}

// Debug logs a message at level Debug
func (logger *LogrusLogger) Debug(args ...interface{}) {
	logger.internalLog.Debug(args...)
}

// Debugf logs a message at level Debug with a format string
func (logger *LogrusLogger) Debugf(format string, args ...interface{}) {
	logger.internalLog.Debugf(format, args...)
}

// Info logs a message at level Info
func (logger *LogrusLogger) Info(args ...interface{}) {
	logger.internalLog.Info(args...)
}

// Infof logs a message at level Info with a format string
func (logger *LogrusLogger) Infof(format string, args ...interface{}) {
	logger.internalLog.Infof(format, args...)
}

// Warn logs a message at level Warn
func (logger *LogrusLogger) Warn(args ...interface{}) {
	logger.internalLog.Warn(args...)
}

// Warnf logs a message at level Warn with a format string
func (logger *LogrusLogger) Warnf(format string, args ...interface{}) {
	logger.internalLog.Warnf(format, args...)
}

// Error logs a message at level Error
func (logger *LogrusLogger) Error(args ...interface{}) {
	logger.internalLog.Error(args...)
}

// Errorf logs a message at level Error with a format string
func (logger *LogrusLogger) Errorf(format string, args ...interface{}) {
	logger.internalLog.Errorf(format, args...)
}

// Fatal logs a message at level Fatal
func (logger *LogrusLogger) Fatal(args ...interface{}) {
	logger.internalLog.Fatal(args...)
}

// Fatalf logs a message at level Fatal with a format string
func (logger *LogrusLogger) Fatalf(format string, args ...interface{}) {
	logger.internalLog.Fatalf(format, args...)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined as a key
func (logger *LogrusLogger) WithError(err error) Entry {
	return &LogrusEntry{
		internalEntry: logger.internalLog.WithError(err),
	}
}

// WithField creates an entry from the standard logger and adds a field to it
func (logger *LogrusLogger) WithField(key string, value interface{}) Entry {
	return &LogrusEntry{
		internalEntry: logger.internalLog.WithField(key, value),
	}
}

// WithFields creates an entry from the standard logger and adds multiple fields to it
func (logger *LogrusLogger) WithFields(fields Fields) Entry {
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}

	return &LogrusEntry{
		internalEntry: logger.internalLog.WithFields(data),
	}
}

// SetLevel sets the standard logger level
func (logger *LogrusLogger) SetLevel(level Level) {
	logger.internalLog.SetLevel(logrus.Level(level))
}

// Writer LogrusLogger implementation that returns the write half of the pipe
func (logger *LogrusLogger) Writer() *io.PipeWriter {
	return logger.internalLog.Writer()
}

// Debug LogrusEntry implementation that logs a message at level Debug
func (entry *LogrusEntry) Debug(args ...interface{}) {
	entry.internalEntry.Debug(args...)
}

// Debugf LogrusEntry implementation that logs a message at level Debug with a format string
func (entry *LogrusEntry) Debugf(format string, args ...interface{}) {
	entry.internalEntry.Debugf(format, args...)
}

// Info LogrusEntry implementation that logs a message at level Info
func (entry *LogrusEntry) Info(args ...interface{}) {
	entry.internalEntry.Info(args...)
}

// Infof LogrusEntry implementation that logs a message at level Info with a format string
func (entry *LogrusEntry) Infof(format string, args ...interface{}) {
	entry.internalEntry.Infof(format, args...)
}

// Warn LogrusEntry implementation that logs a message at level Warn
func (entry *LogrusEntry) Warn(args ...interface{}) {
	entry.internalEntry.Warn(args...)
}

// Warnf LogrusEntry implementation that logs a message at level Warn with a format string
func (entry *LogrusEntry) Warnf(format string, args ...interface{}) {
	entry.internalEntry.Warnf(format, args...)
}

// Error LogrusEntry implementation that logs a message at level Error
func (entry *LogrusEntry) Error(args ...interface{}) {
	entry.internalEntry.Error(args...)
}

// Errorf LogrusEntry implementation that logs a message at level Error with a format string
func (entry *LogrusEntry) Errorf(format string, args ...interface{}) {
	entry.internalEntry.Errorf(format, args...)
}

// Fatal LogrusEntry implementation that logs a message at level Fatal
func (entry *LogrusEntry) Fatal(args ...interface{}) {
	entry.internalEntry.Fatal(args...)
}

// Fatalf LogrusEntry implementation that logs a message at level Fatal with a format string
func (entry *LogrusEntry) Fatalf(format string, args ...interface{}) {
	entry.internalEntry.Fatalf(format, args...)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined as a key
func (entry *LogrusEntry) WithError(err error) Entry {
	return &LogrusEntry{
		internalEntry: entry.internalEntry.WithError(err),
	}
}

// WithField creates an entry from a LogrusEntry and adds a field to it
func (entry *LogrusEntry) WithField(key string, value interface{}) Entry {
	return &LogrusEntry{
		internalEntry: entry.internalEntry.WithField(key, value),
	}
}

// WithFields creates an entry from a LogrusEntry and adds multiple fields to it
func (entry *LogrusEntry) WithFields(fields Fields) Entry {
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}

	return &LogrusEntry{
		internalEntry: entry.internalEntry.WithFields(data),
	}
}

var logLevel logrus.Level

func init() {
	logLevel = logrus.InfoLevel
}

// SetGlobalLevel sets the level for all new loggers
func SetGlobalLevel(level Level) {
	logLevel = logrus.Level(level)
}

// NewLogrusLogger instantiates a new logrus type logger
func NewLogrusLogger() Logger {
	logrusInst := logrus.New()
	//logrusInst.Formatter = &logrus.JSONFormatter{}
	logrusInst.Out = os.Stdout
	logrusInst.Level = logLevel
	return &LogrusLogger{
		internalLog: logrusInst,
	}
}

// NewLogrusDiscardLogger instantiates a new logrus type logger that has output going to /dev/null
func NewLogrusDiscardLogger() Logger {
	logrusInst := logrus.New()
	logrusInst.Out = ioutil.Discard
	logrusInst.Level = logLevel
	return &LogrusLogger{
		internalLog: logrusInst,
	}
}
