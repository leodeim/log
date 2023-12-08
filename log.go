package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Logger interface {
	NewLocal(opts ...Op) Logger
	SetLevel(level Level) error
	Close()
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
	Fatal   Level = "FATAL"
	Error   Level = "ERROR"
	Warning Level = "WARN"
	Info    Level = "INFO"
	Debug   Level = "DEBUG"
)

var levels = map[Level]int{
	Debug:   0,
	Info:    1,
	Warning: 2,
	Error:   3,
	Fatal:   4,
}

const (
	DefaultLevel      = Info
	DefaultName       = "<...>"
	DefaultFormat     = FormatTextColor
	DefaultWriteMode  = ModeBlocking
	DefaultDateFormat = "2006/01/02 15:04:05"
)

type globalProps struct {
	writers    []*writer
	mode       WriteMode
	dateFormat string
	sync.WaitGroup
	buf chan *logLine
}

type localProps struct {
	writers []*writer
	name    string
	level   int
}

type writer struct {
	writer io.Writer
	format Format
	sync.Mutex
}

type Op func(*globalProps, *localProps)

// WithName (local logger option), set logger/module name up to 7 characters (default: "<...>")
func WithName(n string) Op {
	return func(gp *globalProps, lp *localProps) {
		if lp == nil {
			return
		}
		lp.name = n
	}
}

// WithLevel (local logger option), set logger/module log level (default: Info)
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

// WithMode (global logger option), set blocking or non blocking logger mode (default: log.ModeBlocking)
func WithMode(m WriteMode) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp == nil {
			return
		}
		gp.mode = m
	}
}

// WithDateFormat (global logger option), set date/time format (default: "2006/01/02 15:04:05")
func WithDateFormat(f string) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp == nil {
			return
		}
		gp.dateFormat = f
	}
}

// WithWriter (global/local logger option), set custom log writer (default: os.Stdout)
// Several WithWriter options could be added
func WithWriter(w io.Writer, f Format) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp != nil {
			gp.writers = append(gp.writers, &writer{writer: w, format: f})
		} else if lp != nil {
			lp.writers = append(lp.writers, &writer{writer: w, format: f})
		}
	}
}

type log struct {
	global *globalProps
	local  *localProps
}

// New main logger instance, accepts both global and local options:
// WithName(string), WithLevel(log.Level), WithMode(log.WriteMode), WithDateFormat(string), WithWriter(io.Writer, log.Format)
func New(opts ...Op) Logger {
	gp := &globalProps{
		dateFormat: DefaultDateFormat,
		mode:       DefaultWriteMode,
	}

	lp := &localProps{
		name:  DefaultName,
		level: levels[DefaultLevel],
	}

	for _, opt := range opts {
		opt(gp, lp)
	}

	if len(gp.writers) == 0 {
		gp.writers = append(gp.writers, &writer{
			writer: os.Stdout,
			format: DefaultFormat,
		})
	}

	l := &log{
		global: gp,
		local:  lp,
	}

	if l.global.mode == ModeNonBlocking {
		l.run()
	}

	return l
}

// New local logger instance, accepts local options:
// WithName(string), WithLevel(log.Level), WithWriter(io.Writer, log.Format)
func (l *log) NewLocal(opts ...Op) Logger {
	lp := *l.local

	for _, opt := range opts {
		opt(nil, &lp)
	}

	return &log{
		global: l.global,
		local:  &lp,
	}
}

// Set log level for current logger instance
func (l *log) SetLevel(level Level) error {
	v, err := parseLevel(level)
	if err != nil {
		return err
	}

	l.local.level = v
	return nil
}

// Close logger, should be closed before application exit in case of non blocking mode
func (l *log) Close() {
	l.global.Wait()
}

func parseLevel(level Level) (int, error) {
	if v, ok := levels[Level(strings.ToUpper(string(level)))]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("incorrect log level type: %s", string(level))
}
