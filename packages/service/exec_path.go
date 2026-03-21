package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func buildCommand(name string, args ...string) (*exec.Cmd, error) {
	resolved, err := resolveExecutable(name)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(resolved, args...)
	cmd.Env = commandEnv()
	return cmd, nil
}

func resolveExecutable(name string) (string, error) {
	if strings.ContainsRune(name, filepath.Separator) {
		return name, nil
	}

	if resolved, err := exec.LookPath(name); err == nil {
		return resolved, nil
	}

	for _, dir := range commandSearchDirs() {
		candidate := filepath.Join(dir, name)
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("exec: %q: executable file not found in $PATH", name)
}

func commandEnv() []string {
	pathValue := strings.Join(commandSearchDirs(), string(os.PathListSeparator))
	env := os.Environ()

	for i, item := range env {
		if strings.HasPrefix(item, "PATH=") {
			env[i] = "PATH=" + pathValue
			return env
		}
	}

	return append(env, "PATH="+pathValue)
}

func commandSearchDirs() []string {
	seen := make(map[string]struct{})
	var dirs []string

	addDir := func(dir string) {
		if dir == "" {
			return
		}
		if _, ok := seen[dir]; ok {
			return
		}
		seen[dir] = struct{}{}
		dirs = append(dirs, dir)
	}

	for _, dir := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		addDir(dir)
	}

	if runtime.GOOS == "darwin" {
		addDir("/opt/homebrew/bin")
		addDir("/usr/local/bin")
		addDir("/opt/homebrew/sbin")
	}

	addDir("/usr/bin")
	addDir("/bin")
	addDir("/usr/sbin")
	addDir("/sbin")

	return dirs
}
