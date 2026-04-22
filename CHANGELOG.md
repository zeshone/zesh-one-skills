# Changelog

All notable changes to this repository are documented in this file.

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
