package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type RunGitResult struct {
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

func isGitRepository(path string) (isGitRepository bool, subdirectories []string) {
	for _, file := range must2(os.ReadDir(path)) {
		if !file.IsDir() {
			continue
		}
		if file.Name() == ".git" {
			return true, []string{}
		}
		subdirectory := filepath.Join(path, file.Name())
		subdirectories = append(subdirectories, subdirectory)
	}
	return false, subdirectories
}

func findGitRepositories() []string {
	currentDirectory := must2(os.Getwd())
	searchPaths := []string{currentDirectory}
	repositoryPaths := []string{}
	for len(searchPaths) > 0 {
		searchIndex := len(searchPaths) - 1
		searchPath := searchPaths[searchIndex]
		searchPaths = searchPaths[:searchIndex]
		isGitRepository, subdirectories := isGitRepository(searchPath)
		if isGitRepository {
			repositoryPaths = append(repositoryPaths, searchPath)
		} else {
			searchPaths = append(searchPaths, subdirectories...)
		}
	}
	return repositoryPaths
}

func run(cmd *exec.Cmd) (output string, err error) {
	bs, err := cmd.CombinedOutput()
	output = string(bs)
	if err != nil {
		return output, err
	}
	return output, nil
}

func runGit(path string, args ...string) RunGitResult {
	args = append([]string{"-C", path}, args...)
	cmd := exec.Command("git", args...)
	output, err := run(cmd)
	return RunGitResult{path, output, err}
}

func gitResetHard(path string) RunGitResult {
	return runGit(path, "reset", "--hard")
}

func gitListBranches(path string) (branches []string, runGitResult RunGitResult) {
	runGitResult = runGit(path, "branch", "--no-color")
	if runGitResult.err != nil {
		return nil, runGitResult
	}
	for _, line := range strings.Split(runGitResult.output, "\n") {
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "* ", "", 1)
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches, runGitResult
}

func gitCheckoutMaster(path string) RunGitResult {
	branches, runGitResult := gitListBranches(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	var branch string
	if sliceContains(branches, "master") {
		branch = "master"
	} else if sliceContains(branches, "main") {
		branch = "main"
	}
	if branch == "" {
		return RunGitResult{
			runGitResult.path,
			runGitResult.output,
			fmt.Errorf("no master/main branch found"),
		}
	}
	return runGit(path, "checkout", branch)
}

func gitClean(path string) RunGitResult {
	return runGit(path, "clean", "-fd")
}

func gitPull(path string) RunGitResult {
	return runGit(path, "pull", "-p")
}

func gitRemoveLocalBranches(path string) RunGitResult {
	branches, runGitResult := gitListBranches(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	for _, branch := range branches {
		if branch == "master" || branch == "main" {
			continue
		}
		runGitResult = runGit(path, "branch", "-D", branch)
		if runGitResult.err != nil {
			return runGitResult
		}
	}
	return RunGitResult{path, "", nil}
}

func cleanGitRepository(path string) RunGitResult {
	fmt.Println("Cleaning at", path)
	runGitResult := gitResetHard(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	runGitResult = gitCheckoutMaster(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	runGitResult = gitClean(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	runGitResult = gitPull(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	runGitResult = gitRemoveLocalBranches(path)
	if runGitResult.err != nil {
		return runGitResult
	}
	fmt.Println("Cleaned at", path)
	return RunGitResult{path, "", nil}
}

func main() {
	gitRepositories := findGitRepositories()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(gitRepositories))
	results := make(chan RunGitResult, len(gitRepositories))
	for _, path := range gitRepositories {
		go func(path string, results chan RunGitResult) {
			results <- cleanGitRepository(path)
			waitGroup.Done()
		}(path, results)
	}
	waitGroup.Wait()
	for {
		select {
		case result := <-results:
			if result.err != nil {
				fmt.Println("Error at", result.path)
				fmt.Println(result.err)
				fmt.Print(result.output)
			}
		default:
			goto done
		}
	}
done:
	close(results)
}
