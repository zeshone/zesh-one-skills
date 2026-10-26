package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestArchitectureCatalogExplainsCommonAndHybridProfiles(t *testing.T) {
	profiles := architectureProfiles()
	if len(profiles) < 15 {
		t.Fatal("architecture catalog is too limited")
	}
	seen := map[string]bool{}
	hybrids := 0
	for _, p := range profiles {
		if seen[p.ID] || p.ID == "" || p.Name == "" || p.Description == "" || p.NameES == "" || p.DescriptionES == "" {
			t.Fatalf("invalid profile: %+v", p)
		}
		seen[p.ID] = true
		if p.Hybrid {
			hybrids++
		}
	}
	if hybrids < 4 || profiles[len(profiles)-1].ID != "custom" {
		t.Fatal("hybrid/custom profiles missing")
	}
}
func TestScopePickerSupportsMixedProjects(t *testing.T) {
	m := newOptionPicker(scopeOptions(false), map[string]string{"frontend": "Detected project scope"}, "Choose project scope", false)
	m.apply("down")
	m.apply("space")
	m.apply("enter")
	if !m.accepted || !reflect.DeepEqual(m.ids(), []string{"frontend", "backend"}) {
		t.Fatal(m.ids())
	}
	if !strings.Contains(m.view(100, 24), "Backend") {
		t.Fatal("scope labels missing")
	}
}
func TestArchitecturePickerIsExclusive(t *testing.T) {
	options := architectureOptions(false)
	m := newOptionPicker(options, nil, "Choose architecture", true)
	for m.options[m.cursor].ID != "clean-vertical-slices" {
		m.apply("down")
	}
	m.apply("space")
	m.apply("all")
	m.apply("enter")
	if !m.accepted || !reflect.DeepEqual(m.ids(), []string{"clean-vertical-slices"}) {
		t.Fatal(m.ids())
	}
}
func TestCLISetupMenusSupportMixedScopeAndCustomArchitecture(t *testing.T) {
	root := t.TempDir()
	input := "\n \x1b[B \n\x1b[F\nMy hybrid profile\npi\n \n\ny\n"
	code, out, stderr := invoke([]string{"--root", root}, input, testSource(t), true)
	if code != 0 || stderr != "" || !strings.Contains(out, "Choose project scope") || !strings.Contains(out, "Choose architecture") {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	if !strings.Contains(readInstalled(t, root, "AGENTS.md"), "Architecture: My hybrid profile") {
		t.Fatal("custom architecture lost")
	}
}
func TestCLIArchitectureChoiceStoresStableID(t *testing.T) {
	root := t.TempDir()
	options := architectureOptions(false)
	index := 0
	for options[index].ID != "clean-vertical-slices" {
		index++
	}
	args := []string{"--root", root, "--scope", "frontend", "--agents", "pi", "--skills", "frontend/design", "--instructions", "AGENTS.md", "--yes"}
	code, out, stderr := invoke(args, "\n"+strings.Repeat("\x1b[B", index)+"\n", testSource(t), true)
	if code != 0 || stderr != "" || !strings.Contains(out, "Architecture: Clean + vertical slices") {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	if !strings.Contains(readInstalled(t, root, "AGENTS.md"), "Architecture: clean-vertical-slices") {
		t.Fatal("stored profile ID is not stable")
	}
}
func TestSetupMenuCancellationDoesNotFetchOrWrite(t *testing.T) {
	for _, scope := range []string{"", "frontend"} {
		t.Run(scope, func(t *testing.T) {
			root := t.TempDir()
			args := []string{"--root", root}
			if scope != "" {
				args = append(args, "--scope", scope)
			}
			code, out, stderr := invoke(args, "\nq", Source{}, true)
			if code != 0 || stderr != "" || !strings.Contains(out, "Cancelled; no files changed.") {
				t.Fatalf("%d %q %q", code, out, stderr)
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatal("cancel mutated project", err)
			}
		})
	}
}
