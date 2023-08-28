package log

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	blue  = color.New(color.FgBlue).SprintFunc()
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
			green(level),
			blue(name),
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
