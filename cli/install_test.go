package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func specimen() (Catalog, map[string][]byte) {
	body := []byte("---\nname: design\ndescription: UI design\n---\n# Design\nUse [baseline](../../shared/security/baseline.md).\n")
	resource := []byte("# Baseline\n")
	files := map[string][]byte{"skills/frontend/design/SKILL.md": body, "skills/shared/security/baseline.md": resource}
	skill := Skill{ID: "frontend/design", Name: "design", Description: "UI layout", Path: "skills/frontend/design", Version: "1", Scopes: []string{"frontend"}, Topics: []string{"design"}, Triggers: []string{"UI layout"}}
	for _, name := range []string{"skills/frontend/design/SKILL.md", "skills/shared/security/baseline.md"} {
		skill.Files = append(skill.Files, File{Path: name, SHA256: fmt.Sprintf("%x", sha256.Sum256(files[name]))})
	}
	security := Skill{ID: "frontend/security", Name: "security", Description: "Security checks", Path: "skills/frontend/security", Version: "1", Scopes: []string{"frontend"}, Topics: []string{"security"}, Triggers: []string{"security"}, Files: []File{{Path: "skills/frontend/security/SKILL.md", SHA256: strings.Repeat("a", 64)}}}
	return Catalog{SchemaVersion: 1, Skills: []Skill{skill, security}}, files
}
func installerFor(files map[string][]byte) Installer {
	return Installer{Fetch: func(commit, name string) ([]byte, error) {
		data, ok := files[name]
		if !ok {
			return nil, fmt.Errorf("missing file")
		}
		return data, nil
	}}
}
func readInstalled(t *testing.T, root, name string) string {
	t.Helper()
	b, e := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if e != nil {
		t.Fatal(e)
	}
	return string(b)
}
func TestInstallSelectedAndPreserveInstructions(t *testing.T) {
	root := t.TempDir()
	cat, files := specimen()
	existing := "# My project\r\nKeep this.\r\n"
	os.WriteFile(filepath.Join(root, "claude.md"), []byte(existing), 0644)
	in := installerFor(files)
	for i := 0; i < 2; i++ {
		if err := in.Install(root, cat, strings.Repeat("a", 40), []string{"frontend/design"}, []string{"claude", "codex", "opencode", "pi"}, "landing mockup", func(string) bool { t.Fatal("unexpected conflict"); return false }); err != nil {
			t.Fatal(err)
		}
	}
	content := readInstalled(t, root, "claude.md")
	if !strings.HasPrefix(content, existing) || strings.Count(content, "<!-- zesh-skills:start -->") != 1 || !strings.Contains(content, "UI layout") {
		t.Fatalf("instructions changed: %q", content)
	}
	for _, p := range []string{".claude/skills/frontend-design/SKILL.md", ".agents/skills/frontend-design/SKILL.md", ".opencode/skills/frontend-design/SKILL.md", ".pi/skills/frontend-design/SKILL.md"} {
		if !strings.Contains(readInstalled(t, root, p), "name: frontend-design") {
			t.Fatal(p)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".zesh-skills/content/skills/frontend/security")); !os.IsNotExist(err) {
		t.Fatal("unselected security installed")
	}
}
func TestModifiedFileConflictAndLatestContent(t *testing.T) {
	root := t.TempDir()
	cat, files := specimen()
	in := installerFor(files)
	commit := strings.Repeat("a", 40)
	install := func(confirm func(string) bool) error {
		return in.Install(root, cat, commit, []string{"frontend/design"}, []string{"pi"}, "mockup", confirm)
	}
	if err := install(nil); err != nil {
		t.Fatal(err)
	}
	name := ".zesh-skills/content/skills/frontend/design/SKILL.md"
	os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), []byte("my local edits"), 0644)
	before := readInstalled(t, root, ".zesh-skills/manifest.json")
	if err := install(func(string) bool { return false }); err == nil {
		t.Fatal("declined conflict accepted")
	}
	if got := readInstalled(t, root, name); got != "my local edits" {
		t.Fatal(got)
	}
	if got := readInstalled(t, root, ".zesh-skills/manifest.json"); got != before {
		t.Fatal("manifest changed after rejection")
	}
	if err := install(func(string) bool { return true }); err != nil {
		t.Fatal(err)
	}
	files["skills/frontend/design/SKILL.md"] = []byte("new upstream version\n")
	cat.Skills[0].Files[0].SHA256 = fmt.Sprintf("%x", sha256.Sum256(files["skills/frontend/design/SKILL.md"]))
	if err := install(nil); err != nil {
		t.Fatal(err)
	}
	if got := readInstalled(t, root, name); got != "new upstream version\n" {
		t.Fatal(got)
	}
}
func TestRejectedInputsLeaveProjectUntouched(t *testing.T) {
	for _, kind := range []string{"traversal", "bad-hash", "unknown-skill", "bad-schema", "unknown-agent", "malformed-block", "locked"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			cat, files := specimen()
			ids := []string{"frontend/design"}
			agents := []string{"pi"}
			switch kind {
			case "traversal":
				cat.Skills[0].Files[0].Path = "../escape"
			case "bad-hash":
				cat.Skills[0].Files[0].SHA256 = strings.Repeat("0", 64)
			case "unknown-skill":
				ids = []string{"missing"}
			case "bad-schema":
				cat.SchemaVersion = 99
			case "unknown-agent":
				agents = []string{"missing"}
			case "malformed-block":
				os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("mine\n<!-- zesh-skills:start -->\n"), 0644)
			case "locked":
				os.WriteFile(filepath.Join(root, ".zesh-skills.lock"), []byte("other process"), 0600)
			}
			if err := installerFor(files).Install(root, cat, strings.Repeat("a", 40), ids, agents, "mockup", nil); err == nil {
				t.Fatal("accepted rejected input")
			}
			entries, _ := os.ReadDir(root)
			want := 0
			if kind == "malformed-block" || kind == "locked" {
				want = 1
			}
			if len(entries) != want {
				t.Fatalf("side effects: %v", entries)
			}
			if kind == "malformed-block" && readInstalled(t, root, "AGENTS.md") != "mine\n<!-- zesh-skills:start -->\n" {
				t.Fatal("changed instructions")
			}
		})
	}
}
func TestNoWritesThroughSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".zesh-skills")); err != nil {
		t.Skip("symlink privilege unavailable: " + err.Error())
	}
	cat, files := specimen()
	if err := installerFor(files).Install(root, cat, strings.Repeat("a", 40), []string{"frontend/design"}, []string{"pi"}, "mockup", nil); err == nil {
		t.Fatal("followed symlink")
	}
	entries, _ := os.ReadDir(outside)
	if len(entries) != 0 {
		t.Fatal("wrote outside project")
	}
}
