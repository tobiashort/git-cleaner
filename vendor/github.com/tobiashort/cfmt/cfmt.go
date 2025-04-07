package cfmt

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type AnsiColor = string

const (
	AnsiRed    AnsiColor = "\033[0;31m"
	AnsiGreen  AnsiColor = "\033[0;32m"
	AnsiYellow AnsiColor = "\033[1;33m"
	AnsiBlue   AnsiColor = "\033[1;34m"
	AnsiPurple AnsiColor = "\033[1;35m"
	AnsiCyan   AnsiColor = "\033[1;36m"
	AnsiReset  AnsiColor = "\033[0m"
)

var regexps = map[*regexp.Regexp]AnsiColor{
	makeRegexp("r"): AnsiRed,
	makeRegexp("g"): AnsiGreen,
	makeRegexp("y"): AnsiYellow,
	makeRegexp("b"): AnsiBlue,
	makeRegexp("p"): AnsiPurple,
	makeRegexp("c"): AnsiCyan,
}

func makeRegexp(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("#%s\\{([^}]*)\\}", name))
}

func Print(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), AnsiReset)
	}
	fmt.Print(a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(clr(format, AnsiReset), a...)
}

func Println(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), AnsiReset)
	}
	fmt.Println(a...)
}

func Fprint(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), AnsiReset)
	}
	fmt.Fprint(w, a...)
}

func Fprintf(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, clr(format, AnsiReset), a...)
}

func Fprintln(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), AnsiReset)
	}
	fmt.Fprintln(w, a...)
}

func stoc(s string) AnsiColor {
	switch s {
	case "r":
		return AnsiRed
	case "g":
		return AnsiGreen
	case "y":
		return AnsiYellow
	case "b":
		return AnsiBlue
	case "p":
		return AnsiPurple
	case "c":
		return AnsiCyan
	default:
		panic(fmt.Errorf("cannot map string '%s' to ansi color", s))
	}
}

func CPrint(s string, a ...any) {
	c := stoc(s)
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), c)
	}
	fmt.Print(c)
	fmt.Print(a...)
	fmt.Print(AnsiReset)
}

func CPrintf(s string, format string, a ...any) {
	c := stoc(s)
	fmt.Print(c)
	fmt.Printf(clr(format, c), a...)
	fmt.Print(AnsiReset)
}

func CPrintln(s string, a ...any) {
	c := stoc(s)
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), c)
	}
	fmt.Print(c)
	fmt.Println(a...)
	fmt.Print(AnsiReset)
}

func clr(str string, reset AnsiColor) string {
	for regex, color := range regexps {
		matches := regex.FindAllStringSubmatch(str, -1)
		for _, match := range matches {
			str = strings.Replace(str, match[0], color+match[1]+reset, 1)
		}
	}
	return str
}
