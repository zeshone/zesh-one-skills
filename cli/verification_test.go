package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func verifyArgs(root string) []string {
	return []string{"--root", root, "--non-interactive", "--agents", "pi", "--architecture", "mockup", "--skills", "frontend/design", "--yes"}
}
func verifyPreview(root string) string {
	return fmt.Sprintf("\nInstall 1 selected skill(s) into %s\nSource: main @ %s\nArchitecture: mockup\nAgents: pi\nSkills: frontend/design\nContent: .zesh-skills/content/\nInstruction files: \n  .pi/skills/<scope>-<name>/SKILL.md\n  AGENTS.md (managed index and triggers)\n", root, strings.Repeat("a", 40))
}
func TestVerifyCaseOnlyIndexConflict(t *testing.T) {
	root := t.TempDir()
	source := testSource(t)
	args := verifyArgs(root)
	code, out, stderr := invoke(args, "", source, false)
	success := verifyPreview(root) + "Installed 1 selected skill(s). Existing selections retained.\n"
	if code != 0 || out != success || stderr != "" {
		t.Fatalf("setup: %d %q %q", code, out, stderr)
	}
	original := readInstalled(t, root, "AGENTS.md")
	modified := strings.Replace(original, "UI layout", "UI LAYOUT", 1)
	if modified == original {
		t.Fatal("fixture did not edit")
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(modified), 0644); err != nil {
		t.Fatal(err)
	}
	before := readInstalled(t, root, manifestPath)
	code, out, stderr = invoke(args, "", source, false)
	wantErr := "zesh-skills: cancelled: local changes preserved in AGENTS.md\n"
	if code != 1 || out != verifyPreview(root) || stderr != wantErr {
		t.Errorf("expected rejection exit1 with exact preview/error; observed exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if readInstalled(t, root, "AGENTS.md") != modified {
		t.Error("case-only user index edit overwritten without confirmation")
	}
	if readInstalled(t, root, manifestPath) != before {
		t.Error("manifest changed")
	}
	if _, err := os.Stat(filepath.Join(root, ".zesh-skills.lock")); !os.IsNotExist(err) {
		t.Error("lock retained")
	}
}
func TestVerifyConcurrentLockRejectsSecondRun(t *testing.T) {
	root := t.TempDir()
	cat, files := specimen()
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	first := true
	in := Installer{Fetch: func(_ string, name string) ([]byte, error) {
		if first {
			first = false
			close(entered)
			<-release
		}
		return files[name], nil
	}}
	go func() {
		done <- in.Install(root, cat, strings.Repeat("a", 40), []string{"frontend/design"}, []string{"pi"}, "mockup", nil)
	}()
	<-entered
	code, out, stderr := invoke(verifyArgs(root), "", testSource(t), false)
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Error(err)
	}
	if code != 1 || out != verifyPreview(root) || stderr != "zesh-skills: project is locked or not writable; inspect .zesh-skills.lock\n" {
		t.Errorf("second run: %d %q %q", code, out, stderr)
	}
	if len(entries) != 1 || entries[0].Name() != ".zesh-skills.lock" {
		t.Errorf("rejected run side effects: %v", entries)
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".zesh-skills.lock")); !os.IsNotExist(err) {
		t.Fatal("lock not released")
	}
}

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

func TestVerifyInstructionEditsDuringConfirmationArePreserved(t *testing.T) {
	root := t.TempDir()
	source := testSource(t)
	if code, _, stderr := invoke(verifyArgs(root), "", source, false); code != 0 || stderr != "" {
		t.Fatal(code, stderr)
	}
	name := ".zesh-skills/content/skills/frontend/design/SKILL.md"
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), []byte("local skill edit"), 0644); err != nil {
		t.Fatal(err)
	}
	before := readInstalled(t, root, manifestPath)
	edited := readInstalled(t, root, "AGENTS.md") + "\nNew unrelated user instruction.\n"
	var stdout, stderr strings.Builder
	step, preview := 0, ""
	input := readerFunc(func(p []byte) (int, error) {
		step++
		line := "\n"
		if step == 3 {
			preview = stdout.String()
			if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(edited), 0644); err != nil {
				t.Fatal(err)
			}
			line = "y\n"
		}
		if step > 3 {
			return 0, fmt.Errorf("unexpected prompt")
		}
		return copy(p, line), nil
	})
	args := []string{"--root", root, "--scope", "frontend", "--agents", "pi", "--architecture", "mockup", "--skills", "frontend/design", "--yes"}
	code := Run(args, input, &stdout, &stderr, source, true)
	if code != 1 || preview == "" || stdout.String() != preview || stderr.String() != "zesh-skills: instruction changed during installation: AGENTS.md; retry\n" {
		t.Fatalf("%d %q %q", code, stdout.String(), stderr.String())
	}
	if readInstalled(t, root, "AGENTS.md") != edited || readInstalled(t, root, manifestPath) != before || readInstalled(t, root, name) != "local skill edit" {
		t.Fatal("rejected installation changed local data")
	}
	if _, err := os.Stat(filepath.Join(root, ".zesh-skills.lock")); !os.IsNotExist(err) {
		t.Fatal("lock retained")
	}
}

func TestVerifyChecksumRejectsWithoutWrites(t *testing.T) {
	root := t.TempDir()
	cat, files := specimen()
	commit := strings.Repeat("a", 40)
	files[cat.Skills[0].Files[0].Path] = []byte("tampered content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/commits/main":
			fmt.Fprintf(w, "{\"sha\":%q}", commit)
		case "/" + commit + "/skills-index.json":
			json.NewEncoder(w).Encode(cat)
		default:
			w.Write(files[strings.TrimPrefix(r.URL.Path, "/"+commit+"/")])
		}
	}))
	defer server.Close()
	code, out, stderr := invoke(verifyArgs(root), "", Source{Client: server.Client(), API: server.URL, Raw: server.URL}, false)
	if code != 1 || out != verifyPreview(root) || stderr != "zesh-skills: checksum mismatch: skills/frontend/design/SKILL.md\n" {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("rejected checksum side effects %v %v", entries, err)
	}
}
