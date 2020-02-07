// +build !windows,!plan9,!solaris

package tabpp

import (
	"os"

	"golang.org/x/sys/unix"
)

func ncolumns(f *os.File) (int, error) {
	ws, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return -1, err
	}
	return int(ws.Col), nil
}
