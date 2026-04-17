# Frontend Domain — ZeshOne Development Guidelines

## How to Use This Guide

- Start here for cross-project norms on frontend development.
- Each skill has detailed patterns available on-demand.
- Domain skills override this file when guidance conflicts.

---

## Available Skills

### Frontend Skills (ZeshOne)

| Skill | Description | Path |
|-------|-------------|------|
| `typescript` | Const types, flat interfaces, no any | [SKILL.md](skills/frontend/typescript/SKILL.md) |
| `react-19` | React Compiler, use(), useActionState, no forwardRef | [SKILL.md](skills/frontend/react-19/SKILL.md) |
| `nextjs-15` | App Router, Vertical Slices, Server Actions as API proxy | [SKILL.md](skills/frontend/nextjs-15/SKILL.md) |
| `tailwind-4` | cn(), no var() in className, @theme config | [SKILL.md](skills/frontend/tailwind-4/SKILL.md) |
| `zod-4` | Breaking changes from v3, safeParse, error param | [SKILL.md](skills/frontend/zod-4/SKILL.md) |
| `zustand-5` | Per-feature stores, selectors, useShallow v5 path | [SKILL.md](skills/frontend/zustand-5/SKILL.md) |

### Shared Skills (ZeshOne)

| Skill | Description | Path |
|-------|-------------|------|
| `github-pr` | Branch naming, PR conventions, base branch alfa | [SKILL.md](skills/shared/github-pr/SKILL.md) |

### Vendor Skills (Reference Only — do not modify)

| Vendor | Skills | Path |
|--------|--------|------|
| Gentleman | typescript, react-19, nextjs-15, tailwind-4, zod-4, zustand-5, github-pr | [skills/vendor/gentleman/](skills/vendor/gentleman/) |

---

## Auto-invoke Skills

When performing these actions, **ALWAYS** invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| Writing TypeScript types or interfaces | `typescript` |
| Writing React components | `react-19` + `typescript` |
| App Router / Server Actions / Next.js routing | `nextjs-15` |
| Working with Tailwind CSS classes | `tailwind-4` |
| Creating or updating Zod schemas | `zod-4` |
| Building forms with validation | `zod-4` + `react-19` |
| Using or configuring Zustand stores | `zustand-5` |
| Creating or reviewing a Pull Request | `github-pr` |
