package log

import "fmt"

type message struct {
	super *log
	level Level
	text  string
	props map[string]string
}

func (m *message) Msg(text string) {
	m.text = text
	m.super.processor.Do(m)
}

func (m *message) Msgf(format string, args ...any) {
	m.text = fmt.Sprintf(format, args...)
	m.super.processor.Do(m)
}

func (m *message) Prop(key string, value string) {
	m.props[key] = value
}
