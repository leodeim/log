package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger interface {
	Local(opts ...Op) Logger
	SetLevel(level Level) error
	Sync()
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

type Format int

const (
	FormatText Format = iota
	FormatJson
)

type WriteMode int

const (
	ModeNonBlocking WriteMode = iota
	ModeBlocking
)

const (
	TextLogFormat = "%s | %7s | %7s | %s"

	DefaultLevel      = Info
	DefaultName       = "<...>"
	DefaultFormat     = FormatText
	DefaultWriteMode  = ModeBlocking
	DefaultDateFormat = "2006/01/02 15:04:05"
)

type log struct {
	global *globalProps
	local  *localProps
}

type globalProps struct {
	writers    []io.Writer
	writeMode  WriteMode
	format     Format
	dateFormat string
	wg         sync.WaitGroup
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

func WithWriteMode(m WriteMode) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp == nil {
			return
		}
		gp.writeMode = m
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

func WithWriter(w io.Writer) Op {
	return func(gp *globalProps, lp *localProps) {
		if gp == nil {
			return
		}
		gp.writers = append(gp.writers, w)
	}
}

func New(opts ...Op) Logger {
	gp := &globalProps{
		format:     DefaultFormat,
		dateFormat: DefaultDateFormat,
		writeMode:  DefaultWriteMode,
	}

	lp := &localProps{
		name:  DefaultName,
		level: levels[DefaultLevel],
	}

	for _, opt := range opts {
		opt(gp, lp)
	}

	if len(gp.writers) == 0 {
		gp.writers = append(gp.writers, os.Stdout)
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
	l.write(Info, fmt.Sprintf(format, v...))
}

func (l *log) Info(message string) {
	l.write(Info, message)
}

func (l *log) Errorf(format string, v ...any) {
	l.write(Error, fmt.Sprintf(format, v...))
}

func (l *log) Error(message string) {
	l.write(Error, message)
}

func (l *log) Warningf(format string, v ...any) {
	l.write(Warning, fmt.Sprintf(format, v...))
}

func (l *log) Warning(message string) {
	l.write(Warning, message)
}

func (l *log) Debugf(format string, v ...any) {
	l.write(Debug, fmt.Sprintf(format, v...))
}

func (l *log) Debug(message string) {
	l.write(Debug, message)
}

func (l *log) Fatalf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.write(Fatal, message)
	panic(message)
}

func (l *log) Fatal(message string) {
	l.write(Fatal, message)
	panic(message)
}

func (l *log) write(level Level, message string) {
	if v, ok := levels[level]; !ok || v < l.local.level {
		return
	}

	if len(message) == 0 {
		return
	}

	log, err := l.formatter(level, message)
	if err != nil {
		return
	}

	if log[len(log)-1] != '\n' {
		log = log + "\n"
	}

	for _, w := range l.global.writers {
		switch l.global.writeMode {
		case ModeBlocking:
			fmt.Fprint(w, log)
		case ModeNonBlocking:
			l.global.wg.Add(1)
			go func(w io.Writer) {
				defer l.global.wg.Done()
				fmt.Fprint(w, log)
			}(w)
		}
	}
}

func (l *log) Sync() {
	l.global.wg.Wait()
}

func (l *log) formatter(level Level, message string) (string, error) {
	switch l.global.format {
	case FormatText:
		name := l.local.name
		if len(name) > 7 {
			name = name[:7]
		}
		return fmt.Sprintf(
			TextLogFormat,
			time.Now().Format(l.global.dateFormat),
			level,
			name,
			message,
		), nil
	case FormatJson:
		b, err := json.Marshal(map[string]string{
			"time":    time.Now().Format(l.global.dateFormat),
			"level":   string(level),
			"module":  l.local.name,
			"message": message,
		})
		if err != nil {
			return "", fmt.Errorf("error while formatting the log message: %v", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("incorrect log format: %v", l.global.format)
	}
}

func parseLevel(level Level) (int, error) {
	if v, ok := levels[Level(strings.ToUpper(string(level)))]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("incorrect log level type: %s", string(level))
}
