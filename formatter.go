package log

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Format int

const (
	FormatText Format = iota
	FormatTextColor
	FormatJson
)

const (
	ColorTextLogFormat = "%s | %14s | %16s | %s"
	TextLogFormat      = "%s | %5s | %7s | %s"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

type _formatter struct{}

var formatter = _formatter{}

func (_formatter) Text(l *log, level Level, message string) (string, error) {
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
}

func (_formatter) TextColor(l *log, level Level, message string) (string, error) {
	name := l.local.name
	if len(name) > 7 {
		name = name[:7]
	}
	return fmt.Sprintf(
		ColorTextLogFormat,
		yellow(time.Now().Format(l.global.dateFormat)),
		green(level),
		blue(name),
		message,
	), nil
}

func (_formatter) Json(l *log, level Level, message string) (string, error) {
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
}
