package tabpp

import (
	"bytes"
	"io"
	"os"
	"text/tabwriter"
)

type Writer interface {
	io.Writer
	WriteTab() (int, error)
	Flush() error
}

type plainWriter struct{ io.Writer }

func (w *plainWriter) WriteTab() (int, error) {
	return w.Writer.Write([]byte{' '})
}

func (w *plainWriter) Flush() error { return nil }

type tabWriter struct {
	*tabwriter.Writer
	innerw io.Writer
}

func newTabWriter(innerw io.Writer) *tabWriter {
	return &tabWriter{
		Writer: tabwriter.NewWriter(
			innerw, 0, 0, 1, ' ', 0,
		),
		innerw: innerw,
	}
}

func (w *tabWriter) WriteTab() (int, error) {
	return w.Writer.Write([]byte{'\t'})
}

type alternativeWriter struct {
	multiw  io.Writer
	plainw  *plainWriter
	tabw    *tabWriter
	targetw io.Writer
	width   int
}

func Wrap(targetw *os.File) Writer {
	ncol, err := ncolumns(targetw)
	if err != nil {
		// not a tty, then just use tabWriter
		return newTabWriter(targetw)

	}

	plainw := &plainWriter{new(bytes.Buffer)}
	tabw := newTabWriter(new(bytes.Buffer))

	return &alternativeWriter{
		multiw:  io.MultiWriter(plainw, tabw),
		plainw:  plainw,
		tabw:    tabw,
		targetw: targetw,
		width:   ncol,
	}
}

func (w *alternativeWriter) Write(p []byte) (int, error) {
	return w.multiw.Write(p)
}

func (w *alternativeWriter) WriteTab() (n int, err error) {
	n, err = w.plainw.WriteTab()
	if err != nil {
		return
	}
	return w.tabw.WriteTab()
}

func (w *alternativeWriter) Flush() (err error) {
	w.plainw.Flush()
	w.tabw.Flush()
	plainb := w.plainw.Writer.(*bytes.Buffer)
	tabb := w.tabw.innerw.(*bytes.Buffer)

	fit, total := bufStats(tabb, w.width)
	if fit <= total/2 {
		_, err = io.Copy(w.targetw, plainb)
	} else {
		_, err = io.Copy(w.targetw, tabb)
	}
	return
}

func bufStats(b *bytes.Buffer, width int) (linesFit, linesTotal int) {
	lines := bytes.Split(b.Bytes(), []byte{'\n'})
	trimLastIfEmpty(&lines)
	linesTotal = len(lines)
	for _, l := range lines {
		if len(l) < width {
			linesFit++
		}
	}
	return
}

func trimLastIfEmpty(bb *[][]byte) {
	if len(*bb) == 0 {
		return
	}
	l := len(*bb) - 1
	b := (*bb)[l]
	if len(b) == 0 {
		*bb = (*bb)[:l]
	}
}
