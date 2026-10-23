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
│       ├── cli-design/
│       └── github-pr/
├── docs/
│   ├── skills-review-rubric.md
│   └── skills-sanitization-plan.md
├── tools/
├── ADOPTING.md
├── CHANGELOG.md
└── CREDITS.md
```

## Adoption

Use the onboarding guide:

-> **[ADOPTING.md](./ADOPTING.md)**

## Validation

```bash
npm install
npm run verify
```

`verify` runs fixture tests, semantic checks, skill lint, and release-readiness checks. It remains Node-only and does not execute a consumer application.

For CLI interface work, load [cli-design](skills/shared/cli-design/SKILL.md) alongside the applicable domain. `npm run test:cli-instructions` runs its focused instruction/routing checks, also included in `test` and `verify`; passing them does not establish model compliance or executable CLI behavior.

With a .NET 10 SDK/reference packs already installed, `npm run test:cli` runs the separate [executable CLI acceptance specimen](skills/shared/cli-design/assets/line-tally/README.md). It uses project-local disposable state and no external package feeds; it is not part of default verification and does not install an SDK or packages.

For React SPA test selection and observable UI/form/API behavior, load [frontend testing](skills/frontend/testing/SKILL.md). `npm run test:frontend-instructions` checks its instruction contract and conditional route; these Node-only checks do not render React or prove browser behavior. Reuse consumer tooling; missing prerequisites require explicit approval, not automatic installation. The optional [React DOM trial](tools/fixtures/frontend-dom/README.md) runs with `npm run test:frontend-dom` after its separately authorized isolated prerequisite setup; it is not part of default verification or actual-browser proof.

The existing [.NET 8 unit-testing](skills/backend/testing-unit/SKILL.md) and [.NET 10 testing](skills/dotnet-10/testing/SKILL.md) profiles cover applicable rule boundaries, rejection effects and deterministic time/cancellation without unifying their libraries or adding infrastructure. `npm run test:dotnet-instructions` checks these instruction contracts; it does not execute xUnit.

## Changelog

-> **[CHANGELOG.md](./CHANGELOG.md)**

## Credits

-> **[CREDITS.md](./CREDITS.md)**
