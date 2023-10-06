package log

import (
	"fmt"
	"io"
)

type WriteMode int

const (
	ModeNonBlocking WriteMode = iota
	ModeBlocking
)

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
	v, ok := levels[level]
	if !ok {
		return
	}

	if v < l.local.level {
		return
	}

	if len(message) == 0 {
		return
	}

	for _, w := range l.global.writers {
		m, err := l.format(w.format, level, message)
		if err != nil {
			continue
		}

		l.out(w.writer, l.global.mode, m)
	}

	for _, w := range l.local.writers {
		m, err := l.format(w.format, level, message)
		if err != nil {
			continue
		}

		l.out(w.writer, l.global.mode, m)
	}
}

func (l *log) out(writer io.Writer, mode WriteMode, message string) {
	switch mode {
	case ModeBlocking:
		fmt.Fprintln(writer, message)
	case ModeNonBlocking:
		l.global.Add(1)
		go func(w io.Writer) {
			defer l.global.Done()
			fmt.Fprintln(w, message)
		}(writer)
	}
}

func (l *log) format(format Format, level Level, message string) (string, error) {
	switch format {
	case FormatText:
		return formatter.text(l, level, message)
	case FormatTextColor:
		return formatter.textColor(l, level, message)
	case FormatJson:
		return formatter.json(l, level, message)
	default:
		return "", fmt.Errorf("incorrect log format: %v", format)
	}
}
