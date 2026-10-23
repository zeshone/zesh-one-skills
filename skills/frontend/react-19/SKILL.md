---
name: react-19
description: >
  Operational React 19 rules for component boundaries, forms, and modern APIs.
  Trigger: When building or refactoring React 19 components, hooks, refs, or form actions.
license: Apache-2.0
allowed-tools: Read Edit Write Glob Grep
metadata:
  author: Zesh-One
  version: "1.5"
  inspired-by: gentleman-programming/react-19
---

## When to Use

Load this skill for React 19 components, hooks, forms and refs. First identify the installed version, bundler/framework, React Compiler status and existing data/mutation approach. Do not upgrade or introduce a framework merely to use an example.

**Runtime gate:** server/client directives and Server Components apply only to an RSC-capable runtime (for example the selected Next.js App Router profile). A React SPA such as Vite uses ordinary client components: no Server Components, Server Actions or `"use client"` scaffolding is required. Preserve TanStack Query and the existing client mutation path; React 19 alone does not supply server infrastructure.

## Critical Patterns

1. **Do** default to Server Components in an RSC-capable runtime; **don't** add `"use client"` by habit in that profile. **Why:** keeps bundles smaller.
2. **Do** add `"use client"` at an RSC client boundary only for state, events, refs, or browser APIs; **don't** use it for static render-only views. **Why:** client boundaries are expensive.
3. **Do** verify React Compiler activation before removing memoization; **don't** strip `useMemo`/`useCallback` blindly. **Why:** non-compiled apps can regress.
4. **Do** use named React imports; **don't** use default `React` import. **Why:** aligns with modern JSX runtime and linting.
5. **Do** preserve the existing form architecture; for React Actions, prefer `useActionState`. `ReactDOM.useFormState` was renamed to `React.useActionState` and is deprecated, not removed. **Don't** migrate working forms solely for this hook. **Why:** compatibility without changing form ownership.
6. **Do** treat `ref` as a prop in new components; **don't** introduce `forwardRef` for React 19-first code. **Why:** simpler component signatures.
7. **Do** use `use()` for supported promise reads where Suspense is already the chosen approach; **don't** replace a working Query/cache flow merely to adopt this API. **Why:** native async read semantics.
8. **Do** keep `use()` reads close to render usage; **don't** pass unresolved promises across deep trees without boundaries. **Why:** improves failure locality.
9. **Do** model writes as explicit mutations with server-side validation and suitable idempotency; use Server Actions only when supported by the chosen runtime, otherwise the existing client mutation path. **Don't** rely only on client-side optimistic assumptions. **Why:** prevents divergence.
10. **Do** follow this skill's current React 19 patterns; **don't** use legacy/deprecated React APIs from older majors; **do** verify exact-version syntax/examples/definitions in Context7 before version-sensitive work. **Why:** avoids stale guidance and runtime mismatches.

For an established React Actions form (not a required migration):

```typescript
import { useActionState } from "react";

export function ProfileForm() {
  const [state, action, isPending] = useActionState(saveProfileAction, null);
  return <form action={action}><button disabled={isPending}>{isPending ? "Saving..." : "Save"}</button></form>;
}
```

## Constraints & Tradeoffs

- Server-first composition improves performance but can constrain direct browser API access.
- React Compiler reduces manual memo work, but only when explicitly enabled in tooling.
- `use()` improves async ergonomics yet requires clear Suspense/error boundary strategy.
- `ref` as prop simplifies new code, but migration may coexist with legacy `forwardRef` wrappers.
- Form action flows centralize mutation handling, but can feel indirect versus inline handlers.

Boundary decision checklist:
- In an RSC-capable runtime, mark as server by default and justify each client boundary; this distinction is not required for a browser-only SPA.
- Confirm compiler status before deleting memoization primitives.
- Ensure every `use()` path has surrounding Suspense and error boundaries.
- Choose one owner for submission, pending state and errors in each form. React Actions, Query mutations and an existing form library need not replace one another; avoid duplicate submissions or competing pending/error state.
- Avoid mixed legacy/new ref patterns unless migration demands coexistence.

Compiler status checks:
- Only check whether React Compiler is already enabled in the selected framework/bundler configuration; verify against the installed version.
- Enabling it requires separate user authorization. If disabled, retain necessary memoization.
- Treat compiler status as environment-specific, not repository-wide by assumption.

## Anti-Patterns

1. Client-marking entire trees because one child uses events.
2. Adding competing pending/error state beside the established form submission owner.
3. Removing memoization without confirming compiler enablement.
4. Using `use()` outside controlled Suspense/error boundaries.
5. Mixing legacy and modern ref patterns in the same component without intent.
6. Importing broad client-only libraries into server component entry points.
7. Converting server components to client components only to run simple formatting logic.

## Progressive Disclosure

1. Confirm installed React version and check Context7 exact-version docs before API-sensitive refactors.
2. Apply the runtime gate; adopt server-vs-client boundaries only where supported.
3. Preserve the existing form architecture and explicit pending UX; use `useActionState` when React Actions already own submission.
4. Introduce `use()` for async reads where Suspense already exists.
5. Verify existing React Compiler activation before memo simplification; do not enable it without separate user authorization.
6. Only in the matching Next.js profile, cross-check [../nextjs-15/SKILL.md](../nextjs-15/SKILL.md); do not load it as generic React routing guidance.

Team rollout pattern:
- **Phase 1**: boundary rules + import hygiene.
- **Phase 2**: form ownership and pending/error UX verified within the existing flow.
- **Phase 3**: selective `use()` adoption in high-impact async reads.
- **Phase 4**: optimization cleanup only with verified compiler activation.

Execution guardrails:
- Keep one concern per component boundary refactor PR.
- Validate behavior snapshots before and after client/server boundary changes.
- Prefer incremental migration over whole-tree rewrites.

## Resources

- [React 19 Blog Announcement](https://react.dev/blog/2024/12/05/react-19)
- [React `useActionState`](https://react.dev/reference/react/useActionState)
- [React `use`](https://react.dev/reference/react/use)
- [React Server Components Reference](https://react.dev/reference/rsc/server-components)
- Only for a Next.js 15 runtime: [Next.js 15 conventions](../nextjs-15/SKILL.md)
- Every runtime: [shared security baseline](../../shared/security/baseline.md)

Related operational pairings:
- Use with [../typescript/SKILL.md](../typescript/SKILL.md) for typed component contracts.
- For sensitive mutations, load shared [application checks](../../shared/security/application-checks.md); add [../security/SKILL.md](../security/SKILL.md) only for its Next.js profile.
- Use with [../tanstack-query/SKILL.md](../tanstack-query/SKILL.md) for client cache orchestration.
- Keep boundary and data-fetching concerns separated during refactors.

## Changelog

### v1.5 — 2026-10-23
- Preserved form ownership, corrected useFormState deprecation, and required separate authorization for compiler activation.

### v1.4 — 2026-10-22
- Scoped RSC/directive/server-action guidance to supported runtimes; preserved SPA clients and existing query/mutation flows.

### v1.3 — 2026-04-22
- Aligned versioning metadata and changelog with today's frontend skill updates.

### v1.2 — 2026-04-21
- Standardized to operational format with exact required section set.
- Replaced long explanatory blocks with 10 atomic Do/Don't rules.
- Reduced snippets to one compact reference while preserving key React 19 guidance.
- Added boundary checklist and rollout notes to keep guidance practical.

### v1.1 — 2026-04-09
- Reworded memoization rule to make compiler prerequisite explicit.
