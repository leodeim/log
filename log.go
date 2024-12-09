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
	Level() Level
	Close()
	Info(text string, props ...any)
	Infof(text string, args ...any)
	Error(text string, props ...any)
	Errorf(text string, args ...any)
	Err(text string, err error)
	Warning(text string, props ...any)
	Warningf(text string, args ...any)
	Warn(text string, err error)
	Debug(text string, props ...any)
	Debugf(text string, args ...any)
	Fatal(text string, err error)
	Fatalf(text string, args ...any)
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

var levelsRev = map[int]Level{
	0: Debug,
	1: Info,
	2: Warning,
	3: Error,
	4: Fatal,
}

const (
	DefaultLevel      = Info
	DefaultName       = "<...>"
	DefaultFormat     = FormatTextColor
	DefaultWriteMode  = ModeBlocking
	DefaultDateFormat = "2006/01/02 15:04:05"
)

type WriteMode int

const (
	ModeNonBlocking WriteMode = 1 + iota
	ModeBlocking
)

type globalProps struct {
	writers    []*writer
	mode       WriteMode
	dateFormat string
}

type localProps struct {
	writers []*writer
	name    string
	level   int
}

type writer struct {
	writer io.Writer
	format Format
	mu     sync.Mutex
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
	global    *globalProps
	local     *localProps
	processor *processor
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
		global:    gp,
		local:     lp,
		processor: NewProcessor(gp.mode),
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
		global:    l.global,
		local:     &lp,
		processor: l.processor,
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

// Get log level for current logger instance
func (l *log) Level() Level {
	if v, ok := levelsRev[l.local.level]; ok {
		return v
	}

	return ""
}

// Close logger, should be closed before application exit in case of non blocking mode
func (l *log) Close() {
	defer func() {
		if recover() != nil {
			fmt.Println("log: error closing buffer channel")
		}
	}()

	close(l.processor.buf)
}

func (l *log) Info(text string, props ...any) {
	m := l.newMessage(Info)
	m.msg(text, props...)
}

func (l *log) Infof(text string, args ...any) {
	m := l.newMessage(Info)
	m.msgf(text, args...)
}

func (l *log) Error(text string, props ...any) {
	m := l.newMessage(Error)
	m.msg(text, props...)
}

func (l *log) Errorf(text string, args ...any) {
	m := l.newMessage(Error)
	m.msgf(text, args...)
}

func (l *log) Err(text string, err error) {
	m := l.newMessage(Error)
	if err == nil {
		m.msg(text)
	} else {
		m.msg(text, "error", err.Error())
	}
}

func (l *log) Warning(text string, props ...any) {
	m := l.newMessage(Warning)
	m.msg(text, props...)
}

func (l *log) Warningf(text string, args ...any) {
	m := l.newMessage(Warning)
	m.msgf(text, args...)
}

func (l *log) Warn(text string, err error) {
	m := l.newMessage(Warning)
	if err == nil {
		m.msg(text)
	} else {
		m.msg(text, "error", err.Error())
	}
}

func (l *log) Debug(text string, props ...any) {
	m := l.newMessage(Debug)
	m.msg(text, props...)
}

func (l *log) Debugf(text string, args ...any) {
	m := l.newMessage(Debug)
	m.msgf(text, args...)
}

func (l *log) Fatal(text string, err error) {
	m := l.newMessage(Fatal)
	if err == nil {
		m.msg(text)
	} else {
		m.msg(text, "error", err.Error())
	}
}

func (l *log) Fatalf(text string, args ...any) {
	m := l.newMessage(Fatal)
	m.msgf(text, args...)
}

func (l *log) newMessage(level Level) *message {
	return &message{
		super: l,
		level: level,
		props: make(map[string]interface{}),
	}
}

func parseLevel(level Level) (int, error) {
	if v, ok := levels[Level(strings.ToUpper(string(level)))]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("incorrect log level type: %s", string(level))
}
