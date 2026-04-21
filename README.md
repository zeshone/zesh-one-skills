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
│   │   ├── typescript/
│   │   ├── zod-4/
│   │   └── zustand-5/
│   ├── apps/
│   │   ├── capacitor/
│   │   └── ionic-angular/
│   └── shared/
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

`verify` runs fixture tests, semantic checks, skill lint, and release-readiness checks.

## Changelog

-> **[CHANGELOG.md](./CHANGELOG.md)**

## Credits

-> **[CREDITS.md](./CREDITS.md)**
