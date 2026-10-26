# Zesh-One-Skills

> AI agent skill suite for backend (.NET 8 REST), frontend (Next.js 15 + React 19), and mobile apps (Ionic + Capacitor), built for OpenCode + Engram + Agent Teams Lite workflows.

## What this repository is

`zesh-one-skills` centralizes domain-specific execution guidance for AI agents.

The suite is intentionally **operational**:
- concise defaults,
- explicit constraints,
- anti-patterns and progressive disclosure,
- minimal context noise.

All skill and AGENTS domain files follow an **English-only** policy.

## Tooling stack

| Tool | Role |
|------|------|
| [OpenCode](https://opencode.ai) | Agent runtime |
| [Engram](https://github.com/gentleman-programming/engram) | Persistent memory across sessions |
| Agent Teams Lite | Multi-agent orchestration/delegation |
| Node.js | Local validation pipeline (`npm run verify`) |

## Repository layout

```text
zesh-one-skills/
├── AGENTS.backend.md
├── AGENTS.frontend.md
├── AGENTS.apps.md
├── skills/
│   ├── backend/
│   │   ├── dataaccess/
│   │   ├── general/
│   │   ├── logging/
│   │   ├── mapping/
│   │   ├── performance/
│   │   ├── requests/
│   │   ├── responses/
│   │   ├── security/
│   │   ├── testing-unit/
│   │   └── validations/
│   ├── frontend/
│   │   ├── nextjs-15/
│   │   ├── react-19/
│   │   ├── security/
│   │   ├── shadcn-ui/
│   │   ├── tanstack-query/
│   │   ├── tailwind-4/
│   │   ├── testing/
│   │   ├── typescript/
│   │   ├── zod-4/
│   │   └── zustand-5/
│   ├── apps/
│   │   ├── capacitor/
│   │   └── ionic-angular/
│   └── shared/
│       ├── accessibility/
│       ├── cli-design/
│       ├── diagnosing-bugs/
│       └── github-pr/
├── docs/
│   ├── MAINTAINING.md
│   ├── skills-review-rubric.md
│   └── skills-sanitization-plan.md
├── tools/
├── ADOPTING.md
├── CHANGELOG.md
└── CREDITS.md
```

## Project-local CLI installer

Run `zesh-skills` from a consumer project's root to choose skills by scope and topic, with descriptions and no mandatory bundles. It supports Claude Code, Codex, OpenCode and Pi, preserves existing instructions, and reads the latest catalog/content from `main` at one pinned commit per run.

See [CLI usage, standalone builds and installation](cli/README.md). Build locally with `npm run build:zesh-cli`; no publication or global installation is performed. The production download flow becomes available after the new `skills-index.json` is merged into `main`.

Maintainers: run `npm run catalog:generate` after changing skills. The existing Node-only checks validate catalog freshness; executable Go checks run separately with `npm run test:zesh-cli`.

## Adoption

Use the onboarding guide:

-> **[ADOPTING.md](./ADOPTING.md)**

For repository changes, use the [maintainer workflow](docs/MAINTAINING.md).

## Validation

```bash
npm install
npm run verify
```

`verify` runs fixture tests, semantic checks, skill lint, and release-readiness checks. It remains Node-only and does not execute a consumer application.

For CLI interface work, load [cli-design](skills/shared/cli-design/SKILL.md) alongside the applicable domain. `npm run test:cli-instructions` runs its focused instruction/routing checks, also included in `test` and `verify`; passing them does not establish model compliance or executable CLI behavior.

With a .NET 10 SDK/reference packs already installed, `npm run test:cli` runs the separate [executable CLI acceptance specimen](skills/shared/cli-design/assets/line-tally/README.md). It uses project-local disposable state and no external package feeds; it is not part of default verification and does not install an SDK or packages.

For React SPA test selection and observable UI/form/API behavior, load [frontend testing](skills/frontend/testing/SKILL.md). `npm run test:frontend-instructions` checks its instruction contract and conditional route; these Node-only checks do not render React or prove browser behavior. Reuse consumer tooling; missing prerequisites require explicit approval, not automatic installation. The optional [React DOM trial](tools/fixtures/frontend-dom/README.md) runs with `npm run test:frontend-dom` after its separately authorized isolated prerequisite setup; it is not part of default verification or actual-browser proof.

The existing [.NET 8 unit-testing](skills/backend/testing-unit/SKILL.md) and [.NET 10 testing](skills/dotnet-10/testing/SKILL.md) profiles cover applicable rule boundaries, rejection effects and deterministic time/cancellation without unifying their libraries or adding infrastructure. `npm run test:dotnet-instructions` checks these instruction contracts; it does not execute xUnit. The optional [isolated xUnit trial](tools/fixtures/dotnet-testing/README.md) runs with `npm run test:dotnet-unit` after explicitly authorized local NuGet setup; ordinary execution never downloads packages.

For defect investigation across profiles, load [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md): reproduce, separate hypotheses from observations, and verify the smallest correction. `npm run test:debugging-instructions` checks its instructions/routes, not a consumer bug or agent compliance; see the [bounded proposal evidence](odd/evidence/debugging-accessibility-adoption.md).

For changed web/hybrid interactions, load [accessibility](skills/shared/accessibility/SKILL.md) through the frontend/apps guide. `npm run test:accessibility-instructions` checks semantic, keyboard/focus and feedback instructions and routing; it does not run a browser, assistive technology or compliance audit. Preserve existing UI libraries and use real consumer journeys for those checks; no tool installation is automatic.

## Changelog

-> **[CHANGELOG.md](./CHANGELOG.md)**

## Credits

-> **[CREDITS.md](./CREDITS.md)**
