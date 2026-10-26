# Install only the skills your project needs

`zesh-skills` is a standalone, project-local skill installer for Claude Code, Codex, OpenCode and Pi. The executable can live on your PATH; skills never install into your home-level agent configuration.

## Quick path

```text
cd your-project
zesh-skills
```

The wizard confirms the current directory as the project root, detects framework manifests and existing instruction files, asks for the architecture, then opens a keyboard-driven checkbox picker. Every relevant skill remains available with its topic and a brief description; no skill IDs need to be typed.

| Key | Action |
| --- | --- |
| Up/Down or `j`/`k` | Move through the scrolling skill list |
| Space | Check/uncheck the focused skill |
| Enter | Accept the current selection |
| `a` / `n` / `r` | Select all / clear selection / restore suggestions |
| `q` or Ctrl-C | Cancel without installing anything |

A small suggested set is visibly checked and marked `*`. Suggestions include framework-neutral debugging, accessibility for frontend/apps, and matching framework/design foundations when package declarations provide evidence. At most five are suggested. Versioned skills require the same declared package major; unknown ranges and conflicting monorepo versions do not select a versioned profile. Missing evidence never forces .NET 8/10 or a new UI library. All suggestions can be unchecked; security/data-access skills are not mandatory. For a landing-page mockup, keep only the design skills you want. Descriptions and suggestion reasons stay visible, and the screen/cursor/input mode are restored on confirmation, cancellation or input failure. Scope and agent setup still use the short initial questions.

The final preview lists source commit, selected IDs, destination directories and instruction files. Answer `y` to install or `n` to leave the project unchanged. Confirming an empty picker selection also cancels. Architecture accepts the displayed options or your own text; it does not scaffold or alter your application's architecture. Clearing choices is not an uninstall: previously installed selections are retained.

```bash
zesh-skills --scope frontend --list
zesh-skills --help

# Automation: explicit choices; no prompts, including on local conflicts.
zesh-skills --non-interactive --agents pi --architecture mockup \
  --skills frontend/tailwind-4 --yes

# Mixed project; multiple agent targets and an additional instruction file.
zesh-skills --scope frontend,backend --agents claude,codex \
  --architecture "modular monolith" --skills frontend/react-19,backend/general \
  --instructions GEMINI.md --yes
```

`--root` overrides the current directory. `--yes` approves the displayed selection, not overwriting local edits. On a terminal it can still ask about locally modified files; in automation a conflict fails without changing installed data. `--list`, `--help` and `--version` never install anything; help/version do not access the network. Exit codes: `0` success or deliberate cancellation, `1` network/filesystem/conflict failure, `2` invalid arguments or incomplete input.

## Build a dependency-free executable

Development requires Go 1.25+; `golang.org/x/term` and its platform support are downloaded as pinned Go modules on the first build. They are compiled into the executable: the installed binary requires neither Go, Node, Git nor a .NET runtime. Network access is required to install/update skills. The picker supports Windows consoles/Windows Terminal and ANSI terminals on Linux/macOS; redirected input/output uses the non-interactive flag contract instead of changing console modes.

```bash
# Go only; terminal support is compiled into the executable.
CGO_ENABLED=0 go build -trimpath -o dist/zesh-skills ./cli

# Repository convenience scripts use Node only during development.
npm run build:zesh-cli
npm run build:zesh-cli -- --all
```

The convenience builder writes `dist/zesh-skills/<os>-<arch>/zesh-skills[.exe]` and `SHA256SUMS`. `--all` cross-compiles Windows, Linux and macOS for amd64 and arm64 with CGO disabled. Cross-compilation is not runtime testing on those platforms. There is no publishing or automatic global installation step.

To make `zesh-skills` available from any terminal, copy **the executable for your OS/architecture** into a directory already on your PATH, or add its directory to your user PATH. On Linux/macOS, `~/.local/bin` is a conventional choice (create it, copy the binary, make it executable and add it to PATH if needed). On Windows, choose a user-owned tools folder and add it through the user Environment Variables settings; open a new terminal afterward. Do not add a writable project checkout itself to a machine-wide PATH. These are manual user actions, not actions performed by the build or installer.

## Project files and update behavior

