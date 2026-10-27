package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
)

// CLIUserPath is a user-only PATH store. Kind 0 represents an absent value;
// writing kind 0 restores absence. Implementations must not modify machine PATH.
type CLIUserPath interface {
	Read() (string, uint32, error)
	Write(string, uint32) error
}

// CLIInstallService has no ambient environment or registry fallback. Callers must
// supply all host dependencies explicitly, so tests cannot reach the real PATH.
type CLIInstallService struct {
	OS, Home, Executable string
	Env                  map[string]string
	WindowsPath          CLIUserPath
}

type CLIInstallConsent struct{ Install, Replace bool }

// CLIInstallPlan exposes preview data, while retaining private snapshots used by
// Apply. Changing a preview field cannot redirect a write.
type CLIInstallPlan struct {
	Destination           string
	ConfigFiles           []string
	Replaces              bool
	PATHChanged           bool
	service               CLIInstallService
	source                installFile
	files                 []installChange
	pathBefore, pathAfter string
	pathKind              uint32
}

type installFile struct {
	path   string
	data   []byte
	mode   os.FileMode
	info   os.FileInfo
	exists bool
}
type installChange struct {
	before installFile
	after  []byte
	mode   os.FileMode
}

// Plan validates and snapshots the complete operation without creating anything.
func (s CLIInstallService) Plan() (CLIInstallPlan, error) {
	p := CLIInstallPlan{service: s}
	p.service.Env = make(map[string]string, len(s.Env))
	for k, v := range s.Env {
		p.service.Env[k] = v
	}
	if s.OS != "windows" && s.OS != "linux" && s.OS != "darwin" {
		return p, errors.New("unsupported operating system")
	}
	if err := installSafePath(s.Home, s.OS); err != nil {
		return p, err
	}
	source, err := installSnapshot(s.Executable)
	if err != nil {
		return p, err
	}
	if !source.exists {
		return p, errors.New("source executable is missing")
	}
	p.source = source
	var bin string
	var configs []string
	if s.OS == "windows" {
		local := s.Env["LOCALAPPDATA"]
		if err := installSafePath(local, s.OS); err != nil {
			return p, err
		}
		bin = filepath.Join(local, "Programs", "zesh-skills")
		p.Destination = filepath.Join(bin, "zesh-skills.exe")
		if s.WindowsPath == nil {
			return p, errors.New("user PATH store is unavailable")
		}
		p.pathBefore, p.pathKind, err = s.WindowsPath.Read()
		if err != nil {
			return p, err
		}
		if p.pathKind != 0 && p.pathKind != 1 && p.pathKind != 2 {
			return p, errors.New("user PATH is not a string")
		}
		if p.pathKind == 0 && p.pathBefore != "" {
			return p, errors.New("absent user PATH has content")
		}
		if strings.ContainsAny(p.pathBefore, "\x00\r\n") {
			return p, errors.New("unsafe user PATH")
		}
		p.pathAfter = p.pathBefore
		present := false
		for _, entry := range strings.Split(p.pathBefore, ";") {
			// Do not expand variables or rewrite any pre-existing entry.
			if strings.EqualFold(strings.TrimRight(strings.Trim(entry, " \""), `/\`), strings.TrimRight(bin, `/\`)) {
				present = true
			}
		}
		if !present {
			if p.pathAfter != "" && !strings.HasSuffix(p.pathAfter, ";") {
				p.pathAfter += ";"
			}
			p.pathAfter += bin
		}
		p.PATHChanged = p.pathAfter != p.pathBefore
	} else {
		bin = filepath.Join(s.Home, ".local", "bin")
		p.Destination = filepath.Join(bin, "zesh-skills")
		shell := filepath.Base(s.Env["SHELL"])
		switch shell {
		case "sh":
			configs = []string{filepath.Join(s.Home, ".profile")}
		case "bash":
			login := filepath.Join(s.Home, ".profile")
			for _, name := range []string{".bash_profile", ".bash_login"} {
				path := filepath.Join(s.Home, name)
				if _, err := os.Lstat(path); err == nil {
					login = path
					break
				} else if !os.IsNotExist(err) {
					return p, err
				}
			}
			configs = []string{filepath.Join(s.Home, ".bashrc"), login}
		case "zsh":
			root := s.Home
			if value := s.Env["ZDOTDIR"]; value != "" {
				root = value
			}
			if err := installConfigRoot(s.Home, root, s.OS); err != nil {
				return p, err
			}
			configs = []string{filepath.Join(root, ".zshrc"), filepath.Join(root, ".zprofile")}
		case "fish":
			root := filepath.Join(s.Home, ".config")
			if value := s.Env["XDG_CONFIG_HOME"]; value != "" {
				root = value
			}
			if err := installConfigRoot(s.Home, root, s.OS); err != nil {
				return p, err
			}
			configs = []string{filepath.Join(root, "fish", "conf.d", "zesh-skills.fish")}
		default:
			return p, errors.New("unsupported shell; expected sh, bash, zsh or fish in SHELL")
		}
		block := installShellBlock(bin, shell)
		for _, path := range configs {
			before, err := installSnapshot(path)
			if err != nil {
				return p, err
			}
			content := string(before.data)
			if strings.Contains(content, "# >>> zesh-skills PATH >>>") || strings.Contains(content, "# <<< zesh-skills PATH <<<") {
				if strings.Count(content, block) != 1 || strings.Count(content, "# >>> zesh-skills PATH >>>") != 1 || strings.Count(content, "# <<< zesh-skills PATH <<<") != 1 {
					return p, fmt.Errorf("modified PATH block: %s", path)
				}
			} else {
				if content != "" && !strings.HasSuffix(content, "\n") {
					content += "\n"
				}
				content += block
			}
			mode := os.FileMode(0600)
			if before.exists {
				mode = before.mode
			}
			p.files = append(p.files, installChange{before, []byte(content), mode})
			p.PATHChanged = p.PATHChanged || !bytes.Equal(before.data, []byte(content))
		}
	}
	if err := installSafePath(bin, s.OS); err != nil {
		return p, err
	}
	target, err := installSnapshot(p.Destination)
	if err != nil {
		return p, err
	}
	p.Replaces = target.exists && !bytes.Equal(target.data, source.data)
	// Binary precedes configuration so a failed config write can restore it.
	p.files = append([]installChange{{target, source.data, 0755}}, p.files...)
	p.ConfigFiles = append([]string(nil), configs...)
	return p, nil
}

// Apply requires a fresh, matching plan and explicit replacement consent. A
// cooperative per-home lock serializes installers. It detects drift after the
// prompt and before each write; it is not a sandbox against a malicious user
// concurrently replacing their own home directory.
func (s CLIInstallService) Apply(p CLIInstallPlan, consent CLIInstallConsent) (result error) {
	if !consent.Install {
		return errors.New("installation was not confirmed")
	}
	if p.source.path == "" || s.OS != p.service.OS || s.Home != p.service.Home || s.Executable != p.service.Executable || !reflect.DeepEqual(s.Env, p.service.Env) || s.WindowsPath != p.service.WindowsPath {
		return errors.New("installation plan does not match service")
	}
	// Preview fields are caller-owned display data, never authorization.
	if len(p.files) == 0 {
		return errors.New("installation plan has no destination")
	}
	target := p.files[0].before
	if target.exists && !bytes.Equal(target.data, p.source.data) && !consent.Replace {
		return errors.New("existing executable differs; explicit replacement confirmation is required")
	}
	if err := installUnchanged(p.source); err != nil {
		return err
	}
	for _, change := range p.files {
		if err := installUnchanged(change.before); err != nil {
			return err
		}
	}
	if err := s.checkPath(p); err != nil {
		return err
	}
	lock := filepath.Join(s.Home, ".zesh-skills-install.lock")
	if err := installNoLinks(lock); err != nil {
		return err
	}
	f, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("cannot acquire installation lock: %w", err)
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(lock)
		return err
	}
	defer func() { result = errors.Join(result, os.Remove(lock)) }()
	var applied []installChange
	var dirs []string
	defer func() {
		if result == nil {
			return
		}
		for i := len(applied) - 1; i >= 0; i-- {
			change := applied[i]
			current, err := installSnapshot(change.before.path)
			if err == nil && (!current.exists || !bytes.Equal(current.data, change.after)) {
				err = errors.New("rollback refused: file changed after installation")
			}
			if err == nil {
				if change.before.exists {
					err = installAtomicWrite(change.before.path, change.before.data, change.before.mode)
				} else {
					err = os.Remove(change.before.path)
				}
			}
			result = errors.Join(result, err)
		}
		for i := len(dirs) - 1; i >= 0; i-- {
			result = errors.Join(result, os.Remove(dirs[i]))
		}
	}()
	if err := installUnchanged(p.source); err != nil {
		return err
	}
	for _, change := range p.files {
		if err := installUnchanged(change.before); err != nil {
			return err
		}
		if change.before.exists && bytes.Equal(change.before.data, change.after) && (s.OS == "windows" || change.before.mode.Perm() == change.mode.Perm()) {
			continue
		}
		if err := installMakeDirs(filepath.Dir(change.before.path), &dirs); err != nil {
			return err
		}
		if err := installUnchanged(change.before); err != nil {
			return err
		}
		if err := installAtomicWrite(change.before.path, change.after, change.mode); err != nil {
			return err
		}
		applied = append(applied, change)
	}
	if err := s.checkPath(p); err != nil {
		return err
	}
	if s.OS == "windows" && p.pathAfter != p.pathBefore {
		kind := p.pathKind
		if kind == 0 {
			kind = 2
		}
		if err := s.WindowsPath.Write(p.pathAfter, kind); err != nil {
			// A store may fail after writing. Restore only our exact value; never
			// overwrite a concurrent unrelated edit while rolling back.
			value, currentKind, readErr := s.WindowsPath.Read()
			if readErr != nil {
				return errors.Join(err, readErr)
			}
			if value == p.pathAfter && currentKind == kind {
				return errors.Join(err, s.WindowsPath.Write(p.pathBefore, p.pathKind))
			}
			return err
		}
	}
	return nil
}

func (s CLIInstallService) checkPath(p CLIInstallPlan) error {
	if s.OS != "windows" {
		return nil
	}
	value, kind, err := s.WindowsPath.Read()
	if err != nil {
		return err
	}
	if value != p.pathBefore || kind != p.pathKind {
		return errors.New("user PATH changed after preview")
	}
	return nil
}

func installSafePath(path, platform string) error {
	if path == "" || !filepath.IsAbs(path) || strings.IndexFunc(path, unicode.IsControl) >= 0 {
		return errors.New("unsafe or relative installation path")
	}
	forbidden := ":"
	if platform == "windows" {
		forbidden = ";%\""
		if strings.HasPrefix(path, `\\`) {
			return errors.New("network installation paths are unsupported")
		}
	}
	if strings.ContainsAny(path, forbidden) && !(platform != "windows" && filepath.VolumeName(path) != "" && !strings.ContainsAny(strings.TrimPrefix(path, filepath.VolumeName(path)), forbidden)) {
		return errors.New("installation path contains a PATH delimiter")
	}
	return installNoLinks(path)
}

func installConfigRoot(home, root, platform string) error {
	if err := installSafePath(root, platform); err != nil {
		return err
	}
	rel, err := filepath.Rel(home, root)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return errors.New("shell configuration must be inside the user home")
	}
	return nil
}

func installNoLinks(path string) error {
	if !filepath.IsAbs(path) {
		return errors.New("expected an absolute file path")
	}
	for current := filepath.Clean(path); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symbolic link is not allowed: %s", current)
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return nil
}

func installSnapshot(path string) (installFile, error) {
	state := installFile{path: path}
	if err := installNoLinks(path); err != nil {
		return state, err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if !info.Mode().IsRegular() {
		return state, fmt.Errorf("not a regular file: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return state, err
	}
	state.data, state.mode, state.info, state.exists = data, info.Mode(), info, true
	return state, nil
}

func installUnchanged(before installFile) error {
	now, err := installSnapshot(before.path)
	if err != nil {
		return err
	}
	if now.exists != before.exists || !bytes.Equal(now.data, before.data) || now.mode != before.mode || (now.exists && (!os.SameFile(now.info, before.info) || !now.info.ModTime().Equal(before.info.ModTime()))) {
		return fmt.Errorf("file changed after preview: %s", before.path)
	}
	return nil
}

func installMakeDirs(path string, created *[]string) error {
	if err := installNoLinks(path); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return errors.New("installation parent is not a directory")
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	if err := installMakeDirs(filepath.Dir(path), created); err != nil {
		return err
	}
	if err := os.Mkdir(path, 0700); err != nil {
		return err
	}
	*created = append(*created, path)
	return nil
}

func installAtomicWrite(path string, data []byte, mode os.FileMode) (result error) {
	if err := installNoLinks(path); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".zesh-install-*")
	if err != nil {
		return err
	}
	temp := file.Name()
	defer func() {
		if result != nil {
			_ = os.Remove(temp)
		}
	}()
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Chmod(mode.Perm()); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temp, path)
}

func installShellBlock(bin, shell string) string {
	// filepath.ToSlash also lets shell-plan tests run on Windows fixtures.
	bin = filepath.ToSlash(bin)
	quoted := "'" + strings.ReplaceAll(bin, "'", "'\"'\"'") + "'"
	body := "case \":$PATH:\" in\n  *:" + quoted + ":*) ;;\n  *) export PATH=" + quoted + ":\"$PATH\" ;;\nesac\n"
	if shell == "fish" {
		quoted = "'" + strings.ReplaceAll(strings.ReplaceAll(bin, `\`, `\\`), "'", `\'`) + "'"
		body = "if not contains -- " + quoted + " $PATH\n    set -gx PATH " + quoted + " $PATH\nend\n"
	}
	return "# >>> zesh-skills PATH >>>\n" + body + "# <<< zesh-skills PATH <<<\n"
}
