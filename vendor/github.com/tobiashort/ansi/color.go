package ansi

type Color = string

const (
	Red    Color = "\033[0;31m"
	Green  Color = "\033[0;32m"
	Yellow Color = "\033[1;33m"
	Blue   Color = "\033[1;34m"
	Purple Color = "\033[1;35m"
	Cyan   Color = "\033[1;36m"
	Reset  Color = "\033[0m"
)
