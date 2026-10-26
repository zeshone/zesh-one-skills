package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Recommendations are conservative hints, never a hidden installation rule.
// Catalog names/topics drive discovery; no inventory of skill IDs lives here.
func basicSuggestions(skills []Skill, scopes []string, packages map[string]int) map[string]string {
	type candidate struct {
		id, reason string
		rank       int
	}
	candidates := []candidate{}
	ui := false
	for _, scope := range scopes {
		ui = ui || scope == "frontend" || scope == "app"
	}
	for _, s := range skills {
		topics := map[string]bool{}
		for _, topic := range s.Topics {
			topics[topic] = true
		}
		neutral := false
		for _, technology := range s.Technologies {
			neutral = neutral || technology == "any"
		}
		switch {
		case neutral && topics["debugging"]:
			candidates = append(candidates, candidate{s.ID, "Framework-neutral debugging baseline", 1})
		case neutral && topics["accessibility"] && ui:
			candidates = append(candidates, candidate{s.ID, "UI keyboard, focus and feedback baseline", 2})
		case topics["framework"] || topics["design"]:
			base, major := skillTechnology(s.Name)
			for _, dependency := range keys(packages) {
				name := strings.ReplaceAll(strings.TrimPrefix(dependency, "@"), "/", "-")
				match := name == base || name == base+"css" || base == name+"js" || strings.HasPrefix(name, base+"-")
				if match && (major == 0 || packages[dependency] == major) {
					candidates = append(candidates, candidate{s.ID, "Matches project dependency " + dependency, 3})
					break
				}
			}
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].rank != candidates[j].rank {
			return candidates[i].rank > candidates[j].rank
		}
		return candidates[i].id < candidates[j].id
	})
	suggestions := map[string]string{}
	for _, c := range candidates {
		if len(suggestions) == 5 {
			break
		}
		suggestions[c.id] = c.reason
	}
	return suggestions
}

var skillMajor = regexp.MustCompile(`^(.*)-([0-9]+)$`)
var packageMajor = regexp.MustCompile(`^[~^]?([0-9]+)(?:\.[0-9xX*]+)*(?:-[a-zA-Z0-9.-]+)?$`)

func skillTechnology(name string) (string, int) {
	if parts := skillMajor.FindStringSubmatch(name); parts != nil {
		major, _ := strconv.Atoi(parts[2])
		return parts[1], major
	}
	return name, 0
}
func declaredMajor(version string) int {
	if parts := packageMajor.FindStringSubmatch(version); parts != nil {
		major, _ := strconv.Atoi(parts[1])
		return major
	}
	return 0
}

type pickerOption struct {
	ID, Label, Description string
	Topics                 []string
}

type pickerModel struct {
	options        []pickerOption
	selected       map[string]bool
	suggestions    map[string]string
	cursor         int
	done, accepted bool
	title          string
	single         bool
	lang           language
}

