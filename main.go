package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type ExecutionResult struct {
	path   string
	output string
	err    error
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2[T any](v T, err error) T {
	must(err)
	return v
}

func sliceContains[T comparable](slice []T, search T) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}
	return false
}

func findGitRepositories() []string {
	var paths []string
	wd := must2(os.Getwd())
	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == ".git" {
			paths = append(paths, filepath.Dir(path))
		}
		return nil
	})
	return paths
}

func run(cmd *exec.Cmd) (output string, err error) {
	bs, err := cmd.CombinedOutput()
	output = string(bs)
	if err != nil {
		return output, err
	}
	return output, nil
}

func runGit(path string, args ...string) ExecutionResult {
	args = append([]string{"-C", path}, args...)
	cmd := exec.Command("git", args...)
	output, err := run(cmd)
	return ExecutionResult{path, output, err}
}

func gitResetHard(path string) ExecutionResult {
	return runGit(path, "reset", "--hard")
}

func gitListBranches(path string) (branches []string, executionResult ExecutionResult) {
	executionResult = runGit(path, "branch", "--no-color")
	if executionResult.err != nil {
		return nil, executionResult
	}
	for _, line := range strings.Split(executionResult.output, "\n") {
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "* ", "", 1)
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches, executionResult
}

func gitCheckoutMaster(path string) ExecutionResult {
	branches, executionResult := gitListBranches(path)
	if executionResult.err != nil {
		return executionResult
	}
	var branch string
	if sliceContains(branches, "master") {
		branch = "master"
	} else if sliceContains(branches, "main") {
		branch = "main"
	}
	if branch == "" {
		return ExecutionResult{
			executionResult.path,
			executionResult.output,
			fmt.Errorf("no master/main branch found"),
		}
	}
	return runGit(path, "checkout", branch)
}

func gitClean(path string) ExecutionResult {
	return runGit(path, "clean", "-fd")
}

func gitPull(path string) ExecutionResult {
	return runGit(path, "pull", "-p")
}

func gitRemoveLocalBranches(path string) ExecutionResult {
	branches, executionResult := gitListBranches(path)
	if executionResult.err != nil {
		return executionResult
	}
	for _, branch := range branches {
		if branch == "master" || branch == "main" {
			continue
		}
		executionResult = runGit(path, "branch", "-D", branch)
		if executionResult.err != nil {
			return executionResult
		}
	}
	return ExecutionResult{path, "", nil}
}

func cleanGitRepository(path string, waitGroup *sync.WaitGroup) {
	var executionResult ExecutionResult
	executionResult = gitResetHard(path);           if executionResult.err != nil { goto errorCase }
	executionResult = gitCheckoutMaster(path);      if executionResult.err != nil { goto errorCase }
	executionResult = gitClean(path);               if executionResult.err != nil { goto errorCase }
	executionResult = gitPull(path);                if executionResult.err != nil { goto errorCase }
	executionResult = gitRemoveLocalBranches(path); if executionResult.err != nil { goto errorCase }
	fmt.Println("[DONE]", path)
	waitGroup.Done()
	return
errorCase:
	fmt.Println("[ERROR]", path)
	fmt.Println(executionResult.err)
	fmt.Println(executionResult.output)
	waitGroup.Done()
	return
}

func main() {
	gitRepositories := findGitRepositories()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(gitRepositories))
	for _, path := range gitRepositories {
		go cleanGitRepository(path, &waitGroup)
	}
	waitGroup.Wait()
}
