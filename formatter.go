package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arrayFlags) contains(str string) bool {
	for _, j := range *i {
		if strings.Contains(str, j) {
			return true
		}
	}

	return false
}

var dirToScan, excludeDirs arrayFlags

func main() {

	var defaultExclude = []string{
		".git",
		".idea",
		"vendor",
	}

	flag.Var(&dirToScan, "dir", "dir with go files to format (multi-flag)")
	flag.Var(&excludeDirs, "exclude", "dir to exclude (multi-flag)")
	flag.Parse()

	for _, dir := range dirToScan {
		for _, excl := range defaultExclude {
			excludeDirs = append(excludeDirs, filepath.Join(dir, excl))
		}
	}

	fmt.Println("staring go multi-formatter")
	fmt.Println("dirs to scan (may be added with -dir flag):")
	for _, dir := range dirToScan {
		fmt.Println("----", dir)
	}

	fmt.Println("excluded dirs (may be added with -exclude flag):")
	for _, dir := range excludeDirs {
		fmt.Println("----", dir)
	}

	for _, dir := range dirToScan {
		err := processPath(dir, excludeDirs)
		if err != nil {
			panic(err)
		}
	}
}

func processPath(rootPath string, excludePaths arrayFlags) error {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if excludePaths.contains(path) {
			return nil
		}

		ext := filepath.Ext(info.Name())

		if !info.IsDir() && ext != ".go" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		return formatFile(path)
	})

	if err != nil {
		return err
	}

	return nil
}

func formatFile(path string) error {
	fmt.Printf("formatting file: %s\n", path)

	commands := map[string][]string{
		"gci":       {"write", path},
		"gofumpt":   {"-w", path},
		"goimports": {"-w", path},
		"gofmt":     {"-w", path},
	}
	for c, args := range commands {
		err := run(c, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func run(cmd string, args []string) error {
	var command = exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
