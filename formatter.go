package log

import (
	"encoding/json"
	"fmt"
	"strings"
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
	ColorTextLogFormat = "%s | %14s | %16s | %s %s"
	TextLogFormat      = "%s | %5s | %7s | %s %s"
	TextPropFormat     = " | %s=%s"
)

var (
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
)

type _formatter struct{}

var formatter = _formatter{}

type formatterProps struct {
	message *message
	format  Format
}

func (_formatter) Get(props *formatterProps) (string, error) {
	switch props.format {
	case FormatText:
		return formatter.text(props)
	case FormatTextColor:
		return formatter.textColor(props)
	case FormatJson:
		return formatter.json(props)
	default:
		return "", fmt.Errorf("incorrect log format: %v", props.format)
	}
}

func (_formatter) text(props *formatterProps) (string, error) {
	logger := props.message.super

	name := logger.local.name
	if len(name) > 7 {
		name = name[:7]
	}

	msgProps := messageProps(props.message.props)

	return fmt.Sprintf(
		TextLogFormat,
		time.Now().Format(logger.global.dateFormat),
		props.message.level,
		name,
		props.message.text,
		msgProps,
	), nil
}

func messageProps(props map[string]string) string {
	all := strings.Builder{}

	for k, v := range props {
		all.WriteString(fmt.Sprintf(TextPropFormat, k, v))
	}

	return all.String()
}

func (_formatter) textColor(props *formatterProps) (string, error) {
	logger := props.message.super

	name := logger.local.name
	if len(name) > 7 {
		name = name[:7]
	}

	msgProps := messagePropsColor(props.message.props)

	return fmt.Sprintf(
		ColorTextLogFormat,
		time.Now().Format(logger.global.dateFormat),
		levelToColor(props.message.level),
		blue(name),
		props.message.text,
		msgProps,
	), nil
}

func messagePropsColor(props map[string]string) string {
	all := strings.Builder{}

	for k, v := range props {
		all.WriteString(fmt.Sprintf(TextPropFormat, k, cyan(v)))
	}

	return all.String()
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
	logger := props.message.super

	b, err := json.Marshal(map[string]string{
		"time":    time.Now().Format(logger.global.dateFormat),
		"level":   string(props.message.level),
		"module":  logger.local.name,
		"message": props.message.text,
	})

	if err != nil {
		return "", fmt.Errorf("error while formatting log message: %v", err)
	}

	return string(b), nil
}
