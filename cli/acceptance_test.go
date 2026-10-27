package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIUsesInvocationDirectoryByDefault(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	code, out, err := invoke([]string{"--agents", "pi", "--architecture", "landing mockup", "--skills", "frontend/design", "--yes"}, "", testSource(t), false)
	if code != 0 || err != "" || !strings.Contains(out, root) {
		t.Fatalf("%d %q %q", code, out, err)
	}
	if !strings.Contains(readInstalled(t, root, "AGENTS.md"), "frontend/design") {
		t.Fatal("default root missing index")
	}
}
func TestRepositoryCatalogAndAllSelectedResources(t *testing.T) {
	data, err := os.ReadFile("../skills-index.json")
	if err != nil {
		t.Fatal(err)
	}
	var c Catalog
	if err = json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	if err = c.Validate(); err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, s := range c.Skills {
		ids = append(ids, s.ID)
	}
	in := Installer{Fetch: func(_ string, name string) ([]byte, error) {
		return os.ReadFile(filepath.Join("..", filepath.FromSlash(name)))
	}}
	root := t.TempDir()
	if err = in.Install(root, c, strings.Repeat("a", 40), ids, []string{"claude", "codex", "opencode", "pi"}, "acceptance trial", nil); err != nil {
		t.Fatal(err)
	}
	index := readInstalled(t, root, "AGENTS.md")
	for _, s := range c.Skills {
		if !strings.Contains(index, s.ID) {
			t.Fatal("missing index entry", s.ID)
		}
		for _, f := range s.Files {
			if _, err = os.Stat(filepath.Join(root, ".zesh-skills/content", filepath.FromSlash(f.Path))); err != nil {
				t.Fatal(err)
			}
		}
	}
}
func TestExtraInstructionIndexIsRetainedAndConflictProtected(t *testing.T) {
	root := t.TempDir()
	cat, files := specimen()
	in := installerFor(files)
	in.Instructions = []string{"GEMINI.md"}
	commit := strings.Repeat("a", 40)
	if err := in.Install(root, cat, commit, []string{"frontend/design"}, []string{"pi"}, "first", nil); err != nil {
		t.Fatal(err)
	}
	in.Instructions = nil
	if err := in.Install(root, cat, commit, []string{"frontend/design"}, []string{"pi"}, "second", nil); err != nil {
		t.Fatal(err)
	}
	original := readInstalled(t, root, "GEMINI.md")
	if !strings.Contains(original, "Architecture: second") {
		t.Fatal("previously configured index was not updated")
	}
	modified := strings.Replace(original, "UI layout", "my edits", 1)
	os.WriteFile(filepath.Join(root, "GEMINI.md"), []byte(modified), 0644)
	if err := in.Install(root, cat, commit, []string{"frontend/design"}, []string{"pi"}, "third", func(string) bool { return false }); err == nil {
		t.Fatal("modified managed index overwritten")
	}
	if readInstalled(t, root, "GEMINI.md") != modified {
		t.Fatal("rejection changed user content")
	}
}
