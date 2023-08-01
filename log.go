package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Local(opts ...Op) Logger
	SetLevel(level Level) error
	Infof(format string, v ...any)
	Info(message string)
	Errorf(format string, v ...any)
	Error(message string)
	Warningf(format string, v ...any)
	Warning(message string)
	Debugf(format string, v ...any)
	Debug(message string)
	Fatalf(format string, v ...any)
	Fatal(message string)
}

type Level string

const (
	levelFatal   Level = "FATAL"
	levelError   Level = "ERROR"
	levelWarning Level = "WARN"
	levelInfo    Level = "INFO"
	levelDebug   Level = "DEBUG"
)

var logLevels = map[Level]int{
	levelFatal:   0,
	levelError:   1,
	levelWarning: 2,
	levelInfo:    3,
	levelDebug:   4,
}

type Format int

const (
	FormatText Format = iota
	FormatJson
)

const (
	TextLogFormat = "%s | %7s | %7s | %s"

	DefaultLevel      = levelInfo
	DefaultName       = "<...>"
	DefaultFormat     = FormatText
	DefaultDateFormat = "2006/01/02 15:04:05"
)

type logger struct {
	global *globalProps
	local  *localProps
}

type globalProps struct {
	format     Format
	dateFormat string
}

type localProps struct {
	name  string
	level int
}

type Op func(*globalProps, *localProps)

func WithName(n string) Op {
	return func(gp *globalProps, lp *localProps) {
		if lp == nil {
			return
		}
		lp.name = n
	}
}

func WithLevel(l Level) Op {
	return func(gp *globalProps, lp *localProps) {
		if lp == nil {
			return
		}
		v, err := parseLevel(l)
		if err == nil {
			lp.level = v
		}
	}
}

func WithFormat(f Format) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp == nil {
			return
		}
		gp.format = f
	}
}

func WithDateFormat(f string) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp == nil {
			return
		}
		gp.dateFormat = f
	}
}

// name -> module name, will be visible in log message
func New(opts ...Op) Logger {
	gp := globalDefaults()
	lp := localDefaults()

	for _, opt := range opts {
		opt(gp, lp)
	}

	return &logger{
		global: gp,
		local:  lp,
	}
}

func (l *logger) Local(opts ...Op) Logger {
	lp := localDefaults()

	for _, opt := range opts {
		opt(nil, lp)
	}

	return &logger{
		global: l.global,
		local:  lp,
	}
}

func globalDefaults() *globalProps {
	return &globalProps{
		format:     DefaultFormat,
		dateFormat: DefaultDateFormat,
	}
}

func localDefaults() *localProps {
	return &localProps{
		name:  DefaultName,
		level: logLevels[DefaultLevel],
	}
}

func (l *logger) SetLevel(level Level) error {
	v, err := parseLevel(level)
	if err != nil {
		return err
	}

	l.local.level = v
	return nil
}

func parseLevel(level Level) (int, error) {
	if v, ok := logLevels[Level(strings.ToUpper(string(level)))]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("abort SetGlobalLogLevel -> bad log level type: %s", string(level))
}

func (l *logger) Infof(format string, v ...any) {
	l.write(levelInfo, fmt.Sprintf(format, v...))
}

func (l *logger) Info(message string) {
	l.write(levelInfo, message)
}

func (l *logger) Errorf(format string, v ...any) {
	l.write(levelError, fmt.Sprintf(format, v...))
}

func (l *logger) Error(message string) {
	l.write(levelError, message)
}

func (l *logger) Warningf(format string, v ...any) {
	l.write(levelWarning, fmt.Sprintf(format, v...))
}

func (l *logger) Warning(message string) {
	l.write(levelWarning, message)
}

func (l *logger) Debugf(format string, v ...any) {
	l.write(levelDebug, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(message string) {
	l.write(levelDebug, message)
}

func (l *logger) Fatalf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.write(levelFatal, message)
	panic(message)
}

func (l *logger) Fatal(message string) {
	l.write(levelFatal, message)
	panic(message)
}

func (l *logger) write(level Level, message string) {
	if v, ok := logLevels[level]; ok && v <= l.local.level {
		var (
			len  = len(message)
			date = time.Now().Format(DefaultDateFormat)
		)

		if len > 0 && message[len-1] != '\n' {
			message = message + "\n"
		}

		fmt.Fprintf(os.Stdout, TextLogFormat, date, level, l.local.name, message)
	}
}
