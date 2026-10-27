package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCLIInstallPreviewCannotBypassReplacementConsent(t *testing.T) {
	s := installFixture(t, "linux", "/bin/sh")
	target := filepath.Join(s.Home, ".local", "bin", "zesh-skills")
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("previous executable"), 0755); err != nil {
		t.Fatal(err)
	}
	p, err := s.Plan()
	if err != nil {
		t.Fatal(err)
	}
	p.Replaces = false
	err = s.Apply(p, CLIInstallConsent{Install: true})
	if err == nil || err.Error() != "existing executable differs; explicit replacement confirmation is required" {
		t.Errorf("replacement consent bypass: %v", err)
	}
	assertInstallFile(t, target, "previous executable")
	assertInstallFile(t, s.Executable, "original executable")
	if _, err := os.Stat(filepath.Join(s.Home, ".profile")); !os.IsNotExist(err) {
		t.Fatalf("rejection changed profile: %v", err)
	}
	if _, err := os.Stat(filepath.Join(s.Home, ".zesh-skills-install.lock")); !os.IsNotExist(err) {
		t.Fatalf("rejection retained lock: %v", err)
	}
}

func TestCLIInstallPreviewCannotSuppressWindowsPATH(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows filesystem paths require Windows")
	}
	s := installFixture(t, "windows", "")
	store := s.WindowsPath.(*installTestPath)
	p, err := s.Plan()
	if err != nil {
		t.Fatal(err)
	}
	destination := p.Destination
	p.PATHChanged = false
	p.Destination = filepath.Join(s.Home, "redirected.exe")
	if err := s.Apply(p, CLIInstallConsent{Install: true}); err != nil {
		t.Fatal(err)
	}
	if store.value != `C:\Existing;`+filepath.Dir(destination) || store.valueType != 2 || store.writes != 1 {
		t.Fatalf("preview suppressed PATH persistence: %+v", store)
	}
	assertInstallFile(t, destination, "original executable")
	assertInstallFile(t, s.Executable, "original executable")
	if _, err := os.Stat(p.Destination); !os.IsNotExist(err) {
		t.Fatalf("preview redirected write: %v", err)
	}
}
