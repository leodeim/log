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
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
)

type _formatter struct{}

var formatter = _formatter{}

type formatterProps struct {
	level      *Level
	message    *string
	format     *Format
	name       *string
	dateFormat *string
}

func (_formatter) Get(props *formatterProps) (string, error) {
	switch *props.format {
	case FormatText:
		return formatter.text(props)
	case FormatTextColor:
		return formatter.textColor(props)
	case FormatJson:
		return formatter.json(props)
	default:
		return "", fmt.Errorf("incorrect log format: %v", *props.format)
	}
}

func (_formatter) text(props *formatterProps) (string, error) {
	name := *props.name
	if len(name) > 7 {
		name = name[:7]
	}

	return fmt.Sprintf(
		TextLogFormat,
		time.Now().Format(*props.dateFormat),
		*props.level,
		name,
		*props.message,
	), nil
}

func (_formatter) textColor(props *formatterProps) (string, error) {
	name := *props.name
	if len(name) > 7 {
		name = name[:7]
	}

	return fmt.Sprintf(
		ColorTextLogFormat,
		time.Now().Format(*props.dateFormat),
		levelToColor(*props.level),
		blue(name),
		*props.message,
	), nil
}

func levelToColor(level Level) string {
	switch level {
	case Fatal:
		return red(level)
	case Error:
		return red(level)
	case Warning:
		return yellow(level)
	case Info:
		return green(level)
	case Debug:
		return magenta(level)
	default:
		return string(level)
	}
}

func (_formatter) json(props *formatterProps) (string, error) {
	b, err := json.Marshal(map[string]string{
		"time":    time.Now().Format(*props.dateFormat),
		"level":   string(*props.level),
		"module":  *props.name,
		"message": *props.message,
	})

	if err != nil {
		return "", fmt.Errorf("error while formatting log message: %v", err)
	}

	return string(b), nil
}
