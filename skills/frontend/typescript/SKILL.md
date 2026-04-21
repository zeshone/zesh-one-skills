---
name: typescript
description: >
  TypeScript defaults for explicit typing, utility types, and schema-driven contracts.
  Trigger: When authoring or reviewing frontend TypeScript models, APIs, and shared type utilities.
license: Apache-2.0
allowed-tools: Read Edit Write Glob Grep
metadata:
  author: Zesh-One
  version: "2.0"
  inspired-by: gentleman-programming/typescript
---

## When to Use

Use this skill as the default contract when writing or reviewing frontend TypeScript:

- Feature/domain models and API response/request shapes.
- UI props and state typing.
- Shared utilities (mapped/conditional/generic types).
- Refactors where implicit inference is hiding risk.
- Reviewing pull requests that introduce `any`, broad casts, or duplicated model definitions.

## Critical Patterns

1) **DO make public function signatures explicit; DON'T rely on inferred return types for exported APIs.** Why: prevents accidental contract drift.

2) **DO use `unknown` + narrowing for untrusted inputs; DON'T use `any`.** Why: preserves type safety at boundaries.

3) **DO prefer schema-driven inference (`z.infer`) when a Zod schema exists; DON'T duplicate parallel interfaces.** Why: one source of truth.

4) **DO use utility types (`Pick`, `Omit`, `Partial`, `Readonly`, `Record`) for derivations; DON'T copy/paste near-identical models.** Why: reduces drift.

5) **DO use `import type` for type-only imports; DON'T import runtime values by mistake.** Why: cleaner bundles and clearer intent.

6) **DO model finite states via const maps + extracted unions; DON'T handwrite raw repeated string unions across files.** Why: centralizes allowed values.

7) **DO keep interfaces/composed types readable and shallow; DON'T inline deep anonymous object trees in exported contracts.** Why: improves reviewability and reuse.

8) **DO prefer discriminated unions for variant payloads; DON'T rely on optional-property guessing.** Why: gives exhaustive narrowing.

9) **DO annotate async function return types (`Promise<T>`); DON'T leak implicit `Promise<any>`.** Why: protects call sites.

10) **DO use `never` exhaustiveness checks in critical switches; DON'T leave default branches that hide missing cases.** Why: catches future enum/state additions.

```typescript
const STATUS = {
  ACTIVE: "active",
  INACTIVE: "inactive",
  PENDING: "pending",
} as const;

type Status = (typeof STATUS)[keyof typeof STATUS];
```

## Constraints & Tradeoffs

- Explicit typing increases verbosity, but pays off during refactors and API changes.
- Schema-driven typing reduces duplication, but schema ownership must be clear to avoid circular dependencies.
- Utility types reduce maintenance, but over-composition can hurt readability if stacked excessively.
- Strict no-`any` policy improves safety, but may require incremental `unknown` narrowing wrappers for third-party data.
- Exhaustiveness checks improve robustness, but require consistent discriminants in domain models.
- Explicit contracts on exports may require extra upfront type aliases, but they reduce downstream breaking changes.
- Heavy generic abstractions can look elegant but slow onboarding; prefer clarity over type-level cleverness.

## Anti-Patterns

- Exporting inferred types from implementation details without explicit contracts.
- Using `as any` to bypass compile errors instead of narrowing or redesigning types.
- Maintaining separate manual types for data already defined by runtime schemas.
- Repeating string literal unions across components instead of shared const-map derivation.
- Deeply nested inline object types in props that should be extracted and named.
- Default switch branches that swallow newly introduced variants.

## Progressive Disclosure

1. Start with explicit public contracts and `unknown` at boundaries.
2. Apply schema-driven inference where schemas exist.
3. Refactor repetitive models using utility types.
4. Introduce discriminated unions for variant behaviors.
5. Add exhaustiveness checks in high-risk state transitions.
6. Deepen with framework-specific defaults in [../react-19/SKILL.md](../react-19/SKILL.md) and schema alignment in [../zod-4/SKILL.md](../zod-4/SKILL.md).

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [TypeScript Utility Types](https://www.typescriptlang.org/docs/handbook/utility-types.html)
- [Zod Type Inference](https://zod.dev/)
- [React 19 conventions](../react-19/SKILL.md)
- [Zod 4 schema guidance](../zod-4/SKILL.md)

## Changelog

### v2.0 — 2026-04-21
- P2 trim: replaced tutorial-style examples with concise operational rules.
- Added required sections: Constraints & Tradeoffs, Anti-Patterns, Progressive Disclosure.
- Kept a single short snippet and reinforced schema-driven typing + explicit contracts.

### v1.1 — 2026-03-24
- Initial concise TypeScript defaults (const types, flat interfaces, no `any`).
