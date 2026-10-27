package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

// Exercise remote JSON through Source, including types Go's Skill cannot represent.
func translatedSource(t *testing.T, value any, fetches ...*atomic.Int32) Source {
	t.Helper()
	cat, files := specimen()
	data, err := json.Marshal(cat)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	raw["skills"].([]any)[0].(map[string]any)["descriptionEs"] = value
	commit := strings.Repeat("a", 40)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/commits/main":
			fmt.Fprintf(w, `{"sha":%q}`, commit)
		case "/" + commit + "/skills-index.json":
			json.NewEncoder(w).Encode(raw)
		default:
			for _, count := range fetches {
				count.Add(1)
			}
			w.Write(files[strings.TrimPrefix(r.URL.Path, "/"+commit+"/")])
		}
	}))
	t.Cleanup(server.Close)
	return Source{Client: server.Client(), API: server.URL, Raw: server.URL}
}

func TestCatalogDescriptionsOnly(t *testing.T) {
	for _, flag := range []string{"--es", "--es=false", ""} {
		t.Run(flag, func(t *testing.T) {
			root := t.TempDir()
			args := []string{"--root", root, "--scope", "frontend", "--list"}
			if flag != "" {
				args = append(args, flag)
			}
			description := "UI layout"
			if flag == "--es" {
				description = "Diseño de interfaces"
			}
			code, out, stderr := invoke(args, "", translatedSource(t, "Diseño de interfaces"), false)
			want := "\ndesign\n  [ ] 1. frontend/design — " + description + "\n\nsecurity\n  [ ] 2. frontend/security — Security checks\n"
			if code != 0 || stderr != "" || out != want {
				t.Fatalf("%d %q %q; want %q", code, out, stderr, want)
			}
			args = []string{"--root", root, "--scope", "frontend", "--architecture", "clean", "--agents", "pi"}
			if flag != "" {
				args = append(args, flag)
			}
			code, out, stderr = invoke(args, "\nq", translatedSource(t, "Diseño de interfaces"), true)
			if code != 0 || stderr != "" || !strings.Contains(out, "frontend/design [design]") || !strings.Contains(out, description) {
				t.Fatalf("picker: %d %q %q", code, out, stderr)
			}
			if entries, err := os.ReadDir(root); err != nil || len(entries) != 0 {
				t.Fatal("presentation wrote project files", err)
			}
		})
	}
}

func TestCatalogOptionalSpanishFallbacks(t *testing.T) {
	for _, value := range []string{"", strings.Repeat("x", 1024), strings.Repeat("á", 512)} {
		root := t.TempDir()
		code, out, stderr := invoke([]string{"--es", "--root", root, "--scope", "frontend", "--list"}, "", translatedSource(t, value), false)
		want := value
		if value == "" {
			want = "UI layout" // Unknown ID/source stays discoverable in English.
		}
		if code != 0 || stderr != "" || !strings.Contains(out, "frontend/design — "+brief(want)) {
			t.Fatalf("optional/boundary value: %d %q %q", code, out, stderr)
		}
	}
	catalog, files := specimen()
	catalog.Skills[0].Description = "ZeshOne Tailwind CSS 4 conventions. Trigger: When styling UI with Tailwind utilities, semantic tokens, responsive rules, or runtime dynamic values."
	code, out, stderr := invoke([]string{"--es", "--root", t.TempDir(), "--scope", "frontend", "--list"}, "", testCatalogSource(t, catalog, files), false)
	if code != 0 || stderr != "" || !strings.Contains(out, brief(catalog.Skills[0].Description)) {
		t.Fatal("matching source with a different ID must not borrow a translation", code, out, stderr)
	}
}

