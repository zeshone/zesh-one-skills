package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var version = "dev"

const helpText = `Usage: zesh-skills [options]
Install selected project-local skills from zeshone/zesh-one-skills main.
No arguments: interactive wizard and checkbox skill picker on a terminal. Default root: current directory.
Picker: arrows/j/k move; Space toggles; Enter confirms; a all, n none, r reset; q/Ctrl-C cancels.
Setup: scope supports multiple choices; architecture selects one explained profile or custom text.
Basic suggestions are preselected, optional and editable; automation flags do not use suggestions.

  --root PATH             Project root (default: current directory)
  --scope SCOPES          Comma-separated frontend,backend,app (mixed allowed)
  --agents AGENTS         Comma-separated claude,codex,opencode,pi
  --architecture TEXT     Architecture choice or custom description
  --skills IDS            Comma-separated catalog IDs, e.g. frontend/tailwind-4
  --instructions FILE.md  Additional root Markdown instruction file
  --list                  List relevant skills by topic; never write files
  --yes                   Confirm installation, NOT replacement of local edits
  --non-interactive       Never prompt; require skills, agents, architecture and --yes
  --es                    Spanish interface (default: English)
  --version               Print version without network or writes
  --help                  Print this help without network or writes

Example:
  zesh-skills
  zesh-skills --scope frontend --list
  zesh-skills --non-interactive --agents pi --architecture mockup --skills frontend/tailwind-4 --yes

Files: .zesh-skills/content, agent-local skill entries, managed CLAUDE.md/AGENTS.md blocks.
Existing selections are retained; only selected skills are updated. Local edits require confirmation.
Exit codes: 0 success/cancel, 1 download/filesystem/conflict failure, 2 invalid usage/input.
`

func main() {
	tty := term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
	os.Exit(Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, github(), tty))
}

