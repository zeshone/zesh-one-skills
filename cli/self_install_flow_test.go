package main

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var errHostUnavailable = errors.New("cannot determine the current user")

// stubHostInstall replaces host detection for the duration of one test, so the
// public CLI flow never reads this machine's home, environment or user PATH.
func stubHostInstall(t *testing.T, service CLIInstallService, err error) {
	t.Helper()
	previous := hostCLIInstallService
	hostCLIInstallService = func() (CLIInstallService, error) { return service, err }
	t.Cleanup(func() { hostCLIInstallService = previous })
}

func installFlow(t *testing.T, args []string, input string, terminal bool) (int, string, string) {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	source := Source{Client: &http.Client{Transport: failingInstallTransport{t}}, API: "http://127.0.0.1:1", Raw: "http://127.0.0.1:1"}
	code, out, stderr := invoke(args, input, source, terminal)
	if entries, err := os.ReadDir(root); err != nil || len(entries) != 0 {
		t.Fatalf("install flow wrote into the project root: %v %v", entries, err)
	}
	return code, out, stderr
}

type failingInstallTransport struct{ t *testing.T }

func (c failingInstallTransport) RoundTrip(*http.Request) (*http.Response, error) {
	c.t.Fatal("--install-cli must not use the network")
	return nil, nil
}

func installPreviewEN(service CLIInstallService, plan CLIInstallPlan) string {
	preview := "\nInstall this executable for the current user\nSource: " + service.Executable + "\nDestination: " + plan.Destination + "\n"
	if len(plan.ConfigFiles) > 0 {
		preview += "Shell startup files: " + strings.Join(plan.ConfigFiles, ", ") + "\n"
	}
	if plan.PATHChanged {
		preview += "The user PATH will be updated; new terminals only.\n"
	} else {
		preview += "The user PATH already contains the destination.\n"
	}
	if plan.Replaces {
		preview += "An existing different executable will be replaced: " + plan.Destination + "\n"
	}
	return preview + "Skills stay project-local; nothing is installed globally.\n"
}

func TestCLIInstallFlagPreviewsAndInstallsOffline(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	stubHostInstall(t, service, nil)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli", "--yes"}, "", false)
	want := installPreviewEN(service, plan) + "Installed. Open a new terminal to use zesh-skills.\n"
	if code != 0 || out != want || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q\nwant stdout=%q", code, out, stderr, want)
	}
	assertInstallFile(t, service.Executable, "original executable")
	assertInstallFile(t, plan.Destination, "original executable")
	profile, err := os.ReadFile(filepath.Join(service.Home, ".profile"))
	if err != nil || !strings.Contains(string(profile), "# >>> zesh-skills PATH >>>") || !strings.Contains(string(profile), ".local/bin") {
		t.Fatalf("profile %q, error %v", profile, err)
	}
}

func TestCLIInstallFlagDeclineChangesNothing(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	stubHostInstall(t, service, nil)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli"}, "n\n", true)
	want := installPreviewEN(service, plan) + "Proceed? (y/N) [n]: Cancelled; no files changed.\n"
	if code != 0 || out != want || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	for _, path := range []string{plan.Destination, filepath.Join(service.Home, ".profile"), filepath.Join(service.Home, ".local")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("decline wrote %s: %v", path, err)
		}
	}
	assertInstallFile(t, service.Executable, "original executable")
}

func TestCLIInstallFlagRequiresExplicitConfirmation(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	stubHostInstall(t, service, nil)
	code, out, stderr := installFlow(t, []string{"--install-cli", "--non-interactive"}, "", false)
	if code != 2 || out != "" || stderr != "zesh-skills: --install-cli requires --yes to confirm\nUse --help for usage.\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if _, err := os.Stat(filepath.Join(service.Home, ".local")); !os.IsNotExist(err) {
		t.Fatalf("usage failure wrote files: %v", err)
	}
}

func TestCLIInstallFlagRejectsSilentReplacement(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	target := filepath.Join(service.Home, ".local", "bin", "zesh-skills")
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("previous executable"), 0755); err != nil {
		t.Fatal(err)
	}
	stubHostInstall(t, service, nil)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli", "--yes"}, "", false)
	want := installPreviewEN(service, plan)
	if code != 1 || out != want || stderr != "zesh-skills: existing executable differs; remove it or run interactively to confirm replacement\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	assertInstallFile(t, target, "previous executable")
	if _, err := os.Stat(filepath.Join(service.Home, ".profile")); !os.IsNotExist(err) {
		t.Fatalf("rejection wrote profile: %v", err)
	}
}

func TestCLIInstallFlagReplacesAfterInteractiveConsent(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	target := filepath.Join(service.Home, ".local", "bin", "zesh-skills")
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("previous executable"), 0755); err != nil {
		t.Fatal(err)
	}
	stubHostInstall(t, service, nil)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli"}, "y\ny\n", true)
	want := installPreviewEN(service, plan) +
		"Proceed? (y/N) [n]: Replace the existing executable at " + plan.Destination + "? (y/N) [n]: " +
		"Installed. Open a new terminal to use zesh-skills.\n"
	if code != 0 || out != want || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q\nwant=%q", code, out, stderr, want)
	}
	assertInstallFile(t, target, "original executable")
	assertInstallFile(t, service.Executable, "original executable")
}