func TestCatalogRejectsInvalidSpanishBeforeWrites(t *testing.T) {
	for _, value := range []any{"bad\ntext", strings.Repeat("x", 1025), "á" + strings.Repeat("x", 1023), 42, []string{"texto"}, nil, true} {
		for _, spanish := range []bool{false, true} {
			root := t.TempDir()
			args := []string{"--root", root, "--scope", "frontend", "--list"}
			message := "incomplete skill: frontend/design"
			if _, ok := value.(string); !ok {
				message = "invalid catalog JSON"
			}
			if spanish {
				args = append(args, "--es")
				if message == "invalid catalog JSON" {
					message = "el JSON del catálogo no es válido"
				} else {
					message = "skill incompleta: frontend/design"
				}
			}
			code, out, stderr := invoke(args, "", translatedSource(t, value), false)
			if code != 1 || out != "" || stderr != "zesh-skills: "+message+"\n" {
				t.Fatalf("%T Spanish=%v: %d %q %q", value, spanish, code, out, stderr)
			}
			if entries, err := os.ReadDir(root); err != nil || len(entries) != 0 {
				t.Fatal("invalid metadata wrote files", err)
			}
			if code, out, stderr := invoke(verifyArgs(root), "", testSource(t), false); code != 0 || stderr != "" {
				t.Fatal(code, out, stderr)
			}
			before := installedSnapshot(t, root)
			installArgs := verifyArgs(root)
			if spanish {
				installArgs = append(installArgs, "--es")
			}
			var fetches atomic.Int32
			code, out, stderr = invoke(installArgs, "", translatedSource(t, value, &fetches), false)
			if code != 1 || out != "" || stderr != "zesh-skills: "+message+"\n" {
				t.Fatalf("install: %d %q %q", code, out, stderr)
			}
			if fetches.Load() != 0 || !reflect.DeepEqual(before, installedSnapshot(t, root)) {
				t.Fatal("invalid metadata fetched content or changed installed state")
			}
		}
	}
}

func TestCatalogLegacySpanishCoverage(t *testing.T) {
	data, err := os.ReadFile("../skills-index.json")
	if err != nil {
		t.Fatal(err)
	}
	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		t.Fatal(err)
	}
	if len(catalog.Skills) != 35 {
		t.Fatal("update explicit legacy coverage", len(catalog.Skills))
	}
	for _, skill := range catalog.Skills {
		t.Run(skill.ID, func(t *testing.T) {
			original := skill.Description
			skill.DescriptionES = "" // Simulate old main with no optional field.
			root := t.TempDir()
			cat := Catalog{SchemaVersion: 1, Skills: []Skill{skill}}
			args := []string{"--es", "--root", root, "--scope", skill.Scopes[0], "--list"}
			code, out, stderr := invoke(args, "", testCatalogSource(t, cat, nil), false)
			if code != 0 || stderr != "" || strings.Contains(out, brief(original)) {
				t.Fatalf("legacy fallback missing: %d %q %q", code, out, stderr)
			}
			// A changed English source must never inherit a stale translation.
			cat.Skills[0].Description = "Changed English source"
			code, out, stderr = invoke(args, "", testCatalogSource(t, cat, nil), false)
			if code != 0 || stderr != "" || !strings.Contains(out, "Changed English source") {
				t.Fatal(code, out, stderr)
			}
			cat.Skills[0].Description = original
			cat.Skills[0].DescriptionES = "Traducción remota prioritaria"
			code, out, stderr = invoke(args, "", testCatalogSource(t, cat, nil), false)
			if code != 0 || stderr != "" || !strings.Contains(out, "Traducción remota prioritaria") {
				t.Fatal(code, out, stderr)
			}
		})
	}
}

func TestBundledSpanishStructure(t *testing.T) {
	data, err := os.ReadFile("skill_descriptions_es.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries map[string]struct{ Source, ES string }
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatal(err)
	}
	catalogData, err := os.ReadFile("../skills-index.json")
	if err != nil {
		t.Fatal(err)
	}
	var catalog Catalog
	if err := json.Unmarshal(catalogData, &catalog); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 35 {
		t.Fatal("expected 35 bundled translations")
	}
	for _, skill := range catalog.Skills {
		entry, ok := entries[skill.ID]
		if !ok || entry.Source != skill.Description || !plainText(entry.Source, 1024) || !plainText(entry.ES, 1024) || entry.Source == entry.ES {
			t.Errorf("invalid or stale translation: %s", skill.ID)
		}
	}
}

func TestCatalogTranslationDoesNotPersist(t *testing.T) {
	roots := []string{t.TempDir(), t.TempDir(), t.TempDir()}
	for i, root := range roots {
		source := testSource(t)
		if i > 0 {
			source = translatedSource(t, "Diseño de interfaces")
		}
		args := verifyArgs(root)
		if i == 2 {
			args = append(args, "--es")
		}
		if code, out, stderr := invoke(args, "", source, false); code != 0 || stderr != "" {
			t.Fatalf("%d %q %q", code, out, stderr)
		}
	}
	original := installedSnapshot(t, roots[0])
	for _, root := range roots[1:] {
		if got := installedSnapshot(t, root); !reflect.DeepEqual(got, original) {
			t.Fatal("localized metadata changed installed artifact bytes or paths")
		}
		if strings.Contains(readInstalled(t, root, manifestPath), "descriptionEs") {
			t.Fatal("localized-only field persisted")
		}
	}
}

func installedSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		files[relative] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
