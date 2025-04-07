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
		a[i] = clr(fmt.Sprint(a[i]))
	}
	fmt.Print(a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(clr(format), a...)
}

func Println(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]))
	}
	fmt.Println(a...)
}

func Fprint(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]))
	}
	fmt.Fprint(w, a...)
}

func Fprintf(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, clr(format), a...)
}

func Fprintln(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]))
	}
	fmt.Fprintln(w, a...)
}

func Sprint(a ...any) string {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]))
	}
	return fmt.Sprint(a...)
}

func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(clr(format), a...)
}

func Sprintln(a ...any) string {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]))
	}
	return fmt.Sprintln(a...)
}

func clr(str string) string {
	for regex, color := range regexps {
		matches := regex.FindAllStringSubmatch(str, -1)
		for _, match := range matches {
			str = strings.Replace(str, match[0], color+match[1]+AnsiReset, 1)
		}
	}
	return str
}
