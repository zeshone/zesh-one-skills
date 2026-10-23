# Frontend Domain Agent Guide (React 19 SPA / Next.js 15)

## Scope
Inspect the installed framework, bundler, versions and existing data-access decisions first. Preserve that stack; this guide does not authorize a migration or new dependencies. If a domain skill conflicts with this file, the skill wins within its applicable runtime profile. For new structural choices, use the [shared selection procedure](skills/shared/architecture-selection.md).

## Always-on security baseline
Load [the shared security baseline](skills/shared/security/baseline.md) at the start of every development task. Load matching [application checks](skills/shared/security/application-checks.md) or [delivery checks](skills/shared/security/delivery-checks.md) only for affected boundaries. Preserve existing controls in every profile. Client route guards are UX, not backend authorization; the API must enforce authentication, resource authorization and input validation.

Copy `skills/shared/security/` with the domain and adjust paths when installed elsewhere. The host must actually load these references; links are not automatic activation or isolation.

## Mandatory Skill Routing (load first, when applicable)

| Work Type | Skill |
|---|---|
| Type modeling, strict typing, interface strategy | `typescript` |
| Bug reproduction, diagnosis, regression investigation | [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md) |
| React component behavior and modern APIs | `react-19` for the installed React 19 runtime |
| Test selection, UI states, forms and API interactions | [testing](skills/frontend/testing/SKILL.md), currently React SPA only; preserve installed tooling and form/client flow |
| Next.js routing, App Router, Server Actions | `nextjs-15` only for the matching Next.js profile |
| Styling with installed Tailwind v4 | `tailwind-4` |
| Schema validation with Zod | `zod-4` |
| Local client state with an applicable Zustand setup | `zustand-5`; preserve its applicability gate |
| Server state, caching, mutations, optional SSR hydration | `tanstack-query` |
| Opted-in shadcn UI primitives/forms | `shadcn-ui` |
| Next.js headers, middleware guards, cookies | `security` for Next.js; shared checks for other runtimes |
| PR preparation/review | `github-pr` |

## Runtime profiles

### React SPA (for example Vite)
- Keep the existing authenticated client: component/query hook -> client/service -> .NET API. Browser access to that API is valid in this profile; preserve server-side controls and the application's established authentication design.
- Do not introduce Next.js, Server Components, Server Actions, middleware or an SSR server merely to follow this guide. Use the existing router.
- Use TanStack Query for API state/cache; keep loading/error/empty states and the existing client mutation path. Client validation helps UX but cannot replace server validation.

### Next.js profile
- Keep the selected Next.js 15 App Router/React 19/TypeScript strict/Tailwind v4 conventions when those are the project's versions.
- API boundary: Client -> Server Action -> .NET API. Do not call .NET APIs directly from client components in this profile. This is encapsulation, not an authorization control.
- Validate user input with Zod at the Server Action/API route boundary. Keep the Next.js security skill's middleware, CSP and cookie protections.
- Do not treat `params`/`searchParams` as synchronous in Next.js 15. Use Server Components/prefetch only where the runtime supports them.

## Shared defaults
- Keep component APIs narrow and type-first; preserve the project's chosen local UI-state approach.
- Do not remove memoization patterns unless React Compiler enablement is verified.
- Do not use `px` for typography/spacing tokens in new code under the applicable styling profile.
- Do not place secrets in client bundles, including `NEXT_PUBLIC_` or `VITE_` variables.
- Do not use React Context as a replacement for server-state tools.
