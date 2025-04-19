package main

import (
	"os"
	"os/exec"
	"runtime"
)

func main() {
	bin := "git-cleaner"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin)
	cmd.Env = append(cmd.Environ(), "CGO_ENABLED=1")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
