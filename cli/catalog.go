package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

type File struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Skill struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	DescriptionES string   `json:"descriptionEs,omitempty"`
	Path          string   `json:"path"`
	Version       string   `json:"version"`
	Scopes        []string `json:"scopes"`
	Topics        []string `json:"topics"`
	Technologies  []string `json:"technologies"`
	Triggers      []string `json:"triggers"`
	Files         []File   `json:"files"`
}

// Reject explicit null as well as non-string optional translations at the JSON boundary.
func (s *Skill) UnmarshalJSON(data []byte) error {
	type skillJSON Skill
	var decoded struct {
		*skillJSON
		DescriptionES json.RawMessage `json:"descriptionEs"`
	}
	decoded.skillJSON = (*skillJSON)(s)
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	if len(decoded.DescriptionES) > 0 {
		if string(decoded.DescriptionES) == "null" {
			return fmt.Errorf("invalid Spanish description type")
		}
		return json.Unmarshal(decoded.DescriptionES, &s.DescriptionES)
	}
	return nil
}

type Catalog struct {
	SchemaVersion int     `json:"schemaVersion"`
	Skills        []Skill `json:"skills"`
}

var identifier = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var digest = regexp.MustCompile(`^[a-f0-9]{64}$`)
var commitID = regexp.MustCompile(`^[a-f0-9]{40}$`)

func safePath(name string) bool {
	if name == "" || path.Clean(name) != name || strings.ContainsAny(name, "\\:\x00") || strings.HasPrefix(name, "/") {
		return false
	}
	for _, part := range strings.Split(name, "/") {
		if part == "" || strings.HasPrefix(part, ".") || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") {
			return false
		}
		for _, r := range part {
			if r < 32 || strings.ContainsRune(`<>"|?*`, r) {
				return false
			}
		}
		base := strings.ToUpper(strings.Split(part, ".")[0])
		if base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" || regexp.MustCompile(`^(COM|LPT)[0-9]$`).MatchString(base) {
			return false
		}
	}
	return true
}
func (c Catalog) Validate() error {
	if c.SchemaVersion != 1 {
		return fmt.Errorf("unsupported catalog schema %d; update zesh-skills", c.SchemaVersion)
	}
	if len(c.Skills) == 0 || len(c.Skills) > 1000 {
		return fmt.Errorf("invalid skill count")
	}
	ids := map[string]bool{}
	aliases := map[string]bool{}
	allFiles := map[string]string{}
	for _, s := range c.Skills {
		parts := strings.Split(s.ID, "/")
		alias := strings.ReplaceAll(s.ID, "/", "-")
		if len(parts) != 2 || !identifier.MatchString(parts[0]) || !identifier.MatchString(parts[1]) || len(alias) > 64 || ids[s.ID] || aliases[alias] {
			return fmt.Errorf("invalid or duplicate skill ID: %q", s.ID)
		}
		ids[s.ID] = true
		aliases[alias] = true
		if s.Path != "skills/"+s.ID || !safePath(s.Path) || s.Name == "" || !plainText(s.Description, 1024) || (s.DescriptionES != "" && !plainText(s.DescriptionES, 1024)) || s.Version == "" || len(s.Scopes) == 0 || len(s.Topics) == 0 || len(s.Triggers) == 0 || len(s.Files) == 0 || len(s.Files) > 500 {
			return fmt.Errorf("incomplete skill: %s", s.ID)
		}
		for _, values := range [][]string{s.Scopes, s.Topics, s.Technologies, s.Triggers} {
			for _, value := range values {
				if !plainText(value, 2048) {
					return fmt.Errorf("invalid metadata: %s", s.ID)
				}
			}
		}
		for _, scope := range s.Scopes {
			if scope != "frontend" && scope != "backend" && scope != "app" {
				return fmt.Errorf("invalid scope: %s", scope)
			}
		}
		found := false
		seen := map[string]bool{}
		for _, f := range s.Files {
			allowed := strings.HasPrefix(f.Path, s.Path+"/") || (strings.HasPrefix(f.Path, "skills/shared/") && strings.HasSuffix(f.Path, ".md") && !strings.HasSuffix(f.Path, "/SKILL.md"))
			if !safePath(f.Path) || !allowed || !digest.MatchString(f.SHA256) || seen[strings.ToLower(f.Path)] {
				return fmt.Errorf("unsafe catalog file: %q", f.Path)
			}
			if old, ok := allFiles[strings.ToLower(f.Path)]; ok && old != f.SHA256 {
				return fmt.Errorf("conflicting file hashes")
			}
			allFiles[strings.ToLower(f.Path)] = f.SHA256
			seen[strings.ToLower(f.Path)] = true
			found = found || f.Path == s.Path+"/SKILL.md"
		}
		if !found {
			return fmt.Errorf("skill entry missing: %s", s.ID)
		}
	}
	return nil
}
func plainText(s string, max int) bool {
	if len(s) == 0 || len(s) > max {
		return false
	}
	for _, r := range s {
		if r < 32 || r == 127 {
			return false
		}
	}
	return true
}
func hash(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }

// Endpoints are fixed in production. Tests substitute isolated HTTP servers.
type Source struct {
	Client   *http.Client
	API, Raw string
}

func github() Source {
	return Source{Client: &http.Client{Timeout: 25 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return fmt.Errorf("redirect refused") }}, API: "https://api.github.com/repos/zeshone/zesh-one-skills", Raw: "https://raw.githubusercontent.com/zeshone/zesh-one-skills"}
}
func (s Source) get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid download request")
	}
	req.Header.Set("User-Agent", "zesh-skills")
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download failed: check network connectivity")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("repository download returned HTTP %d (main may not contain the catalog yet)", res.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, 8*1024*1024+1))
	if err != nil {
		return nil, fmt.Errorf("download interrupted")
	}
	if len(data) > 8*1024*1024 {
		return nil, fmt.Errorf("download exceeds 8 MiB")
	}
	return data, nil
}
func (s Source) Load() (Catalog, string, error) {
	var c Catalog
	data, err := s.get(s.API + "/commits/main")
	if err != nil {
		return c, "", err
	}
	var head struct {
		SHA string `json:"sha"`
	}
	if json.Unmarshal(data, &head) != nil || !commitID.MatchString(head.SHA) {
		return c, "", fmt.Errorf("invalid main commit response")
	}
	data, err = s.Fetch(head.SHA, "skills-index.json")
	if err != nil {
		return c, "", err
	}
	if json.Unmarshal(data, &c) != nil {
		return c, "", fmt.Errorf("invalid catalog JSON")
	}
	return c, head.SHA, c.Validate()
}
func (s Source) Fetch(commit, name string) ([]byte, error) {
	if !commitID.MatchString(commit) || !safePath(name) {
		return nil, fmt.Errorf("unsafe download path")
	}
	return s.get(s.Raw + "/" + commit + "/" + name)
}
