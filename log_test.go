package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalLogger(t *testing.T) {
	testCases := []struct {
		name         string
		level        Level
		mode         WriteMode
		message      string
		expectedLine string
	}{
		{name: "ONE", level: Debug, mode: ModeNonBlocking, message: "hello", expectedLine: "|  INFO |     ONE | hello"},
		{name: "TWO", level: Debug, mode: ModeBlocking, message: "hello", expectedLine: "|  INFO |     TWO | hello"},
		{name: "THREE", level: Info, mode: ModeBlocking, message: "hello", expectedLine: "|  INFO |   THREE | hello"},
		{name: "FOUR", level: Error, mode: ModeBlocking, message: "hello", expectedLine: ""},
		{name: "FIVE", level: Fatal, mode: ModeBlocking, message: "hello", expectedLine: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w1 := &DummyWriter{}
			w2 := &DummyWriter{}

			l := New(
				WithName(tc.name),
				WithLevel(tc.level),
				WithMode(tc.mode),
				WithWriter(w1, FormatText),
				WithWriter(w2, FormatText),
			)

			l.Info(tc.message)
			l.Close()

			if tc.expectedLine != "" {
				require.NotEmpty(t, w1.Lines)
				require.NotEmpty(t, w2.Lines)
				assert.Contains(t, w1.Lines[0], tc.expectedLine)
				assert.Contains(t, w2.Lines[0], tc.expectedLine)
			} else {
				assert.Empty(t, w1.Lines)
				assert.Empty(t, w2.Lines)
			}
		})
	}
}

func TestLocalLogger(t *testing.T) {
	testCases := []struct {
		name         string
		level        Level
		message      string
		expectedLine string
	}{
		{name: "ONE", level: Debug, message: "world", expectedLine: "|  INFO |     ONE | world"},
		{name: "TWO", message: "world", expectedLine: ""},
		{name: "THREE", level: Fatal, message: "world", expectedLine: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w1 := &DummyWriter{}
			w2 := &DummyWriter{}

			g := New(
				WithLevel(Error),
				WithMode(ModeBlocking),
				WithWriter(w1, FormatText),
			)

			l := g.NewLocal(
				WithName(tc.name),
				WithWriter(w2, FormatText),
			)

			if tc.level != "" {
				l.SetLevel(tc.level)
			}

			l.Info(tc.message)

			if tc.expectedLine != "" {
				require.NotEmpty(t, w1.Lines)
				require.NotEmpty(t, w2.Lines)
				assert.Contains(t, w1.Lines[0], tc.expectedLine)
				assert.Contains(t, w2.Lines[0], tc.expectedLine)
			} else {
				assert.Empty(t, w1.Lines)
				assert.Empty(t, w2.Lines)
			}
		})
	}
}