func TestCLIInstallFlagSpanish(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	stubHostInstall(t, service, nil)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli", "--yes", "--es"}, "", false)
	want := "\nInstalar este ejecutable para el usuario actual\nOrigen: " + service.Executable + "\nDestino: " + plan.Destination + "\n" +
		"Archivos de inicio del shell: " + strings.Join(plan.ConfigFiles, ", ") + "\n" +
		"Se actualizará el PATH del usuario; solo aplica a terminales nuevas.\n" +
		"Las skills siguen siendo locales del proyecto; no se instala nada globalmente.\n" +
		"Instalado. Abre una terminal nueva para usar zesh-skills.\n"
	if code != 0 || out != want || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q\nwant=%q", code, out, stderr, want)
	}
	assertInstallFile(t, plan.Destination, "original executable")
}

func TestCLIInstallFlagRespectsHelpAndVersionPrecedence(t *testing.T) {
	service := installFixture(t, "linux", "/bin/sh")
	stubHostInstall(t, service, nil)
	hostCLIInstallService = func() (CLIInstallService, error) {
		t.Fatal("help/version must not detect the host")
		return CLIInstallService{}, nil
	}
	for _, tt := range []struct {
		args []string
		want string
	}{
		{[]string{"--install-cli", "--help"}, helpText},
		{[]string{"--install-cli", "--version"}, "zesh-skills dev\n"},
	} {
		code, out, stderr := installFlow(t, tt.args, "", false)
		if code != 0 || out != tt.want || stderr != "" {
			t.Fatalf("%v: exit=%d stdout=%q stderr=%q", tt.args, code, out, stderr)
		}
	}
}

func TestCLIInstallFlagReportsHostFailure(t *testing.T) {
	stubHostInstall(t, CLIInstallService{}, errHostUnavailable)
	code, out, stderr := installFlow(t, []string{"--install-cli", "--yes"}, "", false)
	if code != 1 || out != "" || stderr != "zesh-skills: "+errHostUnavailable.Error()+"\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
}

func TestCLIInstallFlagWindowsFlow(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows user PATH flow requires a Windows host")
	}
	service := installFixture(t, "windows", "")
	store := service.WindowsPath.(*installTestPath)
	stubHostInstall(t, service, nil)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli", "--yes"}, "", false)
	want := installPreviewEN(service, plan) + "Installed. Open a new terminal to use zesh-skills.\n"
	if code != 0 || out != want || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q\nwant=%q", code, out, stderr, want)
	}
	if store.value != `C:\Existing;`+filepath.Dir(plan.Destination) || store.valueType != 2 || store.writes != 1 {
		t.Fatalf("user PATH %+v", store)
	}
	assertInstallFile(t, plan.Destination, "original executable")
	assertInstallFile(t, service.Executable, "original executable")
	again, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr = installFlow(t, []string{"--install-cli", "--yes"}, "", false)
	if want := installPreviewEN(service, again) + "Installed. Open a new terminal to use zesh-skills.\n"; code != 0 || out != want || stderr != "" {
		t.Fatalf("second run: exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if store.writes != 1 || strings.Contains(out, "will be updated") {
		t.Fatalf("second run changed PATH: %+v\n%s", store, out)
	}
}

func TestCLIInstallFlagWindowsRejectsSilentReplacement(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows user PATH flow requires a Windows host")
	}
	service := installFixture(t, "windows", "")
	store := service.WindowsPath.(*installTestPath)
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(plan.Destination), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(plan.Destination, []byte("previous executable"), 0755); err != nil {
		t.Fatal(err)
	}
	stubHostInstall(t, service, nil)
	plan, err = service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	code, out, stderr := installFlow(t, []string{"--install-cli", "--yes"}, "", false)
	if code != 1 || out != installPreviewEN(service, plan) || stderr != "zesh-skills: existing executable differs; remove it or run interactively to confirm replacement\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	assertInstallFile(t, plan.Destination, "previous executable")
	if store.writes != 0 || store.value != `C:\Existing` {
		t.Fatalf("rejection changed PATH: %+v", store)
	}
}

func TestCLIInstallFlagWindowsInteractiveReplacement(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows user PATH flow requires a Windows host")
	}
	for _, tt := range []struct {
		name      string
		input     string
		wantCode  int
		wantCfg   string
		wantPathW int
		cancelled bool
	}{
		{"accepted", "y\ny\n", 0, "original executable", 1, false},
		{"replacement refused", "y\nn\n", 0, "previous executable", 0, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			service := installFixture(t, "windows", "")
			store := service.WindowsPath.(*installTestPath)
			plan, err := service.Plan()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(plan.Destination), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(plan.Destination, []byte("previous executable"), 0755); err != nil {
				t.Fatal(err)
			}
			plan, err = service.Plan()
			if err != nil {
				t.Fatal(err)
			}
			stubHostInstall(t, service, nil)
			code, out, stderr := installFlow(t, []string{"--install-cli"}, tt.input, true)
			want := installPreviewEN(service, plan) +
				"Proceed? (y/N) [n]: Replace the existing executable at " + plan.Destination + "? (y/N) [n]: "
			if tt.cancelled {
				want += "Cancelled; no files changed.\n"
			} else {
				want += "Installed. Open a new terminal to use zesh-skills.\n"
			}
			if code != tt.wantCode || out != want || stderr != "" {
				t.Fatalf("exit=%d stdout=%q stderr=%q\nwant=%q", code, out, stderr, want)
			}
			assertInstallFile(t, plan.Destination, tt.wantCfg)
			assertInstallFile(t, service.Executable, "original executable")
			if store.writes != tt.wantPathW {
				t.Fatalf("PATH writes=%d, want %d (%+v)", store.writes, tt.wantPathW, store)
			}
			if tt.cancelled && store.value != `C:\Existing` {
				t.Fatalf("refused replacement changed PATH: %+v", store)
			}
		})
	}
}