| Location | Purpose |
| --- | --- |
| `.zesh-skills/content/skills/<domain>/<name>/` | Selected upstream instructions and their bundled resources |
| `.zesh-skills/manifest.json` | Installed selections, last source commit and content hashes |
| `.claude/skills/<domain>-<name>/SKILL.md` | Claude Code discovery entry |
| `.agents/skills/<domain>-<name>/SKILL.md` | Codex discovery entry |
| `.opencode/skills/<domain>-<name>/SKILL.md` | OpenCode discovery entry |
| `.pi/skills/<domain>-<name>/SKILL.md` | Pi discovery entry |
| `CLAUDE.md`, `AGENTS.md` | Managed installed-skill index, descriptions and usage triggers |

Discovery entries link to canonical local instructions, preserving domain-qualified names such as `frontend-security` and `backend-security`. Selecting several agents creates their native entries; agents that also scan compatibility directories may see equivalent duplicate entries. Original upstream skill names/content remain in canonical storage.

Existing case variants such as `claude.md` are reused, not duplicated. Unrelated Markdown outside `<!-- zesh-skills:start -->` / `<!-- zesh-skills:end -->` is preserved byte-for-byte, including its line endings. Malformed/duplicate markers stop installation. Root-level `GEMINI.md` and `copilot-instructions.md` are detected and can be selected as an additional instruction file; arbitrary Markdown files are not guessed or rewritten. Custom nested agent configuration is not modified.

Run the wizard again to add or update a selection. Previously installed skills and agent targets are retained; omitting them is not an uninstall operation. Unmodified files update automatically; local edits and unowned colliding files require explicit confirmation. Declining any conflict cancels the entire write phase. Removing skills/agents is intentionally not implemented in this first version.

Downloads resolve `main` **once** through the public GitHub API. The catalog and every selected file come from that exact commit, with catalog-provided SHA-256 validation. Files are fetched before writing. References to unselected skills or repository-only material become commit-pinned GitHub links rather than implicit skill installations. Shared supporting references required by a selected skill remain local. Remote links require connectivity when opened; a completely offline dependency closure is not promised.

Paths are validated for portability, symbolic-link destinations are refused, and rooted filesystem operations confine writes. A project lock prevents simultaneous installer runs. Ordinary write failures trigger rollback; power loss/process termination is not an atomic transaction. If interrupted, inspect `.zesh-skills.lock` and installed files before manually removing a stale lock and rerunning. Do not remove a lock owned by a live process. No project scripts or downloaded executables run during installation.

## Maintain the remote catalog

Installations read the published `skills-index.json` from the repository's `main`. Future skill changes must publish the regenerated catalog alongside their source changes. Local tests use isolated HTTP fixtures; building a binary alone does not publish skills or a release.

```bash
npm run catalog:generate -- --help  # Side-effect-free catalog command help
npm run catalog:generate
npm run catalog:check
npm run test:catalog
npm test                 # Existing Node-only suite, including catalog freshness
npm run test:zesh-cli     # Separate executable Go acceptance/unit checks
```

`skills-index.json` uses `schemaVersion: 1` and stable `<domain>/<directory>` IDs. Entries contain name, description, scopes, topics, technologies, triggers, source path, version and file hashes. The generator reads all `SKILL.md` YAML frontmatter, including folded descriptions. Add a skill, regenerate the catalog and submit both together; no binary rebuild is needed. The pipeline rejects a stale catalog. A future incompatible catalog schema requires a compatible CLI update.

Optional frontmatter overrides allow precise classification without editing the generator:

```yaml
metadata:
  version: "1.0"
  catalog:
    scopes: [frontend]
    topics: [design]
    technologies: [react]
    triggers: ["Landing-page layout and visual composition"]
```

Without overrides, existing domain metadata and conservative topic defaults are used. Scopes are `frontend`, `backend`, `app`; topics are descriptive and extensible. Shared skills may apply to several scopes. Detection suggests scopes from known package manifests and .NET project files without running tools; it does not hide skills based on an inferred framework version. Review version compatibility when selecting .NET 8 versus .NET 10 or framework-specific skills.

## Integration references and proof limits

Discovery paths follow the documented [Claude Code](https://code.claude.com/docs/en/skills), [Codex](https://developers.openai.com/codex/skills), [OpenCode](https://opencode.ai/docs/skills/) and [Pi](https://github.com/badlogic/pi-mono/blob/main/packages/coding-agent/docs/skills.md) conventions. Generated layouts are checked locally; actual model loading/behavior and every host's project-trust UI require consumer-side verification and are not established by these tests.
