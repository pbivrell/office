package util

// Fields is a map of arbitrary key-values to be logged
type Fields map[string]interface{}

// Level is the representation of logging level
type Level uint32

// Constants defining log levels
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

// Entry defines what functions Entry loggers must implement
type Entry interface {
	WithError(err error) Entry
	WithField(key string, value interface{}) Entry
	WithFields(fields Fields) Entry

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Logger defines what functions loggers themselves must implement
type Logger interface {
	WithError(err error) Entry
	WithField(key string, value interface{}) Entry
	WithFields(fields Fields) Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	SetLevel(level Level)
}

// Formatter defines what functions formatters must implement
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
