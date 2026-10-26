package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpanishHelpAndErrorsNeedNoNetworkOrWrites(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	code, out, stderr := invoke([]string{"--es", "--help"}, "", Source{}, false)
	if code != 0 || stderr != "" || !strings.HasPrefix(out, "Uso: zesh-skills") || !strings.Contains(out, "--es") || !strings.Contains(out, "arquitectura") {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	for _, tt := range []struct {
		args           []string
		code           int
		stdout, stderr string
	}{
		{[]string{"--es", "--version"}, 0, "zesh-skills dev\n", ""},
		{[]string{"--es", "--unknown"}, 2, "", "zesh-skills: opción no definida: -unknown\nUsa --help para ver el uso.\n"},
		{[]string{"--es", "--scope", "wrong", "--list"}, 2, "", "zesh-skills: alcance desconocido: wrong\nUsa --help para ver el uso.\n"},
	} {
		code, stdout, stderr := invoke(tt.args, "", Source{}, false)
		if code != tt.code || stdout != tt.stdout || stderr != tt.stderr {
			t.Fatalf("%v: %d %q %q", tt.args, code, stdout, stderr)
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("side effects", err)
	}
}
func TestSpanishWizardPreservesCustomValuesAndAcceptsSi(t *testing.T) {
	root := t.TempDir()
	input := "\n \x1b[B \n\x1b[F\nMi híbrido personal\npi\n \n\nsí\n"
	code, out, stderr := invoke([]string{"--es", "--root", root}, input, testSource(t), true)
	if code != 0 || stderr != "" {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	for _, text := range []string{"Raíz del proyecto", "Selecciona el alcance del proyecto", "Elige la arquitectura", "Arquitectura personalizada", "Selecciona las skills del proyecto", "Skills a instalar: 1", "¿Continuar?", "Instalación completada"} {
		if !strings.Contains(out, text) {
			t.Fatalf("missing Spanish UI %q", text)
		}
	}
	if strings.Contains(out, "Proceed?") || strings.Contains(out, "Choose architecture") {
		t.Fatal("English UI leaked")
	}
	if !strings.Contains(readInstalled(t, root, "AGENTS.md"), "Architecture: Mi híbrido personal") {
		t.Fatal("custom value translated or protocol changed")
	}
}
func TestLanguageDoesNotChangeInstalledProtocols(t *testing.T) {
	source := testSource(t)
	roots := []string{t.TempDir(), t.TempDir()}
	for i, root := range roots {
		args := verifyArgs(root)
		if i == 1 {
			args = append(args, "--es")
		}
		code, out, stderr := invoke(args, "", source, false)
		if code != 0 || stderr != "" {
			t.Fatal(code, out, stderr)
		}
	}
	for _, file := range []string{"AGENTS.md", manifestPath, ".pi/skills/frontend-design/SKILL.md", ".zesh-skills/content/skills/frontend/design/SKILL.md"} {
		if readInstalled(t, roots[0], file) != readInstalled(t, roots[1], file) {
			t.Fatal("language changed stored content", file)
		}
	}
}
func TestSpanishListAndConflictRejectionPreserveChoices(t *testing.T) {
	root := t.TempDir()
	code, out, stderr := invoke([]string{"--es", "--root", root, "--scope", "frontend", "--list"}, "", testSource(t), false)
	want := "\nDiseño\n  [ ] 1. frontend/design — UI layout\n\nSeguridad\n  [ ] 2. frontend/security — Security checks\n"
	if code != 0 || out != want || stderr != "" {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	args := append(verifyArgs(root), "--es")
	if code, _, stderr = invoke(args, "", testSource(t), false); code != 0 || stderr != "" {
		t.Fatal(code, stderr)
	}
	name := filepath.Join(root, ".zesh-skills/content/skills/frontend/design/SKILL.md")
	os.WriteFile(name, []byte("local edits"), 0644)
	before := readInstalled(t, root, manifestPath)
	code, out, stderr = invoke(args, "", testSource(t), false)
	if code != 1 || stderr != "zesh-skills: cancelado: se conservaron los cambios locales en .zesh-skills/content/skills/frontend/design/SKILL.md\n" {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	data, _ := os.ReadFile(name)
	if string(data) != "local edits" || readInstalled(t, root, manifestPath) != before {
		t.Fatal("Spanish rejection modified data")
	}
}
func TestSpanishHybridShowsTranslatedLabelButStoresCanonicalID(t *testing.T) {
	root := t.TempDir()
	options := architectureOptions(true)
	index := 0
	for options[index].ID != "microservices-events" {
		index++
	}
	args := []string{"--es", "--root", root, "--scope", "frontend", "--agents", "pi", "--skills", "frontend/design", "--instructions", "AGENTS.md", "--yes"}
	code, out, stderr := invoke(args, "\n"+strings.Repeat("\x1b[B", index)+"\n", testSource(t), true)
	if code != 0 || stderr != "" || !strings.Contains(out, "Arquitectura: Microservicios + eventos") {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	if !strings.Contains(readInstalled(t, root, "AGENTS.md"), "Architecture: microservices-events") {
		t.Fatal("translated label leaked into stored profile")
	}
	code, out, stderr = invoke([]string{"--es=false", "--help"}, "", Source{}, false)
	if code != 0 || stderr != "" || out != helpText {
		t.Fatal("English default/explicit opt-out changed")
	}
}

func TestListMultiTopicSpacingPreservesEnglishAndPickerContracts(t *testing.T) {
	for _, spanish := range []bool{false, true} {
		root := t.TempDir()
		catalog, files := specimen()
		catalog.Skills[0].Topics = []string{"design", "testing"}
		args := []string{"--root", root, "--scope", "frontend", "--list"}
		heading, security, pickerTopics := "design, testing", "security", "[design,testing]"
		if spanish {
			args = append(args, "--es")
			heading = "Diseño, Pruebas"
			security = "Seguridad"
			pickerTopics = "[Diseño,Pruebas]"
		}
		code, stdout, stderr := invoke(args, "", testCatalogSource(t, catalog, files), false)
		want := "\n" + heading + "\n  [ ] 1. frontend/design — UI layout\n\n" + security + "\n  [ ] 2. frontend/security — Security checks\n"
		if code != 0 || stdout != want || stderr != "" {
			t.Fatalf("Spanish=%v: %d %q %q", spanish, code, stdout, stderr)
		}
		model := newPicker(catalog.Skills, nil).inLanguage(language(spanish))
		if !strings.Contains(model.view(120, 30), pickerTopics) {
			t.Fatal("picker separator changed")
		}
		entries, err := os.ReadDir(root)
		if err != nil || len(entries) != 0 {
			t.Fatal("listing changed project", err)
		}
	}
}

func TestSpanishConfirmationCanDeclineWithoutWriting(t *testing.T) {
	root := t.TempDir()
	args := []string{"--es", "--root", root, "--scope", "frontend", "--architecture", "Architecture user value", "--agents", "pi", "--skills", "frontend/design", "--instructions", "AGENTS.md"}
	code, out, stderr := invoke(args, "\nn\n", testSource(t), true)
	if code != 0 || stderr != "" || !strings.Contains(out, "Cancelado; no se modificaron archivos.") || !strings.Contains(out, "Architecture user value") {
		t.Fatalf("%d %q %q", code, out, stderr)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatal("decline changed project")
	}
}
