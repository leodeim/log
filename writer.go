package log

import (
	"fmt"
)

const buffSize = 100

type WriteMode int

const (
	ModeNonBlocking WriteMode = iota
	ModeBlocking
)

type logLine struct {
	msg   string
	level Level
}

func (l *log) Infof(format string, v ...any) {
	l.log(logLine{level: Info, msg: fmt.Sprintf(format, v...)})
}

func (l *log) Info(message string) {
	l.log(logLine{level: Info, msg: message})
}

func (l *log) Errorf(format string, v ...any) {
	l.log(logLine{level: Error, msg: fmt.Sprintf(format, v...)})
}

func (l *log) Error(message string) {
	l.log(logLine{level: Error, msg: message})
}

func (l *log) Warningf(format string, v ...any) {
	l.log(logLine{level: Warning, msg: fmt.Sprintf(format, v...)})
}

func (l *log) Warning(message string) {
	l.log(logLine{level: Warning, msg: message})
}

func (l *log) Debugf(format string, v ...any) {
	l.log(logLine{level: Debug, msg: fmt.Sprintf(format, v...)})
}

func (l *log) Debug(message string) {
	l.log(logLine{level: Debug, msg: message})
}

func (l *log) Fatalf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	l.log(logLine{level: Fatal, msg: message})
	panic(message)
}

func (l *log) Fatal(message string) {
	l.log(logLine{level: Fatal, msg: message})
	panic(message)
}

func (l *log) log(line logLine) {
	switch l.global.mode {
	case ModeBlocking:
		l.write(line)
	case ModeNonBlocking:
		select {
		case l.global.buf <- &line:
		default:
		}
	}
}

func (l *log) write(line logLine) {
	v, ok := levels[line.level]
	if !ok {
		return
	}

	if v < l.local.level {
		return
	}

	if len(line.msg) == 0 {
		return
	}

	for _, w := range l.global.writers {
		m, err := l.format(w.format, line.level, line.msg)
		if err != nil {
			continue
		}

		l.writeByMode(w, l.global.mode, m)
	}

	for _, w := range l.local.writers {
		m, err := l.format(w.format, line.level, line.msg)
		if err != nil {
			continue
		}

		l.writeByMode(w, l.global.mode, m)
	}
}

func (l *log) run() {
	l.global.buf = make(chan *logLine, buffSize)

	go func() {
		for {
			select {
			case line := <-l.global.buf:
				if line != nil {
					l.write(*line)
				}
			}
		}
	}()
}

func (l *log) writeByMode(w *writer, mode WriteMode, msg string) {
	switch mode {
	case ModeBlocking:
		l.writeWithLock(w, msg)
	case ModeNonBlocking:
		l.writeDirect(w, msg)
	}
}

func (l *log) writeDirect(w *writer, msg string) {
	fmt.Fprintln(w.writer, msg)
}

func (l *log) writeWithLock(w *writer, msg string) {
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
