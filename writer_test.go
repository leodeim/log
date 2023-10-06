package log

type DummyWriter struct {
	Lines []string
}

func (w *DummyWriter) Write(p []byte) (n int, err error) {
	w.Lines = append(w.Lines, string(p))
	return len(p), nil
}
