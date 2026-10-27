package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCLIInstallCopyPreservesSource exercises the service without using the
// machine's home directory, environment, or registry.
func TestCLIInstallCopyPreservesSource(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX startup-file behavior is exercised on POSIX")
	}
	home := t.TempDir()
	source := filepath.Join(t.TempDir(), "original")
	if err := os.WriteFile(source, []byte("original executable"), 0755); err != nil {
		t.Fatal(err)
	}
	service := CLIInstallService{OS: "linux", Home: home, Executable: source, Env: map[string]string{"SHELL": "/bin/sh"}}
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".local", "bin", "zesh-skills")); !os.IsNotExist(err) {
		t.Fatalf("planning must not write destination: %v", err)
	}
	if err := service.Apply(plan, CLIInstallConsent{Install: true}); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{source, filepath.Join(home, ".local", "bin", "zesh-skills")} {
		got, err := os.ReadFile(file)
		if err != nil || string(got) != "original executable" {
			t.Fatalf("%s: content %q, error %v", file, got, err)
		}
	}
	profile, err := os.ReadFile(filepath.Join(home, ".profile"))
	if err != nil || !strings.Contains(string(profile), ".local/bin") {
		t.Fatalf("profile %q, error %v", profile, err)
	}
}

// TestCLIInstallWindowsPathPreview uses a fake user-PATH store on Windows.
func TestCLIInstallWindowsPathPreview(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows filesystem paths require a Windows host")
	}
	store := &installTestPath{value: `C:\Existing`, valueType: 2}
	service := CLIInstallService{OS: "windows", Home: t.TempDir(), Executable: filepath.Join(t.TempDir(), "source.exe"), Env: map[string]string{"LOCALAPPDATA": t.TempDir()}, WindowsPath: store}
	if err := os.WriteFile(service.Executable, []byte("source"), 0755); err != nil {
		t.Fatal(err)
	}
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if store.writes != 0 || store.value != `C:\Existing` {
		t.Fatalf("planning mutated PATH: %+v", store)
	}
	if err := service.Apply(plan, CLIInstallConsent{Install: true}); err != nil {
		t.Fatal(err)
	}
	if store.writes != 1 || store.valueType != 2 || !strings.HasPrefix(store.value, `C:\Existing;`) {
		t.Fatalf("unexpected PATH after apply: %+v", store)
	}
}

type installTestPath struct {
	value     string
	valueType uint32
	writes    int
}

func (p *installTestPath) Read() (string, uint32, error) { return p.value, p.valueType, nil }
func (p *installTestPath) Write(value string, valueType uint32) error {
	p.value, p.valueType = value, valueType
	p.writes++
	return nil
}
