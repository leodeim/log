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

	if v < m.super.local.level {
		return
	}

	if len(m.text) == 0 {
		return
	}

	for _, w := range m.super.global.writers {
		str, err := formatter.Get(&formatterProps{m, w.format})
		if err != nil {
			continue
		}

		p.writeByMode(w, m.super.global.mode, str)
	}

	for _, w := range m.super.local.writers {
		str, err := formatter.Get(&formatterProps{m, w.format})
		if err != nil {
			continue
		}

		p.writeByMode(w, m.super.global.mode, str)
	}
}

func (p *processor) writeByMode(w *writer, mode WriteMode, msg string) {
	switch mode {
	case ModeBlocking:
		p.writeWithLock(w, msg)
	case ModeNonBlocking:
		p.writeDirect(w, msg)
	}
}

func (p *processor) writeDirect(w *writer, msg string) {
	fmt.Fprintln(w.writer, msg)
}

func (p *processor) writeWithLock(w *writer, msg string) {
	w.mu.Lock()
	fmt.Fprintln(w.writer, msg)
	w.mu.Unlock()
}
