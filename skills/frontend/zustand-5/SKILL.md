---
name: zustand-5
description: >
  Operational Zustand 5 rules for store boundaries, selectors, and async state flow.
  Trigger: When creating, refactoring, or reviewing Zustand 5 stores and selectors.
license: Apache-2.0
allowed-tools: Read Write Edit Bash
metadata:
  author: Zesh-One
  version: "2.2"
  inspired-by: gentleman-programming/zustand-5
---

## Applicability Gate (Mandatory)

Use this gate BEFORE applying any rule in this skill.

### Applies when

- Zustand 5 is installed or explicitly selected for the affected client global/shared UI state.
- The project has no conflicting established state standard, or an approved architectural migration plan covers this scope.

### Does NOT apply when

- Redux or other established state solutions own this scope without an approved architectural migration plan.
- Zustand has not been selected; a need for shared UI state alone does not authorize adoption.

### Required behavior when it does NOT apply

- Return: **"not applicable"**.
- Respect existing Redux legacy architecture and other established state solutions.
- Do **not** force migration, parallel adoption, or coexistence by default.

### Governance Decision Enforced

- No coexistence of state-management packages by default.
- Apply Zustand rules only to an installed or explicitly selected Zustand setup; do not replace other established state solutions.
- If Redux legacy exists, this `zustand-5` skill is out of scope unless an approved architectural migration plan covers the change.

## When to Use

Load this skill when designing global client state, selector usage, async actions, or persisted feature state with Zustand 5.

## Critical Patterns

1. **Do** create one store per feature/domain; **don't** centralize all app state in a mega-store. **Why:** limits coupling and rerender radius.
2. **Do** use selectors for every component read; **don't** subscribe to the whole store object. **Why:** prevents broad invalidations.
3. **Do** use `useShallow` for object/tuple multi-selectors; **don't** return new object literals unguarded. **Why:** avoids avoidable rerenders.
4. **Do** prefer `zustand/react/shallow` for `useShallow`; `zustand/shallow` and `zustand/react/shallow` are both valid in v5. **Don't** rewrite a valid installed import solely for compatibility. **Why:** preference is not an API restriction.
5. **Do** keep async actions in feature action modules; **don't** embed transport/service details directly in UI components. **Why:** separation of concerns.
6. **Do** expose explicit loading/error flags in async store actions; **don't** hide request lifecycle in implicit side effects. **Why:** deterministic UI states.
7. **Do** prefer local `useState` for component-local state; **don't** globalize ephemeral UI toggles by default. **Why:** reduces store complexity.
8. **Do** keep server/API state outside Zustand; **don't** use Zustand as remote cache. **Why:** remote data belongs to the established server-state solution (TanStack Query when selected) and has different invalidation/lifecycle semantics.
9. **Do** use `persist` only for state with clear rehydration value; **don't** persist volatile or sensitive transient data. **Why:** storage correctness and security.
10. **Do** namespace persist keys with app/feature prefix; **don't** use generic keys (`settings`, `store`). **Why:** collision prevention.
11. **Do** follow this skill's current Zustand 5 patterns; **don't** use legacy/deprecated APIs from older Zustand versions; **do** verify exact-version syntax/examples/definitions in Context7 before version-sensitive store changes. **Why:** avoids stale imports and API drift.

```typescript
import { useShallow } from "zustand/react/shallow";

const { name, email } = useUserStore(
  useShallow((state) => ({ name: state.name, email: state.email }))
);
```

## Constraints & Tradeoffs

- Fine-grained selectors improve performance but require selector discipline in every component.
- Feature-scoped stores reduce coupling, but can require orchestration for cross-feature workflows.
- Persisted state improves UX continuity, yet introduces versioning/migration responsibilities.
- Middleware improves ergonomics, but each layer increases cognitive and runtime cost.
- Async action abstraction simplifies UI, while adding an architectural boundary to maintain.

Equality and SSR lifecycle checks:
- In v5, `create` no longer takes an equality function. If custom equality is required, `createWithEqualityFn` from `zustand/traditional` requires `use-sync-external-store`; do not add it without approval.
- Only for an existing SSR runtime, use request-local stores, not a shared server module singleton; keep a stable browser store for the intended provider lifetime.
- Use matching initial snapshots on server and browser. Plan controlled persisted-state hydration so storage does not change the first client render unexpectedly; verify against the installed version for framework/provider and persist recipes.
- Browser-only SPAs do not need SSR scaffolding; validate persisted-state loading after reload when persistence is used.

