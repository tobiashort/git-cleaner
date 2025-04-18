//go:build windows

package isatty

//#include <stdio.h>
//#include <conio.h>
import "C"

import (
	"os"
)

func IsTerminal(f *os.File) bool {
	return IsTerminalFd(int(f.Fd()))
}

func IsTerminalFd(fd int) bool {
	return C._isatty(C.int(fd)) == 1
}
