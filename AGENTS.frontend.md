# Frontend Domain Agent Guide (Next.js 15 + React 19)

## Scope
Use this file as the execution default for frontend work in this repo. Keep outputs focused on shipping features safely and fast. If a domain skill conflicts with this file, the skill wins.

## Mandatory Skill Routing (load first)

| Work Type | Skill |
|---|---|
| Type modeling, strict typing, interface strategy | `typescript` |
| React component behavior and modern APIs | `react-19` |
| Routing, App Router architecture, Server Actions | `nextjs-15` |
| Styling strategy and Tailwind patterns | `tailwind-4` |
| Schema validation | `zod-4` |
| Local client state | `zustand-5` |
| Server state, caching, mutations, hydration | `tanstack-query` |
| UI primitives/forms with shadcn patterns | `shadcn-ui` |
| Security headers, CSP, guards, cookie safety | `security` |
| PR preparation/review | `github-pr` |

## Operational Defaults
- Stack baseline: Next.js 15 App Router, React 19, TypeScript strict, Tailwind v4.
- API boundary is fixed: Client -> Server Action -> .NET API. Do not expose backend endpoints directly to browser clients.
- Validate user input on the server boundary with Zod (Server Action/API route).
- Use TanStack Query for server/API data and caching; use Zustand for local UI state.
- Keep component APIs narrow and type-first.

## Non-Negotiables
- Do not call .NET APIs directly from client components.
- Do not treat `params`/`searchParams` as synchronous in Next.js 15.
- Do not remove memoization patterns unless React Compiler enablement is verified.
- Do not use `px` for typography/spacing tokens in new code.
- Do not place secrets in `NEXT_PUBLIC_` variables.
- Do not use React Context as a replacement for server-state tools.
