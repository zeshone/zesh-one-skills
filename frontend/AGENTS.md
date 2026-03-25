# Frontend Domain — Development Guidelines

## How to Use This Guide

- Start here for cross-project norms on frontend development with the Gentleman stack.
- Each skill has detailed patterns available on-demand.
- Domain skills override this file when guidance conflicts.

---

## Available Skills

Use these skills for detailed patterns on-demand:

### Frontend Skills (Gentleman Stack)
| Skill | Description | URL |
|-------|-------------|-----|
| `typescript` | Const types, flat interfaces, utility types | [SKILL.md](../vendor/gentleman/typescript/SKILL.md) |
| `react-19` | No useMemo/useCallback, React Compiler patterns | [SKILL.md](../vendor/gentleman/react-19/SKILL.md) |
| `nextjs-15` | App Router, Server Actions, streaming | [SKILL.md](../vendor/gentleman/nextjs-15/SKILL.md) |
| `tailwind-4` | cn() utility, no var() in className | [SKILL.md](../vendor/gentleman/tailwind-4/SKILL.md) |
| `zod-4` | New API (z.email(), z.uuid()) | [SKILL.md](../vendor/gentleman/zod-4/SKILL.md) |
| `zustand-5` | Persist, selectors, slices | [SKILL.md](../vendor/gentleman/zustand-5/SKILL.md) |

### Shared Skills
| Skill | Description | URL |
|-------|-------------|-----|
| `github-pr` | Pull Request conventions, branch workflow, PR checklist | [SKILL.md](../shared/skills/github-pr/SKILL.md) |

---

## Auto-invoke Skills

When performing these actions, **ALWAYS** invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| Writing TypeScript types or interfaces | `typescript` |
| Writing React components | `react-19` |
| App Router / Server Actions / Next.js routing | `nextjs-15` |
| Working with Tailwind CSS classes | `tailwind-4` |
| Creating or updating Zod schemas | `zod-4` |
| Using or configuring Zustand stores | `zustand-5` |
| Creating or reviewing a Pull Request | `github-pr` |
