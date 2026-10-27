//go:build !windows

package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// newHostCLIInstallService detects the real host without mutating anything. No
// Windows user PATH store exists in this build, so registry state is unreachable.
func newHostCLIInstallService() (CLIInstallService, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return CLIInstallService{}, err
	}
	executable, err := os.Executable()
	if err != nil {
		return CLIInstallService{}, err
	}
	if resolved, err := filepath.EvalSymlinks(executable); err == nil {
		executable = resolved
	}
	env := map[string]string{}
	for _, name := range []string{"SHELL", "ZDOTDIR", "XDG_CONFIG_HOME"} {
		if value := os.Getenv(name); value != "" {
			env[name] = value
		}
	}
	return CLIInstallService{OS: runtime.GOOS, Home: home, Executable: executable, Env: env}, nil
}
