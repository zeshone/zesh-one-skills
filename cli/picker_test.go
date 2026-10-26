package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func suggestionCatalog() []Skill {
	makeSkill := func(id, topic, technology string) Skill {
		parts := strings.Split(id, "/")
		return Skill{ID: id, Name: parts[1], Description: "Guidance for " + parts[1], Topics: []string{topic}, Technologies: []string{technology}}
	}
	return []Skill{
		makeSkill("shared/diagnostics", "debugging", "any"),
		makeSkill("shared/keyboard", "accessibility", "any"),
		makeSkill("frontend/react-19", "framework", "frontend"),
		makeSkill("frontend/react-18", "framework", "frontend"),
		makeSkill("frontend/tailwind-4", "design", "frontend"),
		makeSkill("frontend/typescript", "framework", "frontend"),
		makeSkill("frontend/security", "security", "frontend"),
		makeSkill("dotnet-10/data-access-persistence", "framework", "dotnet-10"),
	}
}
func TestBasicSuggestionsAreMinimalCompatibleAndCatalogDriven(t *testing.T) {
	skills := suggestionCatalog()
	cases := []struct {
		name     string
		scopes   []string
		packages map[string]int
		want     []string
	}{
		{"unknown frontend", []string{"frontend"}, nil, []string{"shared/diagnostics", "shared/keyboard"}},
		{"known stack", []string{"frontend"}, map[string]int{"react": 19, "tailwindcss": 4, "typescript": 5}, []string{"frontend/react-19", "frontend/tailwind-4", "frontend/typescript", "shared/diagnostics", "shared/keyboard"}},
		{"older framework", []string{"frontend"}, map[string]int{"react": 18}, []string{"frontend/react-18", "shared/diagnostics", "shared/keyboard"}},
		{"unknown framework major", []string{"frontend"}, map[string]int{"react": 0}, []string{"shared/diagnostics", "shared/keyboard"}},
		{"backend no UI or inferred persistence", []string{"backend"}, map[string]int{"dotnet": 10}, []string{"shared/diagnostics"}},
		{"app baseline", []string{"app"}, nil, []string{"shared/diagnostics", "shared/keyboard"}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := basicSuggestions(skills, tt.scopes, tt.packages)
			if !reflect.DeepEqual(keys(got), tt.want) {
				t.Fatalf("got %v; want %v", keys(got), tt.want)
			}
			for _, reason := range got {
				if reason == "" {
					t.Fatal("suggestion lacks explanation")
				}
			}
		})
	}
}
func TestPickerKeyboardSelectionAndReset(t *testing.T) {
	skills := suggestionCatalog()
	m := newPicker(skills, map[string]string{skills[0].ID: "basic"})
	if !m.selected[skills[0].ID] {
		t.Fatal("suggestions not preselected")
	}
	for _, key := range []string{"down", "space", "all", "none", "reset", "down", "space"} {
		m.apply(key)
	}
	if !reflect.DeepEqual(m.ids(), []string{skills[0].ID, skills[2].ID}) {
		t.Fatal(m.ids())
	}
	m.apply("enter")
	if !m.done || !m.accepted {
		t.Fatal("selection not confirmed")
	}
	cancel := newPicker(skills, nil)
	cancel.apply("cancel")
	if !cancel.done || cancel.accepted {
		t.Fatal("cancel accepted")
	}
}
func TestPickerDriverConsumesOnlySelectionKeys(t *testing.T) {
	skills := suggestionCatalog()[:3]
	input := strings.NewReader("\x1b[B \nn\n")
	reader := bufio.NewReader(input)
	var out bytes.Buffer
	ids, accepted, err := selectSkills(reader, input, &out, skills, map[string]string{skills[0].ID: "basic"})
	if err != nil || !accepted || !reflect.DeepEqual(ids, []string{skills[0].ID, skills[1].ID}) {
		t.Fatal(ids, accepted, err)
	}
	remaining, _ := io.ReadAll(reader)
	if string(remaining) != "n\n" {
		t.Fatalf("later answer consumed: %q", remaining)
	}
	for _, text := range []string{"[x]", "Guidance for", "Space", "Enter", "suggested"} {
		if !strings.Contains(out.String(), text) {
			t.Fatal("missing TUI content", text)
		}
	}
	if !strings.HasSuffix(out.String(), "\x1b[?25h\x1b[?1049l") {
		t.Fatal("screen/cursor not restored")
	}
}
func TestPickerCancellationAndEOFRestoreScreen(t *testing.T) {
	for _, input := range []string{"q", "\x03", ""} {
		t.Run(input, func(t *testing.T) {
			var out bytes.Buffer
			in := strings.NewReader(input)
			ids, accepted, err := selectSkills(bufio.NewReader(in), in, &out, suggestionCatalog(), nil)
			if accepted || len(ids) != 0 {
				t.Fatal("cancelled selection accepted")
			}
			if input == "" && err == nil {
				t.Fatal("EOF accepted")
			}
			if !strings.HasSuffix(out.String(), "\x1b[?25h\x1b[?1049l") {
				t.Fatal("terminal not restored")
			}
		})
	}
}
func TestPickerViewFitsSmallWindowAndScrolls(t *testing.T) {
	m := newPicker(suggestionCatalog(), nil)
	m.apply("end")
	view := m.view(32, 12)
	if !strings.Contains(view, "data-access-persistence") && !strings.Contains(view, "dotnet-10") {
		t.Fatal("cursor did not scroll", view)
	}
	for _, line := range strings.Split(view, "\n") {
		if len([]rune(strings.TrimSuffix(line, "\r"))) > 32 {
			t.Fatalf("line overflows: %q", line)
		}
	}
}
func TestCLISuggestedChoicesNeedNoTypedIDsAndCanBeCleared(t *testing.T) {
	for _, choice := range []string{"\n\ny\n", "n\n", "q"} {
		t.Run(choice, func(t *testing.T) {
			root := t.TempDir()
			catalog, files := specimen()
			catalog.Skills[0].Topics = []string{"accessibility"}
			catalog.Skills[0].Technologies = []string{"any"}
			input := "\n" + choice
			code, stdout, stderr := invoke([]string{"--root", root, "--scope", "frontend", "--architecture", "mockup", "--agents", "pi"}, input, testCatalogSource(t, catalog, files), true)
			if code != 0 || stderr != "" || strings.Contains(stdout, "Select skill IDs") {
				t.Fatalf("%d %q %q", code, stdout, stderr)
			}
			if choice == "\n\ny\n" {
				if !strings.Contains(stdout, "1 selected | 1 suggested") || !strings.Contains(readInstalled(t, root, "AGENTS.md"), "frontend/design") {
					t.Fatal("visible suggestion was not installed after Enter")
				}
				if _, err := os.Stat(filepath.Join(root, ".pi/skills/frontend-security")); !os.IsNotExist(err) {
					t.Fatal("unselected security installed")
				}
			} else {
				entries, err := os.ReadDir(root)
				if err != nil || len(entries) != 0 || !strings.Contains(stdout, "Cancelled") {
					t.Fatal("cleared/cancelled picker changed project", err)
				}
			}
		})
	}
}

func TestPackageMajorDetectionAvoidsAmbiguousVersions(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"react":"^19.0.0","tailwindcss":"~4.1.0","next":"workspace:*","typescript":"5.9.0"}}`), 0644)
	detected, err := Detect(root)
	if err != nil {
		t.Fatal(err)
	}
	if detected.PackageMajors["react"] != 19 || detected.PackageMajors["tailwindcss"] != 4 || detected.PackageMajors["typescript"] != 5 || detected.PackageMajors["next"] != 0 {
		t.Fatalf("%v", detected.PackageMajors)
	}
	nested := filepath.Join(root, "other")
	os.Mkdir(nested, 0755)
	os.WriteFile(filepath.Join(nested, "package.json"), []byte(`{"dependencies":{"react":"^18.0.0"}}`), 0644)
	detected, err = Detect(root)
	if err != nil || detected.PackageMajors["react"] != 0 {
		t.Fatal("conflicting project versions must not be guessed", detected.PackageMajors, err)
	}
}
