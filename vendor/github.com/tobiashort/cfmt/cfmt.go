package cfmt

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/tobiashort/ansi"
	"github.com/tobiashort/isatty"
)

var regexps = map[*regexp.Regexp]ansi.Color{
	makeRegexp("r"): ansi.Red,
	makeRegexp("g"): ansi.Green,
	makeRegexp("y"): ansi.Yellow,
	makeRegexp("b"): ansi.Blue,
	makeRegexp("p"): ansi.Purple,
	makeRegexp("c"): ansi.Cyan,
}

func makeRegexp(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("#%s\\{([^}]*)\\}", name))
}

func Print(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.Reset)
	}
	fmt.Print(a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(clr(format, ansi.Reset), a...)
}

func Println(a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.Reset)
	}
	fmt.Println(a...)
}

func Fprint(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.Reset)
	}
	fmt.Fprint(w, a...)
}

func Fprintf(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, clr(format, ansi.Reset), a...)
}

func Fprintln(w io.Writer, a ...any) {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.Reset)
	}
	fmt.Fprintln(w, a...)
}

func Sprint(a ...any) string {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.Reset)
	}
	return fmt.Sprint(a...)
}

func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(clr(format, ansi.Reset), a...)
}

func Sprintln(a ...any) string {
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), ansi.Reset)
	}
	return fmt.Sprintln(a...)
}

func stoc(s string) ansi.Color {
	switch s {
	case "r":
		return ansi.Red
	case "g":
		return ansi.Green
	case "y":
		return ansi.Yellow
	case "b":
		return ansi.Blue
	case "p":
		return ansi.Purple
	case "c":
		return ansi.Cyan
	default:
		panic(fmt.Errorf("cannot map string '%s' to ansi color", s))
	}
}

func CPrint(s string, a ...any) {
	c := stoc(s)
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), c)
	}
	if isatty.IsTerminal() {
		fmt.Print(c)
	}
	fmt.Print(a...)
	if isatty.IsTerminal() {
		fmt.Print(ansi.Reset)
	}
}

func CPrintf(s string, format string, a ...any) {
	c := stoc(s)
	if isatty.IsTerminal() {
		fmt.Print(c)
	}
	fmt.Printf(clr(format, c), a...)
	if isatty.IsTerminal() {
		fmt.Print(ansi.Reset)
	}
}

func CPrintln(s string, a ...any) {
	c := stoc(s)
	for i := range a {
		a[i] = clr(fmt.Sprint(a[i]), c)
	}
	if isatty.IsTerminal() {
		fmt.Print(c)
	}
	fmt.Println(a...)
	if isatty.IsTerminal() {
		fmt.Print(ansi.Reset)
	}
}

func clr(str string, reset ansi.Color) string {
	for regex, color := range regexps {
		matches := regex.FindAllStringSubmatch(str, -1)
		for _, match := range matches {
			if isatty.IsTerminal() {
				str = strings.Replace(str, match[0], color+match[1]+reset, 1)
			} else {
				str = strings.Replace(str, match[0], match[1], 1)
			}
		}
	}
	return str
}
