package isatty

import (
	"os"

	"golang.org/x/term"
)

func IsTerminal(f *os.File) bool {
	return IsTerminalFd(int(f.Fd()))
}

func IsTerminalFd(fd int) bool {
	return term.IsTerminal(fd)
}
