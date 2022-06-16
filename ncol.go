package tabpp

import (
	"os"

	"golang.org/x/term"
)

func ncolumns(f *os.File) (int, error) {
	w, _, err := term.GetSize(int(f.Fd()))
	if err != nil {
		return -1, err
	}
	return w, nil
}