// Run is the CLI boundary, with streams and HTTP source injected for deterministic acceptance checks.
func Run(args []string, input io.Reader, out, errOut io.Writer, source Source, terminal bool) int {
	flags := flag.NewFlagSet("zesh-skills", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var root, scope, agentText, architecture, skillText, instructions string
	var list, yes, noPrompt, help, showVersion bool
	var spanish bool
	var lang language
	flags.StringVar(&root, "root", "", "project root")
	flags.StringVar(&scope, "scope", "", "scopes")
	flags.StringVar(&agentText, "agents", "", "agents")
	flags.StringVar(&architecture, "architecture", "", "architecture")
	flags.StringVar(&skillText, "skills", "", "skill IDs")
	flags.StringVar(&instructions, "instructions", "", "additional instruction file")
	flags.BoolVar(&list, "list", false, "list skills")
	flags.BoolVar(&yes, "yes", false, "confirm installation")
	flags.BoolVar(&noPrompt, "non-interactive", false, "never prompt")
	flags.BoolVar(&spanish, "es", false, "Spanish interface")
	flags.BoolVar(&help, "help", false, "help")
	flags.BoolVar(&showVersion, "version", false, "version")
	usage := func(err error) int {
		fmt.Fprintf(errOut, lang.text("zesh-skills: %s\nUse --help for usage.\n"), lang.errorText(err))
		return 2
	}
	failure := func(err error) int { fmt.Fprintf(errOut, "zesh-skills: %s\n", lang.errorText(err)); return 1 }
	cancelled := func() int { fmt.Fprintln(out, lang.text("Cancelled; no files changed.")); return 0 }
	parseErr := flags.Parse(args)
	lang = language(spanish)
	if parseErr != nil {
		return usage(parseErr)
	}
	if flags.NArg() != 0 {
		return usage(fmt.Errorf("unexpected positional arguments"))
	}
	if help {
		fmt.Fprint(out, lang.help())
		return 0
	}
	if showVersion {
		fmt.Fprintf(out, "zesh-skills %s\n", version)
		return 0
	}
	interactive := terminal && !noPrompt
	if !interactive && !list && (skillText == "" || agentText == "" || architecture == "" || !yes) {
		return usage(fmt.Errorf("non-interactive mode requires --skills, --agents, --architecture and --yes"))
	}
	if scope != "" {
		if _, err := parseScopes(scope); err != nil {
			return usage(err)
		}
	}
	for _, a := range split(agentText) {
		if agentDirs[a] == "" {
			return usage(fmt.Errorf("unknown agent: %s", a))
		}
	}
	if architecture != "" && !plainText(architecture, 300) {
		return usage(fmt.Errorf("invalid architecture text"))
	}
	if instructions != "" && !instructionPath(instructions) {
		return usage(fmt.Errorf("--instructions requires a root-level .md filename"))
	}
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return failure(err)
		}
	}
	reader := bufio.NewReader(input)
	ask := func(label, def string) (string, error) {
		label = lang.text(label)
		if def != "" {
			fmt.Fprintf(out, "%s [%s]: ", label, display(def))
		} else {
			fmt.Fprintf(out, "%s: ", label)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input ended; no installation started")
		}
		line = strings.TrimSpace(line)
		if line == "" {
			line = def
		}
		return line, nil
	}
	if interactive && !list {
		var err error
		root, err = ask("Project root", root)
		if err != nil {
			return usage(err)
		}
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return usage(fmt.Errorf("invalid project root"))
	}
	root = absolute
	detected, err := Detect(root)
	if err != nil {
		return failure(err)
	}
	architectureDisplay := architecture
	if interactive && !list {
		fmt.Fprintf(out, lang.text("Project: %s\nDetected technologies: %s\nExisting instructions: %s\n"), display(filepath.Base(root)), strings.Join(detected.Technologies, ", "), strings.Join(detected.Instructions, ", "))
		if scope == "" {
			defaults := map[string]string{}
			for _, value := range detected.Scopes {
				defaults[value] = "Detected project scope"
			}
			selected, accepted, e := selectOptions(reader, input, out, newOptionPicker(scopeOptions(bool(lang)), defaults, "Choose project scope", false).inLanguage(lang))
			if e != nil {
				return failure(e)
			}
			if !accepted || len(selected) == 0 {
				return cancelled()
			}
			scope = strings.Join(selected, ",")
		}
		if architecture == "" {
			options := architectureOptions(bool(lang))
			selected, accepted, e := selectOptions(reader, input, out, newOptionPicker(options, nil, "Choose architecture", true).inLanguage(lang))
			if e != nil {
				return failure(e)
			}
			if !accepted || len(selected) == 0 {
				return cancelled()
			}
			architecture = selected[0]
			for _, option := range options {
				if option.ID == architecture {
					architectureDisplay = option.Label
				}
			}
			if architecture == "custom" {
				architecture, err = ask("Custom architecture", "")
				if err != nil {
					return usage(err)
				}
				if !plainText(architecture, 300) {
					return usage(fmt.Errorf("invalid architecture text"))
				}
				architectureDisplay = architecture
			}
		}
		if agentText == "" {
			agentText, err = ask("Agents (claude,codex,opencode,pi)", strings.Join(detected.Agents, ","))
			if err != nil {
				return usage(err)
			}
		}
	}
	scopes := detected.Scopes
	if scope != "" {
		scopes, err = parseScopes(scope)
		if err != nil {
			return usage(err)
		}
	}
	if len(scopes) == 0 {
		if list || !interactive {
			scopes = []string{"frontend", "backend", "app"}
		} else {
			return usage(fmt.Errorf("choose at least one scope"))
		}
	}
	agents := split(agentText)
	for _, a := range agents {
		if agentDirs[a] == "" {
			return usage(fmt.Errorf("unknown agent: %s", a))
		}
	}
	if !list && (len(agents) == 0 || !plainText(architecture, 300)) {
		return usage(fmt.Errorf("agents and architecture are required"))
	}
	c, commit, err := source.Load()
	if err != nil {
		return failure(err)
	}
	relevant := relevantSkills(c, scopes)
	if list || (interactive && skillText != "") {
		showSkillsLanguage(out, relevant, lang)
	}
	if list {
		return 0
	}
	if skillText == "" {
		selected, accepted, e := selectOptions(reader, input, out, newPicker(relevant, basicSuggestions(relevant, scopes, detected.PackageMajors)).inLanguage(lang))
		if e != nil {
			return failure(e)
		}
		if !accepted || len(selected) == 0 {
			fmt.Fprintln(out, lang.text("Cancelled; no files changed."))
			return 0
		}
		skillText = strings.Join(selected, ",")
	}
	if skillText == "" || skillText == "none" {
		fmt.Fprintln(out, lang.text("Cancelled; no files changed."))
		return 0
	}
	ids := split(skillText)
	lookup := map[string]bool{}
	for _, s := range c.Skills {
		lookup[s.ID] = true
	}
	for i, id := range ids {
		if n, e := strconv.Atoi(id); e == nil {
			if n < 1 || n > len(relevant) {
				return usage(fmt.Errorf("invalid skill number: %s", id))
			}
			ids[i] = relevant[n-1].ID
		}
		if !lookup[ids[i]] {
			return usage(fmt.Errorf("unknown skill: %s", ids[i]))
		}
	}
	ids = unique(ids)
	if interactive && instructions == "" {
		instructions, err = ask("Additional root instruction file (.md; optional)", "")
		if err != nil {
			return usage(err)
		}
	}
	if instructions != "" && !instructionPath(instructions) {
		return usage(fmt.Errorf("instruction file must be a root-level .md filename"))
	}
	fmt.Fprintf(out, lang.text("\nInstall %d selected skill(s) into %s\nSource: main @ %s\nArchitecture: %s\nAgents: %s\nSkills: %s\nContent: .zesh-skills/content/\n"), len(ids), display(root), commit, display(architectureDisplay), strings.Join(agents, ", "), strings.Join(ids, ", "))
	fmt.Fprint(out, lang.text("Instruction files: "))
	targets := map[string]bool{}
	for _, a := range agents {
		target := "AGENTS.md"
		if a == "claude" {
			target = "CLAUDE.md"
		}
		for _, existing := range detected.Instructions {
			if strings.EqualFold(existing, target) {
				target = existing
			}
		}
		targets[target] = true
		fmt.Fprintf(out, "\n  %s/<scope>-<name>/SKILL.md", agentDirs[a])
	}
	if instructions != "" {
		targets[instructions] = true
	}
	fmt.Fprintf(out, lang.text("\n  %s (managed index and triggers)\n"), strings.Join(keys(targets), ", "))
	if !yes {
		answer, e := ask("Proceed? (y/N)", "n")
		if e != nil {
			return usage(e)
		}
		if !lang.affirmative(answer) {
			fmt.Fprintln(out, lang.text("Cancelled; no files changed."))
			return 0
		}
	}
	confirm := func(name string) bool {
		if !interactive {
			return false
		}
		answer, e := ask(fmt.Sprintf(lang.text("Replace local changes in %s? (y/N)"), display(lang.confirmationName(name))), "n")
		return e == nil && lang.affirmative(answer)
	}
	in := Installer{Fetch: source.Fetch}
	if instructions != "" {
		in.Instructions = []string{instructions}
	}
	if err = in.Install(root, c, commit, ids, agents, architecture, confirm); err != nil {
		return failure(err)
	}
	fmt.Fprintf(out, lang.text("Installed %d selected skill(s). Existing selections retained.\n"), len(ids))
	return 0
}
func split(s string) []string {
	parts := []string{}
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			parts = append(parts, p)
		}
	}
	return unique(parts)
}
func unique(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, s := range values {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
func parseScopes(s string) ([]string, error) {
	out := []string{}
	for _, p := range split(s) {
		switch p {
		case "frontend", "front":
			out = append(out, "frontend")
		case "backend", "back":
			out = append(out, "backend")
		case "app", "apps":
			out = append(out, "app")
		case "mixed":
			out = append(out, "frontend", "backend", "app")
		default:
			return nil, fmt.Errorf("unknown scope: %s", p)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("scope is required")
	}
	return unique(out), nil
}
func affirmative(s string) bool { return strings.EqualFold(s, "y") || strings.EqualFold(s, "yes") }
func display(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return '?'
		}
		return r
	}, s)
}
func relevantSkills(c Catalog, scopes []string) []Skill {
	out := []Skill{}
	for _, s := range c.Skills {
		match := false
		for _, a := range scopes {
			for _, b := range s.Scopes {
				match = match || a == b
			}
		}
		if match {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := strings.Join(out[i].Topics, ","), strings.Join(out[j].Topics, ",")
		if a == b {
			return out[i].ID < out[j].ID
		}
		return a < b
	})
	return out
}
func brief(s string) string {
	runes := []rune(s)
	if len(runes) > 180 {
		return string(runes[:177]) + "..."
	}
	return s
}

