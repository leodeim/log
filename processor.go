package log

import (
	"fmt"
)

const buffSize = 1024

type processor struct {
	mode WriteMode
	buf  chan *message
}

func NewProcessor(mode WriteMode) *processor {
	p := &processor{
		mode: mode,
	}

	if mode == ModeNonBlocking {
		p.run()
	}

	return p
}

func (p *processor) run() {
	p.buf = make(chan *message, buffSize)

	go func() {
		for msg := range p.buf {
			if msg != nil {
				p.write(msg)
			}
		}
	}()
}

func (p *processor) Do(m *message) {
	defer func() {
		if recover() != nil {
			fmt.Println("log: error writing to buffer")
		}
	}()

	switch p.mode {
	case ModeBlocking:
		p.write(m)
	case ModeNonBlocking:
		select {
		case p.buf <- m:
		default:
		}
	}
}

func (p *processor) write(m *message) {
	v, ok := levels[m.level]
	if !ok {
		return
	}

	logger := m.super

	if v < logger.local.level {
		return
	}

	if len(m.text) == 0 {
		return
	}

	for _, w := range logger.global.writers {
		str, err := formatter.Get(&formatterProps{m, w.format})
		if err != nil {
			continue
		}

		p.writeByMode(w, str)
	}

	for _, w := range logger.local.writers {
		str, err := formatter.Get(&formatterProps{m, w.format})
		if err != nil {
			continue
		}

		p.writeByMode(w, str)
	}
}

func (p *processor) writeByMode(w *writer, msg string) {
	switch p.mode {
	case ModeBlocking:
		p.writeSync(w, msg)
	case ModeNonBlocking:
		p.writeAsync(w, msg)
	}
}

func (p *processor) writeAsync(w *writer, msg string) {
	fmt.Fprintln(w.writer, msg)
}

func (p *processor) writeSync(w *writer, msg string) {
	w.mu.Lock()
	fmt.Fprintln(w.writer, msg)
	w.mu.Unlock()
}
