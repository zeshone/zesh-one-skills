# Frontend Domain — ZeshOne Development Guidelines

## How to Use This Guide

- Start here for cross-project norms on frontend development.
- Each skill has detailed patterns available on-demand.
- Domain skills override this file when guidance conflicts.

---

## Stack Base

| Layer | Technology |
|-------|-----------|
| Framework | Next.js 15 (App Router) |
| UI Library | React 19 |
| Language | TypeScript 5 (strict mode) |
| Styling | Tailwind CSS v4 (CSS-first, no config file) |
| Server state | TanStack Query v5 |
| Client state | Zustand 5 |
| Forms | React Hook Form + Zod 4 |
| UI Components | shadcn/ui (radix-ui base, oklch colors) |

**API pattern:** All external API calls go through Server Actions — Client → Server Action → .NET API. Endpoints are never exposed to the browser.

---

## Available Skills

### Frontend Skills (ZeshOne)

| Skill | Description | Path |
|-------|-------------|------|
| `typescript` | Const types, flat interfaces, no any | [SKILL.md](skills/frontend/typescript/SKILL.md) |
| `react-19` | React Compiler, use(), useActionState, no forwardRef | [SKILL.md](skills/frontend/react-19/SKILL.md) |
| `nextjs-15` | App Router, Vertical Slices, Server Actions as API proxy, project structure, performance | [SKILL.md](skills/frontend/nextjs-15/SKILL.md) |
| `tailwind-4` | cn(), CSS-first config, fluid typography, container queries, CSS modern features | [SKILL.md](skills/frontend/tailwind-4/SKILL.md) |
| `zod-4` | Breaking changes from v3, safeParse, error param | [SKILL.md](skills/frontend/zod-4/SKILL.md) |
| `zustand-5` | Per-feature stores, selectors, useShallow v5 path | [SKILL.md](skills/frontend/zustand-5/SKILL.md) |
| `tanstack-query` | Query key factories, prefetch + HydrationBoundary, mutations, optimistic updates, infinite queries | [SKILL.md](skills/frontend/tanstack-query/SKILL.md) |
| `shadcn-ui` | shadcn/ui components, React Hook Form + Zod pattern, dark mode via next-themes | [SKILL.md](skills/frontend/shadcn-ui/SKILL.md) |
| `security` | CSP nonce, security headers, auth guards, XSS prevention, cookies, rate limiting | [SKILL.md](skills/frontend/security/SKILL.md) |

### Shared Skills (ZeshOne)

| Skill | Description | Path |
|-------|-------------|------|
| `github-pr` | Branch naming, PR conventions, base branch alfa | [SKILL.md](skills/shared/github-pr/SKILL.md) |

---

## Auto-invoke Skills

When performing these actions, **ALWAYS** invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| Writing TypeScript types or interfaces | `typescript` |
| Writing React components | `react-19` + `typescript` |
| App Router / Server Actions / Next.js routing | `nextjs-15` |
| Defining project structure or feature organization | `nextjs-15` |
| Working with Tailwind CSS classes | `tailwind-4` |
| Creating or updating Zod schemas | `zod-4` |
| Building forms with validation | `shadcn-ui` + `zod-4` |
| Using or configuring Zustand stores | `zustand-5` |
| Fetching data, caching, mutations | `tanstack-query` |
| Building UI components (buttons, dialogs, tables) | `shadcn-ui` |
| Security headers, auth guards, CSP, input validation | `security` |
| Creating or reviewing a Pull Request | `github-pr` |

---

## Critical Cross-Skill Rules

- Server Actions are the ONLY way to call .NET APIs — never from the client.
- `params` and `searchParams` in Next.js 15 are Promises — always `await`.
- React Compiler must be explicitly enabled — verify before removing `useMemo`/`useCallback`.
- NEVER use `px` for font sizes or spacing — use `rem` or `clamp()`.
- NEVER use `var()` inside `className` — use Tailwind semantic classes.
- NEVER use React Context for server/API data — use TanStack Query.
- All user input validated with Zod on the server (Server Action or API route).
- Sensitive cookies: always `httpOnly + secure + sameSite`.
- No secrets in `NEXT_PUBLIC_` env vars.