Store review checklist:
- Ensure each component subscribes only to required state slices.
- Ensure async store actions expose explicit loading and error state.
- Ensure persisted keys are namespaced and migration-safe.
- Ensure local-only UI state remains outside global stores.
- Ensure middleware stack has clear, documented purpose.

Performance notes:
- Measure rerenders in React DevTools before and after selector changes.
- Prefer scalar selectors in hot paths over broad object selections.
- Normalize store shape when deep updates become frequent.

## Anti-Patterns

1. `const store = useStore()` in render paths that only need one field.
2. Shared monolithic store for unrelated domains.
3. Persisting tokens/session-like secrets directly in browser storage.
4. Treating either valid `useShallow` export path as an incompatible v4-only API.
5. Calling API clients directly inside many components instead of central actions.
6. Returning new arrays/objects in selectors without shallow guards.
7. Persisting loading or error flags that should reset per session.
8. Mixing Redux and Zustand in the same project without an approved architectural plan.

## Quick Pre-Merge Verification Checklist

- [ ] Applicability gate executed and documented (applies vs not applicable).
- [ ] If Redux or another solution owns this scope without an approved migration, result is **not applicable** and no forced migration was introduced.
- [ ] Store layout is per-feature/domain (no monolithic store).
- [ ] Components read state through selectors only (no full-store subscription in render paths).
- [ ] `useShallow` imports use either valid path; `zustand/react/shallow` is the convention, not a compatibility mandate.
- [ ] Zustand is not being used for server/API state; the established remote-cache solution retains ownership.
- [ ] Persist usage is justified, safe, and key names are namespaced.

## Progressive Disclosure

1. Confirm installed Zustand version and check Context7 exact-version docs before API-sensitive store refactors.
2. Start with per-feature stores plus scalar selectors.
3. Add `useShallow` where multi-field selection is required.
4. Move async workflows into feature actions and surface loading/error flags.
5. Add persistence only to user-preference state with namespaced keys.
6. Align store usage with React boundaries from [../react-19/SKILL.md](../react-19/SKILL.md).

Adoption maturity model:
- **Stage 1**: domain stores + scalar selectors.
- **Stage 2**: shallow multi-selectors and action extraction.
- **Stage 3**: controlled persistence with key governance.
- **Stage 4**: middleware tuning with performance profiling.

Execution checkpoints:
- Add selector-level tests for critical derived state behavior.
- Validate hydration behavior for each persisted field after app reload.
- Document migration policy when persisted schema changes.

## Resources

- Zustand Docs — https://zustand.docs.pmnd.rs/
- Zustand Getting Started — https://zustand.docs.pmnd.rs/getting-started/introduction
- Zustand Prevent Rerenders with useShallow — https://zustand.docs.pmnd.rs/guides/prevent-rerenders-with-use-shallow
- Zustand Persist Middleware — https://zustand.docs.pmnd.rs/integrations/persisting-store-data
- React boundary alignment — [../react-19/SKILL.md](../react-19/SKILL.md)

Related operational pairings:
- Use with [../typescript/SKILL.md](../typescript/SKILL.md) for typed store interfaces.
- Use with [../tanstack-query/SKILL.md](../tanstack-query/SKILL.md) to separate cache from global UI state.
- Use with [../security/SKILL.md](../security/SKILL.md) before persisting user-sensitive state.
- Keep store modules colocated with feature folders for ownership clarity.

## Changelog

### v2.2 — 2026-10-23
- Preserved established state solutions, corrected shallow/equality API guidance, and scoped SSR/persist hydration lifetimes.

### v2.1 — 2026-04-22
- Added mandatory Applicability Gate with explicit apply/not-apply detection criteria.
- Enforced governance decision: no default coexistence of state-management packages; respect Redux legacy as out-of-scope for this skill.
- Added required not-applicable behavior: return "not applicable" and do not force migration.
- Reinforced Zustand v5 practices (per-feature stores, selectors-only reads, correct `useShallow` path).
- Explicitly prohibited using Zustand for server/API state (use TanStack Query instead).
- Added anti-pattern for mixing Redux + Zustand without approved architecture plan.
- Added quick pre-merge verification checklist.

### v1.3 — 2026-04-21
- Standardized to operational format with required section set.
- Converted core guidance into 10 atomic Do/Don't rules.
- Reduced examples to one short snippet and preserved v5-specific import guidance.
- Added store review checklist and maturity model for adoption planning.
