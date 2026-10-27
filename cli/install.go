package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

const manifestPath = ".zesh-skills/manifest.json"
const blockStart = "<!-- zesh-skills:start -->"
const blockEnd = "<!-- zesh-skills:end -->"

var agentDirs = map[string]string{"claude": ".claude/skills", "codex": ".agents/skills", "opencode": ".opencode/skills", "pi": ".pi/skills"}

type Manifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	Commit        string            `json:"commit"`
	Architecture  string            `json:"architecture"`
	Skills        map[string]Skill  `json:"skills"`
	Agents        []string          `json:"agents"`
	Files         map[string]string `json:"files"`
	Blocks        map[string]string `json:"blocks"`
}
type Installer struct {
	Fetch        func(string, string) ([]byte, error)
	Instructions []string
}
type backup struct {
	data   []byte
	mode   os.FileMode
	exists bool
}

func (in Installer) Install(root string, c Catalog, commit string, ids, agents []string, architecture string, confirm func(string) bool) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if !commitID.MatchString(commit) {
		return fmt.Errorf("invalid source commit")
	}
	if len(ids) == 0 || len(agents) == 0 || !plainText(architecture, 300) {
		return fmt.Errorf("skills, agents and architecture are required")
	}
	available := map[string]Skill{}
	for _, s := range c.Skills {
		available[s.ID] = s
	}
	for _, id := range ids {
		if _, ok := available[id]; !ok {
			return fmt.Errorf("unknown skill: %s", id)
		}
	}
	for _, a := range agents {
		if agentDirs[a] == "" {
			return fmt.Errorf("unknown agent: %s", a)
		}
	}
	// OpenRoot confines all subsequent filesystem operations even if an ancestor changes.
	dir, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer dir.Close()
	if err = rejectLinks(dir, manifestPath); err != nil {
		return err
	}
	lock, err := dir.OpenFile(".zesh-skills.lock", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("project is locked or not writable; inspect .zesh-skills.lock")
	}
	lock.Close()
	defer dir.Remove(".zesh-skills.lock")
	state := Manifest{SchemaVersion: 1, Skills: map[string]Skill{}, Files: map[string]string{}, Blocks: map[string]string{}}
	old, err := dir.ReadFile(manifestPath)
	if err == nil {
		if json.Unmarshal(old, &state) != nil || state.SchemaVersion != 1 || state.Skills == nil || state.Files == nil || state.Blocks == nil {
			return fmt.Errorf("invalid local manifest; preserve it and inspect manually")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for id, s := range state.Skills {
		if id != s.ID {
			return fmt.Errorf("invalid local skill identity")
		}
		if err = (Catalog{SchemaVersion: 1, Skills: []Skill{s}}).Validate(); err != nil {
			return fmt.Errorf("invalid local skill metadata: %w", err)
		}
	}
	for _, a := range state.Agents {
		if agentDirs[a] == "" {
			return fmt.Errorf("invalid local agent")
		}
	}
	for name, h := range state.Files {
		if !managedPath(name) || !digest.MatchString(h) {
			return fmt.Errorf("invalid local manifest file")
		}
	}
	for name, h := range state.Blocks {
		if !instructionPath(name) || !digest.MatchString(h) {
			return fmt.Errorf("invalid local instruction ownership")
		}
	}
	writes := map[string][]byte{}
	total := 0
	for _, id := range ids {
		s := available[id]
		for _, f := range s.Files {
			target := ".zesh-skills/content/" + f.Path
			if _, ok := writes[target]; ok {
				continue
			}
			data, e := in.Fetch(commit, f.Path)
			if e != nil {
				return e
			}
			if hash(data) != f.SHA256 {
				return fmt.Errorf("checksum mismatch: %s", f.Path)
			}
			total += len(data)
			if total > 64*1024*1024 {
				return fmt.Errorf("installation exceeds 64 MiB")
			}
			writes[target] = data
		}
		state.Skills[id] = s
	}
	allAgents := map[string]bool{}
	for _, a := range append(state.Agents, agents...) {
		allAgents[a] = true
	}
	state.Agents = keys(allAgents)
	state.Commit = commit
	state.Architecture = architecture
	// Keep relative links when their destination is installed; pin other repository links remotely.
	for name, data := range writes {
		if strings.HasSuffix(name, ".md") {
			writes[name] = rewriteLinks(data, strings.TrimPrefix(name, ".zesh-skills/content/"), commit, writes, state.Files)
		}
	}
	instructions := map[string]bool{}
	for name := range state.Blocks {
		instructions[name] = true
	}
	for _, a := range state.Agents {
		for _, s := range state.Skills {
			alias := strings.ReplaceAll(s.ID, "/", "-")
			wrapper := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\nRead [the complete skill](../../../.zesh-skills/content/%s/SKILL.md) before working on this topic. Resolve its resource links relative to that file.\n", alias, jsonString(s.Description), s.Path)
			writes[agentDirs[a]+"/"+alias+"/SKILL.md"] = []byte(wrapper)
		}
		name := "AGENTS.md"
		if a == "claude" {
			name = "CLAUDE.md"
		}
		name, err = existingCase(dir, name)
		if err != nil {
			return err
		}
		instructions[name] = true
	}
	for _, name := range in.Instructions {
		if !instructionPath(name) {
			return fmt.Errorf("instruction file must be a root-level .md filename")
		}
		instructions[name] = true
	}
	block := renderBlock(state)
	instructionSnapshots := map[string]string{}
	for name := range instructions {
		if err = rejectLinks(dir, name); err != nil {
			return err
		}
		data, e := dir.ReadFile(name)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		instructionSnapshots[name] = string(data)
		updated, oldBlock, e := mergeBlock(data, block)
		if e != nil {
			return fmt.Errorf("%s: %w", name, e)
		}
		if len(oldBlock) > 0 && (state.Blocks[name] == "" || hash(oldBlock) != state.Blocks[name]) && string(oldBlock) != string(block) {
			if confirm == nil || !confirm(name+" (modified managed index)") {
				return fmt.Errorf("cancelled: local changes preserved in %s", name)
			}
		}
		writes[name] = updated
		state.Blocks[name] = hash(block)
	}
	names := keys(writes)
	backups := map[string]backup{}
	for _, name := range names {
		if err = rejectLinks(dir, name); err != nil {
			return err
		}
		data, e := dir.ReadFile(name)
		exists := e == nil
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		mode := os.FileMode(0644)
		if exists {
			info, e := dir.Stat(name)
			if e != nil {
				return e
			}
			mode = info.Mode().Perm()
		}
		backups[name] = backup{data, mode, exists}
		if !instructions[name] && exists && string(data) != string(writes[name]) && (state.Files[name] == "" || hash(data) != state.Files[name]) {
			if confirm == nil || !confirm(name) {
				return fmt.Errorf("cancelled: local changes preserved in %s", name)
			}
		}
		if !instructions[name] {
			state.Files[name] = hash(writes[name])
		}
	}
	// Confirmations may wait while an editor saves unrelated instruction text.
	// Reject drift after every prompt and before any project content is written.
	for name, snapshot := range instructionSnapshots {
		if err = rejectLinks(dir, name); err != nil {
			return err
		}
		current, e := dir.ReadFile(name)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		if string(current) != snapshot {
			return fmt.Errorf("instruction changed during installation: %s; retry", name)
		}
	}
	names = append(names, manifestPath)
	backups[manifestPath] = backup{old, 0600, old != nil}
	for _, name := range names {
		if err = rejectLinks(dir, name); err != nil {
			return err
		}
		current, e := dir.ReadFile(name)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		snapshot := backups[name]
		if (e == nil) != snapshot.exists || string(current) != string(snapshot.data) {
			return fmt.Errorf("file changed during installation: %s; retry", name)
		}
	}
	encoded, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	writes[manifestPath] = append(encoded, '\n')
	// The manifest is written last. Roll back ordinary write failures; power loss is not a transaction.
	completed := []string{}
	createdDirs := []string{}
	rollback := func(cause error) error {
		var failures []string
		for i := len(completed) - 1; i >= 0; i-- {
			name := completed[i]
			b := backups[name]
			var e error
			if b.exists {
				e = dir.WriteFile(name, b.data, b.mode)
			} else {
				e = dir.Remove(name)
			}
			if e != nil {
				failures = append(failures, name)
			}
		}
		for i := len(createdDirs) - 1; i >= 0; i-- {
			_ = dir.Remove(createdDirs[i])
		}
		if len(failures) > 0 {
			return fmt.Errorf("%w; rollback incomplete for %s", cause, strings.Join(failures, ", "))
		}
		return cause
	}
	for _, name := range names {
		if err = rejectLinks(dir, name); err != nil {
			return rollback(err)
		}
		parent := path.Dir(name)
		parts := strings.Split(parent, "/")
		current := ""
		for _, part := range parts {
			if part == "." {
				continue
			}
			current = path.Join(current, part)
			e := dir.Mkdir(current, 0755)
			if e == nil {
				createdDirs = append(createdDirs, current)
			} else if !errors.Is(e, os.ErrExist) {
				return rollback(e)
			}
		}
		// Write through the rooted handle. Record before writing so partial failures restore originals.
		completed = append(completed, name)
		if err = dir.WriteFile(name, writes[name], backups[name].mode); err != nil {
			return rollback(err)
		}
	}
	return nil
}
func managedPath(name string) bool {
	prefix := ".zesh-skills/content/"
	if strings.HasPrefix(name, prefix) {
		p := strings.TrimPrefix(name, prefix)
		return safePath(p) && strings.HasPrefix(p, "skills/")
	}
	for _, dir := range agentDirs {
		if strings.HasPrefix(name, dir+"/") {
			return safePath(strings.TrimPrefix(name, dir+"/"))
		}
	}
	return false
}
func instructionPath(name string) bool {
	return safePath(name) && !strings.Contains(name, "/") && strings.EqualFold(path.Ext(name), ".md")
}
func rejectLinks(root *os.Root, name string) error {
	current := ""
	for _, part := range strings.Split(name, "/") {
		current = path.Join(current, part)
		info, err := root.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing symbolic link: %s", current)
		}
		if current != name && !info.IsDir() {
			return fmt.Errorf("not a directory: %s", current)
		}
	}
	return nil
}
func existingCase(root *os.Root, want string) (string, error) {
	f, err := root.Open(".")
	if err != nil {
		return "", err
	}
	defer f.Close()
	entries, err := f.ReadDir(-1)
	if err != nil {
		return "", err
	}
	found := ""
	for _, entry := range entries {
		if strings.EqualFold(entry.Name(), want) {
			if found != "" {
				return "", fmt.Errorf("ambiguous instruction filenames: %s", want)
			}
			found = entry.Name()
		}
	}
	if found == "" {
		found = want
	}
	return found, nil
}
func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func jsonString(s string) string { b, _ := json.Marshal(s); return string(b) }
func md(s string) string {
	r := strings.NewReplacer("|", "\\|", "<", "&lt;", ">", "&gt;", "`", "'")
	return r.Replace(s)
}
func renderBlock(m Manifest) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n## Installed Zesh skills\n\nArchitecture: %s\n\nLoad only the relevant selected skill before working on its topic.\n\n| Skill | Description | Usage triggers | Instructions |\n| --- | --- | --- | --- |\n", blockStart, md(m.Architecture))
	for _, id := range keys(m.Skills) {
		s := m.Skills[id]
		fmt.Fprintf(&b, "| %s | %s | %s | [%s](.zesh-skills/content/%s/SKILL.md) |\n", id, md(s.Description), md(strings.Join(s.Triggers, "; ")), id, s.Path)
	}
	b.WriteString(blockEnd)
	return []byte(b.String())
}
func mergeBlock(old, block []byte) ([]byte, []byte, error) {
	s := string(old)
	starts, ends := strings.Count(s, blockStart), strings.Count(s, blockEnd)
	if starts == 0 && ends == 0 {
		sep := "\n\n"
		if len(old) == 0 {
			sep = ""
		}
		return []byte(s + sep + string(block) + "\n"), nil, nil
	}
	start, end := strings.Index(s, blockStart), strings.Index(s, blockEnd)
	if starts != 1 || ends != 1 || start > end {
		return nil, nil, fmt.Errorf("malformed managed index markers; no files changed")
	}
	end += len(blockEnd)
	return []byte(s[:start] + string(block) + s[end:]), []byte(s[start:end]), nil
}

var markdownLink = regexp.MustCompile(`\]\(([^\s)]+)\)`)

func rewriteLinks(data []byte, name, commit string, writes map[string][]byte, installed map[string]string) []byte {
	return markdownLink.ReplaceAllFunc(data, func(match []byte) []byte {
		target := string(match[2 : len(match)-1])
		if strings.HasPrefix(target, "#") || strings.HasPrefix(target, "/") || strings.Contains(target, ":") {
			return match
		}
		parts := strings.SplitN(target, "#", 2)
		dest := path.Clean(path.Join(path.Dir(name), parts[0]))
		local := ".zesh-skills/content/" + dest
		if _, ok := writes[local]; ok {
			return match
		}
		if _, ok := installed[local]; ok {
			return match
		}
		if !safePath(dest) {
			return match
		}
		suffix := ""
		if len(parts) == 2 {
			suffix = "#" + parts[1]
		}
		return []byte("](https://github.com/zeshone/zesh-one-skills/blob/" + commit + "/" + dest + suffix + ")")
	})
}
