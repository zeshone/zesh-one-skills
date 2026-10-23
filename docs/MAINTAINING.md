# Maintaining Zesh-One-Skills

## 1. Scope and authority
- This is the repository maintainer workflow; consumer setup belongs in [ADOPTING](../ADOPTING.md).
- Use the [skills review rubric](skills-review-rubric.md) for spec, editorial quality and final review decisions; do not substitute a lint pass.
- This guide takes precedence over the vendor template in [skill-creator](../.github/skills/skill-creator/SKILL.md) for repository maintenance. Leave that vendor file attributed, not silently rewritten.

## 2. Select the owner
- Find the existing owner in `skills/` and the applicable domain guide before authoring; prefer a targeted edit over a new skill.
- Keep shared architecture, security, testing, diagnosis and accessibility rules with their existing owners; reference rather than duplicate them.
- Add a skill only for a coherent, reusable, hard-to-infer procedure not already owned; use the rubric to judge scope.

## 3. Preserve applicability gates
- Inspect installed versions, runtime/profile, application capabilities, existing architecture and recorded decisions first.
- Preserve SPA versus Next boundaries, separate .NET 8/.NET 10 conventions, valid Angular forms/HttpClient and Ionic lifecycle ownership.
- Retain the [architecture selection gate](../skills/shared/architecture-selection.md): minimum sufficient structure, alternatives and user choice before new structural changes; reuse an applicable decision.
- Keep the [security baseline](../skills/shared/security/baseline.md) active; route application/delivery checks only to changed boundaries. Do not infer approval for migrations, dependency installation or native sync from a UI task.

## 4. First-party frontmatter
Adapt this shape to the actual owner; follow the rubric for naming/directory alignment and field semantics:

```yaml
---
name: example-skill
description: "Trigger: specific work. State the bounded procedure."
license: Apache-2.0
allowed-tools: Read Write Edit Bash
metadata:
  author: Zesh-One
  version: "1.0"
---
```

- Keep `description` specific with a useful exact `Trigger:`. Quote `metadata.version`.
- Use a single space-separated string for `allowed-tools`, not a YAML list. Declare only applicable tools; declarations are not permission or isolation.
- Verified vendor differences: first-party author is `Zesh-One`, not the template's `gentleman-programming`; do not overwrite vendor authorship.
- The template says always `Apache-2.0`, but existing shared skills declare `MIT`; preserve the actual license/provenance rather than applying a blanket license default.
- The vendor's tools use a YAML list; its generic `skills/{name}/`, single `AGENTS.md`, mandatory `Commands` and local-only references are not repository-wide requirements. Use domain routing, operational sections and relevant local/official external resources instead.

## 5. Keep the operational core small
- English-only applies to skills and domain AGENTS files, not every repository document.
- The established seven-section convention is When to Use, Critical Patterns, Constraints & Tradeoffs, Anti-Patterns, Progressive Disclosure, Resources and Changelog; adapt proportionally, not with empty boilerplate.
- Keep defaults, real gotchas and escape hatches in the core; move detailed support material to references with explicit read conditions.
- Rubric length thresholds are advisory: ideal 80–140 lines, acceptable roughly 180–220; over 220 risks noise and over 500 is a strong smell, not automatic spec failure.

## 6. Maintain conditional routes
- Update the applicable `AGENTS.*.md`: [backend](../AGENTS.backend.md), [.NET 10](../AGENTS.dotnet-10.md), [frontend](../AGENTS.frontend.md), [apps](../AGENTS.apps.md).
- Add conditional work-type routes, not blanket loading of inapplicable frameworks or tools; preserve existing route/profile/security ownership.
- Keep references relative to their containing file and resolvable when copied or relocated; include `shared/security/` with the domain and adjust external installation paths as [ADOPTING](../ADOPTING.md) describes.
- A link is not automatic loading, registry discovery or technical isolation; the host must actually load instructions.

## 7. Write instruction contracts first
- Follow [test-accessibility-instructions.js](../tools/test-accessibility-instructions.js) and existing `tools/test-*-instructions.js`: Node built-ins, `node:test`, clear assertions, no consumer prerequisites.
- Pin positive instructions, negative forbidden mandates, preservation of existing profiles and conditional route/link checks.
- Run `node tools/test-<unit>-instructions.js` before editing instructions: record RED, then implement and record GREEN. Do not weaken a failed assertion merely to pass.
- Trace referenced owners manually too; text-presence checks can miss downstream contradictions.

