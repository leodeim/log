package log

import (
	"fmt"
)

type WriteMode int

const (
	ModeNonBlocking WriteMode = iota
	ModeBlocking
)

func (l *log) Infof(format string, v ...any) {
	l.log(Info, fmt.Sprintf(format, v...))
}

func (l *log) Info(message string) {
	l.log(Info, message)
}

func (l *log) Errorf(format string, v ...any) {
	l.log(Error, fmt.Sprintf(format, v...))
}

func (l *log) Error(message string) {
	l.log(Error, message)
}

func (l *log) Warningf(format string, v ...any) {
	l.log(Warning, fmt.Sprintf(format, v...))
}

func (l *log) Warning(message string) {
	l.log(Warning, message)
}

func (l *log) Debugf(format string, v ...any) {
	l.log(Debug, fmt.Sprintf(format, v...))
}

func (l *log) Debug(message string) {
	l.log(Debug, message)
}

func (l *log) Fatalf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.log(Fatal, message)
	panic(message)
}

func (l *log) Fatal(message string) {
	l.log(Fatal, message)
	panic(message)
}

func (l *log) log(level Level, msg string) {
	v, ok := levels[level]
	if !ok {
		return
	}

	if v < l.local.level {
		return
	}

	if len(msg) == 0 {
		return
	}

	for _, w := range l.global.writers {
		m, err := l.format(w.format, level, msg)
		if err != nil {
			continue
		}

		l.writeByMode(w, l.global.mode, m)
	}

	for _, w := range l.local.writers {
		m, err := l.format(w.format, level, msg)
		if err != nil {
			continue
		}

		l.writeByMode(w, l.global.mode, m)
	}
}

func (l *log) writeByMode(w *writer, mode WriteMode, msg string) {
	switch mode {
	case ModeBlocking:
		l.write(w, msg)
	case ModeNonBlocking:
		l.global.Add(1)
		go func(w *writer) {
			defer l.global.Done()
			l.write(w, msg)
		}(w)
	}
}

func (l *log) write(w *writer, msg string) {
	w.Lock()
	defer w.Unlock()

	fmt.Fprintln(w.writer, msg)
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