func showSkills(out io.Writer, skills []Skill) { showSkillsLanguage(out, skills, false) }

func showSkillsLanguage(out io.Writer, skills []Skill, lang language) {
	last := ""
	for i, s := range skills {
		topic := strings.Join(s.Topics, ", ")
		if topic != last {
			// Preserve --list's comma-space contract; picker rows use comma only.
			labels := make([]string, len(s.Topics))
			for index, value := range s.Topics {
				labels[index] = lang.text(value)
			}
			fmt.Fprintf(out, "\n%s\n", strings.Join(labels, ", "))
			last = topic
		}
		fmt.Fprintf(out, "  [ ] %d. %s — %s\n", i+1, s.ID, brief(s.Description))
	}
}

type Detection struct {
	Scopes, Technologies, Instructions, Agents []string
	PackageMajors                              map[string]int
}

func Detect(root string) (Detection, error) {
	result := Detection{PackageMajors: map[string]int{}}
	dir, err := os.OpenRoot(root)
	if err != nil {
		return result, err
	}
	defer dir.Close()
	scopes, tech := map[string]bool{}, map[string]bool{}
	// Inspect manifests only. Never invoke project tools, package hooks or recursively scan dependencies.
	var scan func(string, int) error
	count := 0
	scan = func(folder string, depth int) error {
		f, e := dir.Open(folder)
		if e != nil {
			return e
		}
		entries, e := f.ReadDir(500)
		f.Close()
		if e != nil && !errors.Is(e, io.EOF) {
			return e
		}
		for _, entry := range entries {
			count++
			if count > 1000 {
				return nil
			}
			name := path.Join(folder, entry.Name())
			if entry.Type()&os.ModeSymlink != 0 || strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			lower := strings.ToLower(entry.Name())
			if entry.IsDir() {
				if depth < 2 && lower != "node_modules" && lower != "vendor" && lower != "bin" && lower != "obj" {
					if e = scan(name, depth+1); e != nil {
						return e
					}
				}
				continue
			}
			if strings.HasSuffix(lower, ".csproj") || strings.HasSuffix(lower, ".sln") || strings.HasSuffix(lower, ".slnx") {
				scopes["backend"] = true
				tech["dotnet"] = true
			}
			if lower == "package.json" {
				info, e := dir.Stat(name)
				if e != nil || info.Size() > 1024*1024 {
					continue
				}
				data, e := dir.ReadFile(name)
				if e != nil {
					continue
				}
				var pkg struct {
					Dependencies    map[string]string `json:"dependencies"`
					DevDependencies map[string]string `json:"devDependencies"`
				}
				if json.Unmarshal(data, &pkg) != nil {
					continue
				}
				for dep, version := range pkg.Dependencies {
					if pkg.DevDependencies == nil {
						pkg.DevDependencies = map[string]string{}
					}
					pkg.DevDependencies[dep] = version
				}
				for dep, version := range pkg.DevDependencies {
					major := declaredMajor(version)
					if previous, exists := result.PackageMajors[dep]; exists && previous != major {
						major = 0 // Conflicting monorepo declarations must not force a profile.
					}
					result.PackageMajors[dep] = major
					switch {
					case dep == "react" || dep == "next" || dep == "@angular/core" || dep == "vue" || dep == "vite":
						scopes["frontend"] = true
						tech[dep] = true
					case strings.HasPrefix(dep, "@ionic/") || strings.HasPrefix(dep, "@capacitor/") || dep == "react-native" || dep == "expo":
						scopes["app"] = true
						tech[dep] = true
					case dep == "express" || dep == "fastify" || dep == "@nestjs/core":
						scopes["backend"] = true
						tech[dep] = true
					}
				}
			}
			if folder == "." && (strings.EqualFold(lower, "agents.md") || strings.EqualFold(lower, "claude.md") || strings.EqualFold(lower, "gemini.md") || strings.EqualFold(lower, "copilot-instructions.md")) {
				result.Instructions = append(result.Instructions, entry.Name())
			}
		}
		return nil
	}
	if err = scan(".", 0); err != nil {
		return result, err
	}
	result.Scopes = keys(scopes)
	result.Technologies = keys(tech)
	sort.Strings(result.Instructions)
	for _, a := range []string{"claude", "codex", "opencode", "pi"} {
		if info, e := dir.Lstat(strings.Split(agentDirs[a], "/")[0]); e == nil && info.IsDir() {
			result.Agents = append(result.Agents, a)
		}
	}
	if len(result.Agents) == 0 {
		for _, name := range result.Instructions {
			if strings.EqualFold(name, "CLAUDE.md") {
				result.Agents = append(result.Agents, "claude")
			}
			if strings.EqualFold(name, "AGENTS.md") {
				result.Agents = append(result.Agents, "codex")
			}
		}
	}
	return result, nil
}
