package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testSource(t *testing.T) Source {
	t.Helper()
	cat, files := specimen()
	return testCatalogSource(t, cat, files)
}

func testCatalogSource(t *testing.T, cat Catalog, files map[string][]byte) Source {
	t.Helper()
	commit := strings.Repeat("a", 40)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/commits/main":
			fmt.Fprintf(w, `{"sha":%q}`, commit)
		case "/" + commit + "/skills-index.json":
			json.NewEncoder(w).Encode(cat)
		default:
			data, ok := files[strings.TrimPrefix(r.URL.Path, "/"+commit+"/")]
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Write(data)
		}
	}))
	t.Cleanup(server.Close)
	return Source{Client: server.Client(), API: server.URL, Raw: server.URL}
}
func invoke(args []string, input string, source Source, tty bool) (int, string, string) {
	var out, err bytes.Buffer
	code := Run(args, strings.NewReader(input), &out, &err, source, tty)
	return code, out.String(), err.String()
}
func TestCLIHelpAndUsageAreSideEffectFree(t *testing.T) {
	for _, tt := range []struct {
		args           []string
		code           int
		stdout, stderr string
	}{
		{[]string{"--help"}, 0, helpText, ""}, {[]string{"--version"}, 0, "zesh-skills dev\n", ""},
		{[]string{"--unknown"}, 2, "", "zesh-skills: flag provided but not defined: -unknown\nUse --help for usage.\n"},
		{[]string{"unexpected"}, 2, "", "zesh-skills: unexpected positional arguments\nUse --help for usage.\n"},
	} {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			code, out, err := invoke(tt.args, "", Source{}, false)
			if code != tt.code || out != tt.stdout || err != tt.stderr {
				t.Fatalf("got %d %q %q", code, out, err)
			}
		})
	}
	code, out, err := invoke(nil, "", Source{}, false)
	if code != 2 || out != "" || !strings.Contains(err, "non-interactive mode requires") {
		t.Fatalf("%d %q %q", code, out, err)
	}
}
func TestCLIListsEveryRelevantTopicWithoutWriting(t *testing.T) {
	root := t.TempDir()
	code, out, err := invoke([]string{"--root", root, "--scope", "frontend", "--list"}, "", testSource(t), false)
	want := "\ndesign\n  [ ] 1. frontend/design — UI layout\n\nsecurity\n  [ ] 2. frontend/security — Security checks\n"
	if code != 0 || err != "" || out != want {
		t.Fatalf("%d %q %q", code, out, err)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatal("listing mutated project")
	}
}
func TestCLILandingMockupInstallsDesignOnly(t *testing.T) {
	root := t.TempDir()
	source := testSource(t)
	// Root, scope, architecture, agents; Space/Enter choose design in the TUI; extra file and confirmation.
	input := "\nfrontend\nlanding mockup\nclaude,codex,opencode,pi\n \n\ny\n"
	code, out, err := invoke([]string{"--root", root}, input, source, true)
	if code != 0 || err != "" || !strings.Contains(out, "Architecture") || !strings.Contains(out, "Security checks") || !strings.Contains(out, "Installed 1 selected skill") {
		t.Fatalf("%d %q %q", code, out, err)
	}
	if !strings.Contains(readInstalled(t, root, "AGENTS.md"), "landing mockup") {
		t.Fatal("architecture missing")
	}
	if _, err := os.Stat(filepath.Join(root, ".pi/skills/frontend-security")); !os.IsNotExist(err) {
		t.Fatal("unselected security installed")
	}
}
func TestCLIRejectAndCancelLeaveFilesUnchanged(t *testing.T) {
	root := t.TempDir()
	source := testSource(t)
	args := []string{"--root", root, "--skills", "frontend/design", "--agents", "pi", "--architecture", "mockup"}
	code, out, err := invoke(args, "", Source{}, false)
	if code != 2 || out != "" || !strings.Contains(err, "--yes") {
		t.Fatalf("%d %q %q", code, out, err)
	}
	args = append(args, "--yes")
	code, _, err = invoke(args, "", source, false)
	if code != 0 || err != "" {
		t.Fatal(code, err)
	}
	data := readInstalled(t, root, "AGENTS.md")
	args = append(args, "--scope", "wrong")
	code, _, err = invoke(args, "", Source{}, false)
	if code != 2 || !strings.Contains(err, "unknown scope") {
		t.Fatal(code, err)
	}
	if readInstalled(t, root, "AGENTS.md") != data {
		t.Fatal("rejection changed file")
	}
	empty := t.TempDir()
	code, out, err = invoke([]string{"--root", empty}, "\nfrontend\nmockup\npi\nq", source, true)
	if code != 0 || err != "" || !strings.Contains(out, "Cancelled") {
		t.Fatalf("%d %q %q", code, out, err)
	}
	entries, _ := os.ReadDir(empty)
	if len(entries) != 0 {
		t.Fatal("cancel wrote data")
	}
}
func TestDetectProjectWithoutExecutingTools(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"react":"19","next":"15","@capacitor/core":"7"},"scripts":{"postinstall":"do-not-run"}}`), 0644)
	os.WriteFile(filepath.Join(root, "demo.csproj"), []byte("<Project/>"), 0644)
	os.WriteFile(filepath.Join(root, "claude.md"), []byte("mine"), 0644)
	os.WriteFile(filepath.Join(root, "GEMINI.md"), []byte("mine"), 0644)
	detected, err := Detect(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(detected.Scopes, ",") != "app,backend,frontend" || !strings.Contains(strings.Join(detected.Technologies, ","), "react") || len(detected.Instructions) != 2 {
		t.Fatalf("%+v", detected)
	}
}
