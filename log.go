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

var levels = map[Level]int{
	levelDebug:   0,
	levelInfo:    1,
	levelWarning: 2,
	levelError:   3,
	levelFatal:   4,
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

type log struct {
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
		if v, err := parseLevel(l); err == nil {
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

func New(opts ...Op) Logger {
	gp := &globalProps{
		format:     DefaultFormat,
		dateFormat: DefaultDateFormat,
	}

	lp := &localProps{
		name:  DefaultName,
		level: levels[DefaultLevel],
	}

	for _, opt := range opts {
		opt(gp, lp)
	}

	return &log{
		global: gp,
		local:  lp,
	}
}

func (l *log) Local(opts ...Op) Logger {
	lp := *l.local

	for _, opt := range opts {
		opt(nil, &lp)
	}

	return &log{
		global: l.global,
		local:  &lp,
	}
}

func (l *log) SetLevel(level Level) error {
	v, err := parseLevel(level)
	if err != nil {
		return err
	}

	l.local.level = v
	return nil
}

func (l *log) Infof(format string, v ...any) {
	l.write(levelInfo, fmt.Sprintf(format, v...))
}

func (l *log) Info(message string) {
	l.write(levelInfo, message)
}

func (l *log) Errorf(format string, v ...any) {
	l.write(levelError, fmt.Sprintf(format, v...))
}

func (l *log) Error(message string) {
	l.write(levelError, message)
}

func (l *log) Warningf(format string, v ...any) {
	l.write(levelWarning, fmt.Sprintf(format, v...))
}

func (l *log) Warning(message string) {
	l.write(levelWarning, message)
}

func (l *log) Debugf(format string, v ...any) {
	l.write(levelDebug, fmt.Sprintf(format, v...))
}

func (l *log) Debug(message string) {
	l.write(levelDebug, message)
}

func (l *log) Fatalf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.write(levelFatal, message)
	panic(message)
}

func (l *log) Fatal(message string) {
	l.write(levelFatal, message)
	panic(message)
}

func (l *log) write(level Level, message string) {
	if v, ok := levels[level]; !ok || v < l.local.level {
		return
	}

	if len(message) == 0 {
		return
	}

	if message[len(message)-1] != '\n' {
		message = message + "\n"
	}

	date := time.Now().Format(l.global.dateFormat)
	fmt.Fprintf(os.Stdout, TextLogFormat, date, level, l.local.name, message)
}

func parseLevel(level Level) (int, error) {
	if v, ok := levels[Level(strings.ToUpper(string(level)))]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("incorrect log level type: %s", string(level))
}