func newPicker(skills []Skill, suggestions map[string]string) *pickerModel {
	options := []pickerOption{}
	for _, skill := range skills {
		options = append(options, pickerOption{ID: skill.ID, Label: skill.ID, Description: skill.Description, Topics: skill.Topics})
	}
	return newOptionPicker(options, suggestions, "Select project skills", false)
}
func newOptionPicker(options []pickerOption, suggestions map[string]string, title string, single bool) *pickerModel {
	m := &pickerModel{options: options, selected: map[string]bool{}, suggestions: suggestions, title: title, single: single}
	m.apply("reset")
	return m
}
func (m *pickerModel) apply(key string) {
	switch key {
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor+1 < len(m.options) {
			m.cursor++
		}
	case "home":
		m.cursor = 0
	case "end":
		if len(m.options) > 0 {
			m.cursor = len(m.options) - 1
		}
	case "space":
		if len(m.options) > 0 {
			id := m.options[m.cursor].ID
			m.selected[id] = !m.selected[id]
		}
	case "all":
		if m.single {
			break
		}
		for _, s := range m.options {
			m.selected[s.ID] = true
		}
	case "none":
		if m.single {
			break
		}
		m.selected = map[string]bool{}
	case "reset":
		m.selected = map[string]bool{}
		for _, s := range m.options {
			if m.suggestions[s.ID] != "" {
				m.selected[s.ID] = true
			}
		}
	case "enter":
		m.done = true
		m.accepted = true
	case "cancel":
		m.done = true
	}
	if m.single && len(m.options) > 0 {
		m.selected = map[string]bool{m.options[m.cursor].ID: true}
	}
}
func (m *pickerModel) ids() []string {
	ids := []string{}
	for _, s := range m.options {
		if m.selected[s.ID] {
			ids = append(ids, s.ID)
		}
	}
	return ids
}
func clip(s string, width int) string {
	r := []rune(display(s))
	if width < 1 {
		return ""
	}
	if len(r) <= width {
		return string(r)
	}
	if width < 4 {
		return string(r[:width])
	}
	return string(r[:width-3]) + "..."
}
func (m *pickerModel) view(width, height int) string {
	if width < 1 {
		width = 1
	}
	visible := (height - 9) / 2
	if visible < 1 {
		visible = 1
	}
	if visible > len(m.options) {
		visible = len(m.options)
	}
	start := 0
	if m.cursor >= visible && visible > 0 {
		start = m.cursor - visible + 1
	}
	end := start + visible
	if end > len(m.options) {
		end = len(m.options)
	}
	lines := []string{m.lang.text(m.title), fmt.Sprintf(m.lang.text("%d selected | %d suggested | %d available"), len(m.ids()), len(m.suggestions), len(m.options)), m.lang.text("Suggestions are optional; every choice is editable."), ""}
	if m.single {
		lines[1] = fmt.Sprintf(m.lang.text("Choose one profile | %d available"), len(m.options))
		lines[2] = m.lang.text("This describes the project; it does not scaffold it.")
	}
	for index := start; index < end; index++ {
		s := m.options[index]
		cursor, checked, tag := " ", " ", ""
		if index == m.cursor {
			cursor = ">"
		}
		if m.selected[s.ID] {
			checked = "x"
		}
		if m.suggestions[s.ID] != "" {
			tag = " *"
		}
		label := s.Label
		if label == "" {
			label = s.ID
		}
		topic := ""
		if len(s.Topics) > 0 {
			topic = " [" + m.lang.topics(s.Topics) + "]"
		}
		box := "[" + checked + "]"
		if m.single {
			box = "( )"
			if checked == "x" {
				box = "(o)"
			}
		}
		lines = append(lines, fmt.Sprintf("%s %s %s%s%s", cursor, box, label, topic, tag), "    "+s.Description)
	}
	reason := m.lang.text("* suggested")
	if len(m.options) > 0 && m.suggestions[m.options[m.cursor].ID] != "" {
		reason = fmt.Sprintf(m.lang.text("* suggested: %s"), m.lang.reason(m.suggestions[m.options[m.cursor].ID]))
	}
	if m.single {
		reason = m.lang.text("Only one architecture profile is selected.")
	}
	lines = append(lines, "", reason, m.lang.text("Up/Down or j/k: move | Space: toggle"), m.lang.text("Enter: confirm | a: all | n: none | r: reset"), m.lang.text("q or Ctrl-C: cancel (no files changed)"))
	if m.single {
		lines[len(lines)-2] = m.lang.text("Enter: choose highlighted profile | q/Ctrl-C: cancel")
	}
	for i, line := range lines {
		lines[i] = clip(line, width)
	}
	return strings.Join(lines, "\r\n")
}
func readPickerKey(reader *bufio.Reader) (string, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	switch b {
	case 3, 'q', 'Q':
		return "cancel", nil
	case '\r':
		if reader.Buffered() > 0 {
			next, _ := reader.Peek(1)
			if next[0] == '\n' {
				reader.ReadByte()
			}
		}
		return "enter", nil
	case '\n':
		return "enter", nil
	case ' ':
		return "space", nil
	case 'j':
		return "down", nil
	case 'k':
		return "up", nil
	case 'a':
		return "all", nil
	case 'n':
		return "none", nil
	case 'r':
		return "reset", nil
	case 27:
		next, e := reader.ReadByte()
		if e != nil {
			return "", e
		}
		if next != '[' && next != 'O' {
			return "", nil
		}
		next, e = reader.ReadByte()
		if e != nil {
			return "", e
		}
		switch next {
		case 'A':
			return "up", nil
		case 'B':
			return "down", nil
		case 'H':
			return "home", nil
		case 'F':
			return "end", nil
		}
	}
	return "", nil
}
func selectSkills(reader *bufio.Reader, input io.Reader, out io.Writer, skills []Skill, suggestions map[string]string) (ids []string, accepted bool, err error) {
	return selectOptions(reader, input, out, newPicker(skills, suggestions))
}

func selectOptions(reader *bufio.Reader, input io.Reader, out io.Writer, m *pickerModel) (ids []string, accepted bool, err error) {
	restore, e := preparePickerSession(input, out)
	if e != nil {
		return nil, false, fmt.Errorf("cannot prepare interactive terminal: %w", e)
	}
	defer func() {
		if e := restore(); e != nil {
			ids = nil
			accepted = false
			err = fmt.Errorf("cannot restore terminal: %w", e)
		}
	}()
	for !m.done {
		width, height := pickerSize(out)
		fmt.Fprint(out, "\x1b[H\x1b[2J", m.view(width, height))
		key, e := readPickerKey(reader)
		if e != nil {
			return nil, false, fmt.Errorf("input ended; no installation started")
		}
		m.apply(key)
	}
	if !m.accepted {
		return nil, false, nil
	}
	return m.ids(), true, nil
}
