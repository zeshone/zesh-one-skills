# Changelog

All notable changes to this repository are documented in this file.

## Unreleased

### Added
- Added `docs/MAINTAINING.md`, the repository maintainer workflow (owner selection, gates, first-party frontmatter versus the vendor template, conditional routes, instruction contracts, pipeline registration, records, provenance and exact lint coverage), linked from README; recorded app-skill inspirations in CREDITS without inferred revisions or licenses. `npm run test:maintainer-instructions` also checks that every instruction script runs exactly once in both default pipelines.
- Added shared web/hybrid accessibility guidance and frontend/apps routes with Node-only instruction checks; existing UI tooling stays optional and browser/assistive-technology verification remains explicit.
- Added shared evidence-led bug diagnosis, conditional routes for all four profiles and Node-only instruction checks; proposal observations remain distinct from executed application verification.
- Added an opt-in .NET10/xUnit capacity-rule trial with six public boundary/rejection cases, separately approved pinned local NuGet prerequisites and offline ordinary execution; default checks/CI remain Node-only.
- Added .NET testing instruction checks for applicable boundaries, rejected-operation effects and time/cancellation scenarios; preserved separate .NET 8/.NET 10 conventions and existing CLI/WAF/real-database obligations.
- Added an opt-in React19/jsdom rendered-component trial with six list/form/API-state cases, isolated fixed prerequisites and explicit installation approval; root dependencies and Node-only CI remain unchanged.
- Added React-SPA-scoped frontend testing guidance, a conditional route and Node-only instruction checks; test layers and UI/form/API evidence remain distinct from browser, model-compliance and security guarantees.
- Added an opt-in .NET 10 `line-tally` acceptance specimen and `npm run test:cli` public-process checks, with project-local disposable state and explicit execution limits; default verification still requires only Node.
- Added the shared `cli-design` interaction contract, conditional .NET routes and focused instruction checks (`npm run test:cli-instructions`). Default `test`/`verify` pipelines remain Node-only; architecture selection and security stay with their existing owners.

### Changed
- Corrected and scoped React-side skills: form ownership instead of forced `useActionState`/React Hook Form migrations, `useFormState` deprecation, Zod 4 deprecated-versus-removed APIs and input/output types, both valid Zustand 5 `useShallow` paths, deterministic Query key hashing, Query/Zustand adoption only when installed or selected, SSR client/store lifetimes, separate Next.js and Query cache ownership, concurrency-safe optimistic updates, Tailwind token/dynamic-class rules, installed shadcn primitive variants and no unauthorized React Compiler activation. Checked by `npm run test:react-stack-instructions` (Node-only).
- Scoped Ionic/Angular guidance to the existing standalone or NgModule/RxJS architecture, freshness-aware page entry with leave cleanup, Router/NavController/tab escape hatches and version-gated template/DI idioms; Capacitor now separates platform, plugin availability and plugin-specific permissions. Checked by `npm run test:apps-instructions` (Node-only).
- Integrated the shared everyday security guardrails while retaining distinct .NET 8/.NET 10 security rules and exceptions.
- Added contextual, user-approved architecture selection to existing .NET guidance; removed automatic second-context assembly promotion and scoped HTTP/database recipes away from local CLIs.
- Scoped React/Query/Next guidance by runtime and Ionic/Capacitor guidance by installed capabilities; extended shared security to SPA and mobile without importing server-only infrastructure.
- Added pinned source attribution, six BEFORE scenarios and instruction regression checks (`npm run test:pilot`, included in `test` and `verify`). Model proposals are not execution or security guarantees.

## v2.9.0 — 2026-10-27

### Added
- Added `zesh-skills --install-cli` for the current user: it copies the running executable to a user-owned location (Windows `%LOCALAPPDATA%\Programs\zesh-skills`, macOS/Linux `~/.local/bin`) and persists only a user-level PATH entry or one marked shell startup block, with preview, confirmation, idempotence, post-preview drift detection and rollback. It needs no administrator rights, keeps the original file, installs no skills and never touches the machine PATH. Checked by Go service tests, public CLI flow tests with isolated homes and a fake registry-backed user PATH, `go vet` and the six CGO-disabled targets.

### Changed
- Rewrote the public user documentation (`README.md`, `ADOPTING.md`, `cli/README.md`) in neutral Spanish and removed the Rioplatense voseo it carried; commands, flags, paths, links and proof limits are unchanged. The download section no longer claims an unpublished release: it links the release page and states the source-build fallback instead.

