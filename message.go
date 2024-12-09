package log

import "fmt"

type message struct {
	super *log
	level Level
	text  string
	props map[string]interface{}
}

func (m *message) msg(text string, props ...any) {
	m.text = text
	m.handleProps(props...)
	m.super.processor.Do(m)
}

func (m *message) msgf(format string, args ...any) {
	m.text = fmt.Sprintf(format, args...)
	m.super.processor.Do(m)
}

func (m *message) handleProps(props ...any) {
	switch len(props) {
	case 0:
		return
	case 1:
		m.props["prop"] = props[0]
	default:
		for i := 0; i < len(props); i += 2 {
			key, ok := props[i].(string)
			if ok {
				m.props[key] = props[i+1]
			}
		}
	}
}
