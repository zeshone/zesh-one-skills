package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func installFixture(t *testing.T, platform, shell string) CLIInstallService {
	t.Helper()
	home := t.TempDir()
	source := filepath.Join(t.TempDir(), "original")
	if err := os.WriteFile(source, []byte("original executable"), 0755); err != nil {
		t.Fatal(err)
	}
	return CLIInstallService{OS: platform, Home: home, Executable: source, Env: map[string]string{"SHELL": shell, "LOCALAPPDATA": home}, WindowsPath: &installTestPath{value: `C:\Existing`, valueType: 2}}
}

func assertInstallFile(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil || string(got) != want {
		t.Fatalf("%s: got %q (%v), want %q", path, got, err, want)
	}
}

func TestCLIInstallSupportedShells(t *testing.T) {
	for _, platform := range []string{"linux", "darwin"} {
		for shell, configs := range map[string][]string{
			"/bin/sh": {".profile"}, "/bin/bash": {".bashrc", ".profile"},
			"/bin/zsh": {".zshrc", ".zprofile"}, "/usr/bin/fish": {".config/fish/conf.d/zesh-skills.fish"},
		} {
			t.Run(platform+shell, func(t *testing.T) {
				s := installFixture(t, platform, shell)
				for _, config := range configs {
					p := filepath.Join(s.Home, filepath.FromSlash(config))
					if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(p, []byte("# existing setting\n"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				plan, err := s.Plan()
				if err != nil {
					t.Fatal(err)
				}
				if err := s.Apply(plan, CLIInstallConsent{Install: true}); err != nil {
					t.Fatal(err)
				}
				assertInstallFile(t, s.Executable, "original executable")
				assertInstallFile(t, plan.Destination, "original executable")
				for _, config := range configs {
					p := filepath.Join(s.Home, filepath.FromSlash(config))
					before, err := os.ReadFile(p)
					if err != nil || !strings.HasPrefix(string(before), "# existing setting\n") || !strings.Contains(string(before), "zesh-skills PATH") {
						t.Fatalf("config %q: %s %v", config, before, err)
					}
					next, err := s.Plan()
					if err != nil {
						t.Fatal(err)
					}
					if err := s.Apply(next, CLIInstallConsent{Install: true}); err != nil {
						t.Fatal(err)
					}
					assertInstallFile(t, p, string(before))
				}
			})
		}
	}
}

func TestCLIInstallRejectsUnsafePlans(t *testing.T) {
	for _, scenario := range []string{"unsupported OS", "unsupported shell", "delimiter", "relative home", "outside config", "missing source", "directory config", "symlink"} {
		t.Run(scenario, func(t *testing.T) {
			s := installFixture(t, "linux", "/bin/sh")
			original := s.Executable
			originalHome := s.Home
			switch scenario {
			case "unsupported OS":
				s.OS = "plan9"
			case "unsupported shell":
				s.Env["SHELL"] = "/bin/unknown"
			case "delimiter":
				s.Home = filepath.Join(s.Home, "bad:name")
			case "relative home":
				s.Home = "relative"
			case "outside config":
				s.Env["SHELL"] = "/bin/zsh"
				s.Env["ZDOTDIR"] = t.TempDir()
			case "missing source":
				s.Executable += ".missing"
			case "directory config":
				if err := os.Mkdir(filepath.Join(s.Home, ".profile"), 0700); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				if err := os.Symlink(t.TempDir(), filepath.Join(s.Home, ".local")); err != nil {
					t.Skipf("symlink unavailable: %v", err)
				}
			}
			if _, err := s.Plan(); err == nil {
				t.Fatal("unsafe plan accepted")
			}
			assertInstallFile(t, original, "original executable")
			if _, err := os.Stat(filepath.Join(originalHome, ".local", "bin", "zesh-skills")); !os.IsNotExist(err) {
				t.Fatalf("unexpected destination: %v", err)
			}
		})
	}
}

func TestCLIInstallConsentConflictAndDrift(t *testing.T) {
	for _, scenario := range []string{"decline", "conflict", "explicit replacement", "source drift", "target drift", "config drift"} {
		t.Run(scenario, func(t *testing.T) {
			s := installFixture(t, "linux", "/bin/sh")
			destination := filepath.Join(s.Home, ".local", "bin", "zesh-skills")
			if err := os.MkdirAll(filepath.Dir(destination), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(destination, []byte("previous executable"), 0755); err != nil {
				t.Fatal(err)
			}
			config := filepath.Join(s.Home, ".profile")
			if err := os.WriteFile(config, []byte("# keep\n"), 0600); err != nil {
				t.Fatal(err)
			}
			p, err := s.Plan()
			if err != nil {
				t.Fatal(err)
			}
			consent := CLIInstallConsent{Install: true}
			wantSource, wantTarget, wantConfig := "original executable", "previous executable", "# keep\n"
			switch scenario {
			case "decline":
				consent.Install = false
			case "explicit replacement":
				consent.Replace = true
				wantTarget = wantSource
			case "source drift":
				consent.Replace = true
				wantSource = "changed source"
				if err := os.WriteFile(s.Executable, []byte(wantSource), 0755); err != nil {
					t.Fatal(err)
				}
			case "target drift":
				consent.Replace = true
				wantTarget = "changed target"
				if err := os.WriteFile(destination, []byte(wantTarget), 0755); err != nil {
					t.Fatal(err)
				}
			case "config drift":
				consent.Replace = true
				wantConfig = "# changed\n"
				if err := os.WriteFile(config, []byte(wantConfig), 0600); err != nil {
					t.Fatal(err)
				}
			}
			err = s.Apply(p, consent)
			if scenario == "explicit replacement" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil {
				t.Fatal("expected rejection")
			}
			assertInstallFile(t, s.Executable, wantSource)
			assertInstallFile(t, destination, wantTarget)
			if scenario != "explicit replacement" {
				assertInstallFile(t, config, wantConfig)
			}
		})
	}
}

type failingInstallPath struct {
	installTestPath
	fail bool
}

func (p *failingInstallPath) Write(value string, kind uint32) error {
	if p.fail {
		p.fail = false
		return errors.New("injected PATH write failure")
	}
	return p.installTestPath.Write(value, kind)
}

func TestCLIInstallWindowsRollbackAndDrift(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows paths require Windows")
	}
	for _, scenario := range []string{"rollback", "drift", "idempotent", "invalid type"} {
		t.Run(scenario, func(t *testing.T) {
			s := installFixture(t, "windows", "")
			store := &failingInstallPath{installTestPath: installTestPath{value: `C:\Existing`, valueType: 2}}
			s.WindowsPath = store
			if scenario == "invalid type" {
				store.valueType = 7
				if _, err := s.Plan(); err == nil {
					t.Fatal("accepted non-string PATH")
				}
				if store.writes != 0 {
					t.Fatal("unexpected writes")
				}
				return
			}
			p, err := s.Plan()
			if err != nil {
				t.Fatal(err)
			}
			switch scenario {
			case "rollback":
				store.fail = true
			case "drift":
				store.value = `C:\Changed`
			}
			err = s.Apply(p, CLIInstallConsent{Install: true})
			if scenario == "idempotent" {
				if err != nil {
					t.Fatal(err)
				}
				before, writes := store.value, store.writes
				p, err = s.Plan()
				if err != nil {
					t.Fatal(err)
				}
				if err := s.Apply(p, CLIInstallConsent{Install: true}); err != nil {
					t.Fatal(err)
				}
				if store.value != before || store.writes != writes {
					t.Fatalf("duplicate PATH write: %+v", store)
				}
				assertInstallFile(t, p.Destination, "original executable")
			} else {
				if err == nil {
					t.Fatal("expected failure")
				}
				if _, err := os.Stat(p.Destination); !os.IsNotExist(err) {
					t.Fatalf("destination not rolled back: %v", err)
				}
				want := `C:\Existing`
				if scenario == "drift" {
					want = `C:\Changed`
				}
				if store.value != want || store.writes != 0 {
					t.Fatalf("unexpected PATH mutation: %+v", store)
				}
			}
			assertInstallFile(t, s.Executable, "original executable")
		})
	}
}
