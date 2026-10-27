package main

import (
	_ "embed"
	"encoding/json"
)

// The bundle is a presentation fallback, never a discovery inventory.
//
//go:embed skill_descriptions_es.json
var spanishSkillData []byte

var spanishSkills = func() map[string]struct{ Source, ES string } {
	var entries map[string]struct{ Source, ES string }
	if err := json.Unmarshal(spanishSkillData, &entries); err != nil {
		panic(err)
	}
	return entries
}()

func (l language) skillDescription(s Skill) string {
	if !l {
		return s.Description
	}
	if plainText(s.DescriptionES, 1024) {
		return s.DescriptionES
	}
	if entry, ok := spanishSkills[s.ID]; ok && entry.Source == s.Description && plainText(entry.ES, 1024) {
		return entry.ES
	}
	return s.Description
}