## 8. Register and run checks
- In [package.json](../package.json), register the focused script and add it exactly once to both `test` and `verify`; retain existing checks.
- `tools/test-cli-instructions.js` pins the non-instruction baseline checks in pipeline order; instruction-script parity is checked dynamically by `tools/test-maintainer-instructions.js`.
- Run the focused contract, then `npm run verify`: fixtures, logging semantics, instruction contracts, skill lint, release readiness and the zero-skills guard.
- Run `git diff --check`; inspect changed files and local links. Keep pipeline membership aligned as scripts are added.

## 9. Keep executable trials separate
- Default `test`/`verify` are Node-only; they do not execute a consumer application or require browser/.NET trial prerequisites.
- Existing opt-in trials are not part of default verification: [test:cli](../skills/shared/cli-design/assets/line-tally/README.md), [test:frontend-dom](../tools/fixtures/frontend-dom/README.md), [test:dotnet-unit](../tools/fixtures/dotnet-testing/README.md).
- Follow each trial's isolated prerequisites, synthetic data and cleanup procedure. Missing tools/packages require disclosure and explicit approval before installation; a check command does not authorize downloads.
- Report trial execution and cleanup separately; project-local state is not an OS sandbox.

## 10. Update affected records
- [README](../README.md): discovery links and layout when surfaces change; [ADOPTING](../ADOPTING.md): consumer routes, setup and evidence limits when adoption changes.
- [CHANGELOG](../CHANGELOG.md): notable changes under Unreleased; [CREDITS](../CREDITS.md): new/clarified sources and attribution, not invented historical certainty.
- Update the changed skill's `metadata.version` and compact Changelog consistently with its actual instruction change; do not version untouched skills.

## 11. Preserve provenance
- Record the source URL, inspected revision when known, declared license and copied versus conceptual use in CREDITS/evidence as appropriate.
- Use `not recorded` for missing revision/license facts; distinguish a metadata inspiration label from verified source inspection or permission to copy.
- Keep vendor exceptions and attribution intact. Frontmatter license declarations alone do not settle upstream licensing uncertainty.

## What lint checks / what it does not check
Source of truth: [tools/lint-skills.js](../tools/lint-skills.js); this section describes existing coverage, not additional lint rules.

**Discovery and blocking errors**
- Recursively discovers `SKILL.md` under five skill roots: backend, frontend, apps, shared and dotnet-10; excludes `.github/skills/` by default. `LINT_SKILLS_ROOTS` can override discovery.
- Requires opening/closing `---` delimiters and parseable YAML; missing/unclosed or invalid YAML stops checks for that file.
- Requires nonempty strings for `name`, `description`, `metadata.author`; `metadata.version` must be a present string (an empty string passes; semantic-version validity is not checked).
- Requires When to Use and Critical Patterns headings. Section checks accept levels 1–6 and case-insensitive prefix matches, not exact full heading text.
- Checks relative Markdown-link target existence only in Resources entered by `/^##\s+Resources?/i`, until the next `/^##\s/` heading. Skips HTTP(S)/pure anchors; strips fragments without validating them. Bare paths/URLs are not checked.
- No discovered skills is a blocking failure; read failures are reported as IO errors. Warnings alone exit 0; errors exit 1 (missing YAML dependency exits 2).

**Warnings, not blocking rules**
- Missing exact `Trigger:` in a valid description; missing Resources heading; exact Keywords/Tags equivalents (singular/plural, including Etiqueta(s) and Palabra(s) clave).
- Only for exact author `Zesh-One`: missing/non-string/blank license, missing/invalid/empty tools, or missing Changelog when the leading numeric major version is at least 2. Changelog is a warning despite the rule's "required" wording.
- Tools accept nonempty arrays or strings tokenized by whitespace/commas, including comma-separated legacy strings; lint does not enforce the rubric's string-only format or validate tool names/license identifiers.

**Not checked by lint**
- name syntax/length/directory alignment, semantic-version validity, file length, English-only prose or the full seven-section convention.
- Local links outside Resources, bare paths/URLs, fragments, attribution completeness or instruction correctness. Review these manually and with focused contracts; lint success is not spec/editorial approval.

## 12. Evidence and PR checklist
- [ ] Owner/gates retained; rubric reviewed; conditional routes and local links checked.
- [ ] Focused RED/GREEN and full verification recorded with command, exit code, counts and relevant environment; distinguish proposed checks from executed checks and list unperformed checks.
- [ ] Both pipelines registered; docs/version/provenance updated only where affected; diff and cleanup checked.
- [ ] Use [odd/evidence records](../odd/evidence/) as reporting examples: separate static instruction checks, model proposal observations and executed specimens.
- [ ] Do not equate text checks with model compliance, application correctness, accessibility conformance or a security guarantee; report browser/device/consumer evidence separately.
