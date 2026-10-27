//go:build windows

package main

import (
	"errors"
	"os"
	"runtime"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// newHostCLIInstallService detects the real host without mutating anything.
func newHostCLIInstallService() (CLIInstallService, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return CLIInstallService{}, err
	}
	executable, err := os.Executable()
	if err != nil {
		return CLIInstallService{}, err
	}
	env := map[string]string{}
	if value := os.Getenv("LOCALAPPDATA"); value != "" {
		env["LOCALAPPDATA"] = value
	}
	return CLIInstallService{OS: runtime.GOOS, Home: home, Executable: executable, Env: env, WindowsPath: windowsUserPath{}}, nil
}

// windowsUserPath reads and writes only HKCU\Environment, the current user's
// PATH. It never opens the machine PATH and never requires elevation.
type windowsUserPath struct{}

func (windowsUserPath) Read() (string, uint32, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", 0, err
	}
	defer key.Close()
	buffer := make([]byte, 4096)
	for {
		n, kind, err := key.GetValue("Path", buffer)
		if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) {
			return "", 0, nil
		}
		if errors.Is(err, windows.ERROR_MORE_DATA) {
			buffer = make([]byte, n)
			continue
		}
		if err != nil {
			return "", 0, err
		}
		if kind != registry.SZ && kind != registry.EXPAND_SZ {
			return "", kind, nil
		}
		if n%2 != 0 {
			return "", 0, errors.New("user PATH is not valid UTF-16")
		}
		units := make([]uint16, 0, n/2)
		for i := 0; i+1 < n; i += 2 {
			units = append(units, uint16(buffer[i])|uint16(buffer[i+1])<<8)
		}
		return windows.UTF16ToString(units), kind, nil
	}
}

// Write restores the previous registry value type; kind 0 means "absent".
func (windowsUserPath) Write(value string, kind uint32) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	switch kind {
	case 0:
		if err := key.DeleteValue("Path"); err != nil && !errors.Is(err, windows.ERROR_FILE_NOT_FOUND) {
			return err
		}
		return nil
	case registry.EXPAND_SZ:
		return key.SetExpandStringValue("Path", value)
	case registry.SZ:
		return key.SetStringValue("Path", value)
	}
	return errors.New("unsupported user PATH value type")
}
