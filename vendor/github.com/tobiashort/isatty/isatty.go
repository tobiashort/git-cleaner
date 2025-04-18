//go:build linux || darwin

package isatty

//#include <unistd.h>
import "C"

import (
	"os"
)

func IsTerminal(f *os.File) bool {
	return IsTerminalFd(int(f.Fd()))
}

func IsTerminalFd(fd int) bool {
	return int(C.isatty(C.int(fd))) == 1
}
