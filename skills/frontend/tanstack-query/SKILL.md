---
name: tanstack-query
description: >
  ZeshOne TanStack Query defaults for server/API state management.
  Trigger: When implementing async data fetching, caching, mutations, invalidation, SSR prefetch, or hydration.
license: Apache-2.0
allowed-tools: Read Write Edit Bash
metadata:
  author: Zesh-One
  version: "2.0"
  inspired-by: frontend-nextjs/tanstack-query.md
---

## When to Use

Load this skill whenever state comes from an API, database proxy, or remote service.

Use TanStack Query as the single source of truth for server state in React/Next.js.
Combine with `nextjs-15` for Server Component prefetch flows and `zustand-5` for non-server global UI state only.

## Critical Patterns

Atomic rules (Do/Don't + Why):

1. **Do** use TanStack Query for ALL server/API state. **Don't** use React Context for fetch/cache state. **Why:** Context does not provide cache lifecycle, dedupe, or invalidation semantics.
2. **Do** centralize `QueryClient` defaults in one module. **Don't** instantiate ad-hoc clients in features. **Why:** inconsistent defaults produce unpredictable cache behavior.
3. **Do** define query key factories per feature (`all/lists/list/details/detail`). **Don't** inline array keys repeatedly. **Why:** factories prevent key drift and invalidation bugs.
4. **Do** wrap fetch definitions in `queryOptions`/`infiniteQueryOptions`. **Don't** duplicate `queryFn` + keys across hooks/prefetch calls. **Why:** one contract powers hooks, prefetch, and tests.
5. **Do** prefetch in Server Components when first paint needs data. **Don't** force client-only cold starts for critical lists/details. **Why:** hydration with warm cache removes initial loading gaps.
6. **Do** use `HydrationBoundary` with dehydrated state from the same query options. **Don't** hydrate a different key/options shape. **Why:** mismatch causes duplicate requests and stale UI.
7. **Do** invalidate narrowly after mutations (detail + affected lists). **Don't** invalidate everything by default. **Why:** broad invalidation increases network churn and UI flicker.
8. **Do** use optimistic updates only with rollback context. **Don't** mutate cache without recovery strategy. **Why:** failed mutations otherwise leave corrupted UI state.
9. **Do** keep query functions pure and deterministic. **Don't** hide navigation/toast side effects inside `queryFn`. **Why:** side effects belong in mutation callbacks/effects.
10. **Do** keep ownership boundary clear: server state in TanStack Query, app UI state in Zustand/local state. **Don't** mix responsibilities inside one store. **Why:** mixed ownership destroys maintainability.

Single reference snippet (factory + options + server prefetch):

```ts
// feature-level contract reused by hooks + prefetch
export const postKeys = { all: ["posts"] as const, detail: (id: string) => [...postKeys.all, "detail", id] as const };
export const postDetailOptions = (id: string) => queryOptions({ queryKey: postKeys.detail(id), queryFn: () => api.posts.get(id) });
// Server Component: await queryClient.prefetchQuery(postDetailOptions(id)); then <HydrationBoundary state={dehydrate(queryClient)}>
```

## Constraints & Tradeoffs

- TanStack Query adds conceptual overhead; pay this to gain consistent async-state architecture.
- Aggressive `staleTime` lowers refetches but can show stale data longer.
- Frequent invalidation improves freshness but increases API cost.
- Server prefetch improves first render but can increase server latency if overused.
- Optimistic updates improve perceived speed but require careful rollback paths.
- Hydration payload size grows with over-prefetching; prefetch only what first paint needs.
- Query key object parameters must remain stable/serializable to avoid accidental cache fragmentation.

Operational verification before merge:

- Confirm no API state lives in React Context providers.
- Confirm each feature exposes key factory + options helper pair.
- Confirm hooks and server prefetch both consume the same options helpers.
- Confirm mutation flows invalidate exactly the impacted detail/list keys.
- Confirm SSR pages use `dehydrate` + `HydrationBoundary` when warm-first-render is required.

## Anti-Patterns

- React Context provider storing API payloads and loading flags.
- `useEffect + fetch + useState` for endpoints already managed by Query.
- Query keys constructed inline with changing object literals in components.
- Different key shapes for read and invalidation paths.
- Mutation success handlers that never invalidate dependent caches.
- One giant "global query file" instead of feature-level contracts.
- Prefetching everything "just in case" without first-render need.
- Calling `new QueryClient()` inside render paths.
- Mixing UI-only state (modals, theme toggles) into query caches.
- Manually duplicating optimistic cache writes in multiple components.

## Progressive Disclosure

Implement only the next needed layer:

1. **Foundational contract**
   - feature query key factory
   - feature `queryOptions` helpers
2. **Basic consumption**
   - `useQuery`/`useInfiniteQuery` using helpers only
3. **Mutation coherence**
   - scoped invalidation
   - optimistic update + rollback when UX requires it
4. **SSR performance**
   - prefetch in Server Component
   - `dehydrate` + `HydrationBoundary` reuse same options contract

Do not jump to optimistic patterns before baseline key/options discipline is solid.

Task sizing guidance:

- Small task: add key factory/options for one endpoint.
- Medium task: migrate one feature from `useEffect` fetching to Query hooks + invalidation.
- Large task: cross-feature SSR prefetch and hydration standardization.

## Resources

- TanStack Query overview: https://tanstack.com/query/latest/docs/framework/react/overview
- TanStack Query keys: https://tanstack.com/query/latest/docs/framework/react/guides/query-keys
- TanStack Query advanced SSR: https://tanstack.com/query/latest/docs/framework/react/guides/advanced-ssr
- TanStack Query mutations: https://tanstack.com/query/latest/docs/framework/react/guides/mutations
- Cross-reference: `skills/frontend/nextjs-15/SKILL.md`
- Cross-reference: `skills/frontend/zustand-5/SKILL.md`
- Cross-reference: `skills/frontend/security/SKILL.md` (protect private pages before query prefetch)
- Cross-reference: `skills/frontend/typescript/SKILL.md`

Quick review checklist:
- Every query uses a factory-generated key.
- Every query/mutation is declared via options helper.
- Every mutation has a defined invalidation plan.
- No server state persisted in Context or UI stores.

## Changelog

### v2.0 — 2026-04-21
- Rewritten as concise operational guidance with repo-specific defaults and strict ownership of server/API state.
- Standardized key factories, `queryOptions` helpers, and Server prefetch/HydrationBoundary flow.

### v1.9 — 2026-04-16
- Consolidated prior baseline guidance that preceded the v2.0 operational rewrite.