The existing [v2.8.0 release](https://github.com/zeshone/zesh-one-skills/releases/tag/v2.8.0) remains the source for that release's history.

## v2.8.0 — 2026-10-26

### Added
- Standalone project-local `zesh-skills` CLI for Claude Code, Codex, OpenCode and Pi, with a checkbox wizard, explained architecture choices and explicit preview/confirmation.
- English-default and opt-in Spanish interface (`--es`), including 35 bundled Spanish descriptions; canonical English IDs, metadata and installed formats stay unchanged.
- Local safety boundaries: commit-pinned catalog/content, validated paths and hashes, conflict confirmation, project lock and ordinary-error rollback; no global skill installation or downloaded script execution.
- Explicit `--release 2.8.0` preparation for six raw Windows/Linux/macOS amd64/arm64 assets with injected version, collision-free filenames and `SHA256SUMS`; download/checksum guidance documents publication pending, not an already-public CLI release.
- Child-PR CI coverage and module-guarded Go test/vet jobs for Ubuntu, Windows and macOS; default npm checks remain Node-only. Local fixtures do not establish online CI or all-platform runtime success.

The existing [v2.7.0 release](https://github.com/zeshone/zesh-one-skills/releases/tag/v2.7.0) remains the source for that release's history.

## v2.6.0 — 2026-04-22

### Changed
- Reinforced frontend skill governance with version-aware Context7 checkpoints across `nextjs-15`, `react-19`, `security`, `tailwind-4`, and `zod-4`.
- Upgraded `tanstack-query` to v2.1 with explicit error semantics (`queryFn` throw), key-stability guidance, SSR prefetch decision matrix, mutation template, and incremental migration playbook.
- Upgraded `typescript` to v2.1 with strict defaults: no new `any` without expiring waiver, schema-first `z.infer` contracts, const-type finite domains, exception policy, and pre-merge checklist.
- Upgraded `zustand-5` to v2.1 with mandatory applicability gate (`not applicable` when Redux legacy exists), no-default-coexistence policy, and quick verification checklist.
- Upgraded `shadcn-ui` to v2.1 as explicit opt-in guidance (non-mandatory), focused on shadcn best practices and latest-version checks through Context7.
- Normalized frontend skill metadata/changelogs for same-day edits (`nextjs-15` v2.2, `react-19` v1.3, `security` v2.1, `tailwind-4` v2.2, `zod-4` v1.2).

### Quality Gate
- `npm test` passing with fixtures, semantic checks, lint (`0 errors / 0 warnings`), release-readiness checks, and zero-skills guard checks.

## v2.5.0 — 2026-04-21

### Added
- Added `docs/skills-review-rubric.md` with a 3-layer review model:
  - spec compliance,
  - Agent Skills best-practice quality,
  - repo philosophy fit.
- Added `docs/skills-sanitization-plan.md` with P1/P2/P3 execution phases and completion criteria.
- Added root-level `CHANGELOG.md` as the canonical release history entry point.

### Changed
- Reworked all `skills/**/SKILL.md` into concise operational guidance with consistent structure:
  - `When to Use`
  - `Critical Patterns`
  - `Constraints & Tradeoffs`
  - `Anti-Patterns`
  - `Progressive Disclosure`
  - `Resources`
  - `Changelog`
- Applied English-only normalization across skills and domain AGENTS files.
- Updated `AGENTS.backend.md`, `AGENTS.frontend.md`, and `AGENTS.apps.md` to concise domain-operational routing guides.
- Updated all skill frontmatter to align with current Agent Skills expectations (`name` alignment and `allowed-tools` string format).
- Standardized/compacted changelog sections in SKILL files to reduce context noise.
- Updated linter behavior for FM-007 in `tools/lint-skills.js` to accept spec-format `allowed-tools` strings while preserving legacy compatibility.
- Expanded fixture coverage in `tools/test-fixtures.js` for FM-007 valid-string and empty-string edge cases.
- Updated `README.md` to reflect the new operational-sanitized state and current repo structure.

### Quality Gate
- `npm test` passing with:
  - fixtures passing,
  - semantic checks passing,
  - `lint-skills` reporting `0 errors / 0 warnings`,
  - release-readiness checks passing.

## v2.4.0 — 2026-04-17
- Judgment Day release: apps domain additions, linter hardening, and full-suite audit alignment.
