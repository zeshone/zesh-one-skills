# Install only the skills your project needs

`zesh-skills` is a standalone, project-local skill installer for Claude Code, Codex, OpenCode and Pi. The executable can live on your PATH; skills never install into your home-level agent configuration.

## Download v2.8.0 (publication pending)

The public catalog is already on `main`. The CLI release is prepared locally, **not yet published**; the following canonical links work only after maintainers publish [v2.8.0](https://github.com/zeshone/zesh-one-skills/releases/tag/v2.8.0). Do not confuse a development build with a released binary.

Download one raw executable and [SHA256SUMS](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/SHA256SUMS) from that release:

| Platform | Asset |
| --- | --- |
| Windows x64 | [zesh-skills-windows-amd64.exe](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-windows-amd64.exe) |
| Windows ARM64 | [zesh-skills-windows-arm64.exe](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-windows-arm64.exe) |
| Linux x64 | [zesh-skills-linux-amd64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-linux-amd64) |
| Linux ARM64 | [zesh-skills-linux-arm64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-linux-arm64) |
| macOS Intel | [zesh-skills-darwin-amd64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-darwin-amd64) |
| macOS Apple Silicon | [zesh-skills-darwin-arm64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-darwin-arm64) |

Linux example (replace the filename for your platform):

```bash
base=https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0
curl -fLO "$base/zesh-skills-linux-amd64"
curl -fLO "$base/SHA256SUMS"
grep '  zesh-skills-linux-amd64$' SHA256SUMS | sha256sum --check -
# Continue only if the checksum matches.
chmod +x zesh-skills-linux-amd64
./zesh-skills-linux-amd64 --version  # expected: 2.8.0
```

On macOS use `shasum -a 256 -c` instead of `sha256sum --check -`, with the selected macOS filename. Check only the downloaded asset's line; the manifest also lists the other five platforms.

Windows PowerShell example:

```powershell
$base = 'https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0'
$name = 'zesh-skills-windows-amd64.exe'
Invoke-WebRequest "$base/$name" -OutFile $name
Invoke-WebRequest "$base/SHA256SUMS" -OutFile SHA256SUMS
$line = Get-Content SHA256SUMS | Where-Object { $_.EndsWith("  $name") }
if (@($line).Count -ne 1) { throw 'Missing or duplicate checksum' }
$expected = ($line -split '  ')[0]
if ((Get-FileHash $name -Algorithm SHA256).Hash.ToLowerInvariant() -ne $expected) { throw 'Checksum mismatch' }
& ".\$name" --version  # expected: 2.8.0
```

Checksums detect mismatched bytes, not an independent signature of publisher identity. After verification, optionally rename the executable to `zesh-skills` (`zesh-skills.exe` on Windows) and follow the manual PATH guidance below. No Go, Node, Git or .NET runtime is needed to run it. Skills remain project-local. macOS and ARM artifacts have compile-only evidence so far; the CI matrix is configured, not an observed online runtime pass.

## Quick path

```text
cd your-project
zesh-skills
# Spanish interface; English remains the default without this flag.
zesh-skills --es
```

### Interface language

Use `--es` to localize prompts, scope/architecture choices and descriptions, picker controls, previews, confirmations and app diagnostics into Spanish. There is no extra startup question or saved language preference. Spanish confirmations accept `s`, `si`, `sí`, `y` or `yes`; English behavior remains unchanged. `--es=false` explicitly keeps English.

Architecture names stay recognizable in Spanish: they remain in English or accompany their translation, for example `Arquitectura Limpia (Clean Architecture)`. Hybrid names such as `CQRS + event sourcing` stay in English to keep the menu compact. Explanations remain Spanish; familiar terms such as UI, view models, focus, feedback and triggers retain their English names.

In the skill catalog (`--list` and the interactive skill picker), **only descriptions** are localized: skill names/IDs and topic headings/tags keep their original English. With `--es`, a validated remote `descriptionEs` takes precedence. The binary also bundles 35 Spanish translations for the original catalog, used only when both ID and exact original English description match. Unknown skills or changed descriptions without a translation honestly fall back to English; they remain discoverable. Without `--es` (or with `--es=false`), descriptions stay exactly as supplied in English.

The language flag does not translate flags, JSON formats, custom user values, downloaded skill content or generated agent documents. Localized descriptions are presentation-only: the manifest and installed instruction tables retain canonical original metadata, without `descriptionEs`. Technical OS details may retain their original language and are identified as such. Predefined architecture profiles are stored by stable ID, regardless of the displayed language.

The wizard confirms the current directory as the project root, detects framework manifests and existing instruction files, asks for the architecture, then opens a keyboard-driven checkbox picker. Every relevant skill remains available with its topic and a brief description; no skill IDs need to be typed.

| Key | Action |
| --- | --- |
| Up/Down or `j`/`k` | Move through the scrolling skill list |
| Space | Check/uncheck the focused skill |
| Enter | Accept the current selection |
| `a` / `n` / `r` | Select all / clear selection / restore suggestions |
| `q` or Ctrl-C | Cancel without installing anything |

A small suggested set is visibly checked and marked `*`. Suggestions include framework-neutral debugging, accessibility for frontend/apps, and matching framework/design foundations when package declarations provide evidence. At most five are suggested. Versioned skills require the same declared package major; unknown ranges and conflicting monorepo versions do not select a versioned profile. Missing evidence never forces .NET 8/10 or a new UI library. All suggestions can be unchecked; security/data-access skills are not mandatory. For a landing-page mockup, keep only the design skills you want. Descriptions and suggestion reasons stay visible, and the screen/cursor/input mode are restored on confirmation, cancellation, input failure or catchable SIGTERM termination. Forced process kills cannot run terminal cleanup. Scope uses the same checkbox menu and can combine frontend, backend and app. Architecture uses an exclusive selection menu; agent setup retains its short question.

### Architecture profiles

The architecture menu includes 20 explained choices: minimal monolith, layered/N-tier, vertical slices, Clean, Hexagonal, Onion, modular monolith, microservices, event-driven, serverless functions, MVC, MVVM, component-based UI and feature-based UI. Common hybrids include modular monolith + vertical slices, Clean + vertical slices, Hexagonal + DDD, microservices + events, and CQRS + event sourcing. **Custom architecture** accepts your own text.

Each entry explains its scope and tradeoffs. Deployment styles and pattern combinations are identified as such, not presented as interchangeable domain designs. A built-in selection stores its stable profile ID; custom text is preserved. Selection does not scaffold code, add infrastructure or impose a skill bundle. Profiles are maintained centrally in [`profiles.go`](profiles.go).

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
# Maintainer release preparation; explicit version, no publication:
npm run build:zesh-cli -- --release 2.8.0
```

The convenience builder writes `dist/zesh-skills/<os>-<arch>/zesh-skills[.exe]` and `SHA256SUMS`. `--all` cross-compiles Windows, Linux and macOS for amd64 and arm64 with CGO disabled. Cross-compilation is not runtime testing on those platforms. Default and `--all` builds retain the development version `dev`. `--release <version>` accepts only a plain numeric `major.minor.patch` (no `v` prefix, prerelease or build suffix); it builds all six targets with `CGO_ENABLED=0`, `GOTOOLCHAIN=local` and linker `main.version` set to that one argument. For `2.8.0`, output is `dist/zesh-skills-release/v2.8.0/`: the six flat filenames above plus `SHA256SUMS` (lowercase SHA-256, two spaces, actual asset filename, newline). Only those explicit seven files belong in an upload; unrelated/stale files are neither included nor deleted. A failed build is not a publishable release. Build from the final authorized main commit and verify before publication; this command never creates tags or overwrites remote releases. There is no publishing or automatic global installation step.

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

`skills-index.json` uses `schemaVersion: 1` and stable `<domain>/<directory>` IDs. Entries contain name, description, scopes, topics, technologies, triggers, source path, version and file hashes. Optional `descriptionEs` is a bounded plain-text Spanish description (up to 1024 UTF-8 bytes, no control characters); missing or empty values remain valid legacy input. The generator reads all `SKILL.md` YAML frontmatter, including folded descriptions. Add a skill, regenerate the catalog and submit both together; no binary rebuild is needed. The pipeline rejects a stale catalog. A future incompatible catalog schema requires a compatible CLI update.

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

For a new skill, optional `metadata.catalog.descriptionEs` supplies the Spanish UI description without a binary rebuild. The generator otherwise uses `cli/skill_descriptions_es.json` only for an exact ID/source match; that file is optional in legacy fixtures and never drives discovery. Keep the canonical English `description` unchanged. Invalid authored translations are rejected before catalog generation.

Without overrides, existing domain metadata and conservative topic defaults are used. Scopes are `frontend`, `backend`, `app`; topics are descriptive and extensible. Shared skills may apply to several scopes. Detection suggests scopes from known package manifests and .NET project files without running tools; it does not hide skills based on an inferred framework version. Review version compatibility when selecting .NET 8 versus .NET 10 or framework-specific skills.

## Integration references and proof limits

Discovery paths follow the documented [Claude Code](https://code.claude.com/docs/en/skills), [Codex](https://developers.openai.com/codex/skills), [OpenCode](https://opencode.ai/docs/skills/) and [Pi](https://github.com/badlogic/pi-mono/blob/main/packages/coding-agent/docs/skills.md) conventions. Generated layouts are checked locally; actual model loading/behavior and every host's project-trust UI require consumer-side verification and are not established by these tests.
