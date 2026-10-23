---
name: testing
description: "Trigger: React SPA tests, component behavior, forms, API states, test selection. Verify observable behavior with existing tools; do not import Next or Ionic recipes."
license: MIT
allowed-tools: Read Write Edit Bash Grep Glob
metadata:
  author: Zesh-One
  version: "1.0"
---

# Frontend behavior testing

## When to Use

Load for test selection or behavior changes in a React SPA. Inspect installed versions, existing test runner, DOM/browser environment, form flow and authenticated client before choosing a test. This first profile does not cover Next server rendering/actions or Angular/Ionic lifecycles; retain their own runtime guidance. A testing task is not permission to migrate a stack or install tools.

## Critical Patterns

### Choose evidence, not a quota

Prefer the smallest layer that catches the actual regression; this is not a fixed pyramid or a coverage percentage mandate.

| Layer | Use for | Does not establish |
|---|---|---|
| Unit | Pure mapping, schema, calculation and boundary rules with deterministic inputs | Rendered form behavior or API authorization |
| Component | Rendered UI state and interaction through role, accessible name and visible content | Browser layout, native validation or real networking |
| Integration | Component + actual client/cache/provider wiring with controlled external responses | The real API's authentication, authorization or storage |
| E2E | A critical journey through the deployed/running app in an actual browser | Exhaustive coverage of every branch |

Reuse the installed runner and libraries. If prerequisites are absent, report which proof is unavailable and request explicit approval before adding dependencies; do not replace component/browser proof with source-text assertions. Keep any optional specimen outside default checks that cannot run it.

### Observe the user contract

- Exercise loading, error, empty and success states when the feature has them; cover retry/recovery with the previous error cleared and the expected visible result.
- Prefer role + accessible name queries and visible output. Assert required content, absence of stale/conflicting states, and the interaction's observable effect. Do not assert internal hook state, private methods, render counts, generated classes or incidental DOM structure.
- Keep the existing authenticated client -> service -> API boundary. Replace only the external seam needed by the test; do not bypass production validation/cache behavior and call that integration coverage. A minimal injected-client specimen is component acceptance, not proof of production Query wiring.
- Control async responses with deferred promises or the existing request-mocking facility; await the resulting UI and effects. Never add arbitrary sleep to hide races. Use fake clocks only for time-dependent behavior; settle pending work and restore mocks/clocks/providers during cleanup.
- Isolate per-test data, DOM and caches. If Query is actually used, create a test-owned client outside render, disable irrelevant retries, and clear it afterward without changing production retry/cache policy.

### Forms have effects, not just schemas

Preserve the installed submission flow; do not migrate forms to Server Actions or `useActionState` for testing. Exercise invalid and valid input through rendered labels/controls, not only a schema unit test.

| Case | Observable assertions |
|---|---|
| Invalid input | Useful field error, no API call, unchanged stored data/counters and usable controls |
| Valid submission | Expected client payload, pending indication, disabled duplicate submission, success result and settled controls |
| Recoverable API failure | Safe actionable error, no false success, retained/recoverable input, retry allowed and error cleared on recovery |

Client-side checks do not replace backend authorization, ownership or input validation. Use synthetic records; never live credentials, production APIs or shared user data. A test invoking a real write needs separately authorized isolated state and assertions for every observable effect.

## Constraints & Tradeoffs

A simulated DOM is not browser execution: do not claim layout, focus fidelity, navigation, native events, visual accessibility or real E2E from jsdom. Accessible queries improve behavioral intent but are not an accessibility audit. Test execution is not model compliance or an application-security guarantee. Preserve .NET API tests and the shared security controls; do not add backend infrastructure merely to test a UI.

## Progressive Disclosure

1. State the changed user rule, existing compatibility behavior and selected test layer.
2. Write the failing behavioral test before the fix where runnable; distinguish a missing tool from a meaningful RED.
3. Run the focused check and applicable existing suite; report exact outcomes, unavailable modes and cleanup. Keep instruction checks, rendered-component acceptance, integration and browser evidence separately labeled.

## Commands

In this collection, `npm run test:frontend-instructions` checks the guidance and routes only; it is included in Node-only default verification. Consumer applications use their existing test commands, not a mandated new runner.

## Anti-Patterns

Snapshot-only confidence; asserting implementation details; schema-only claims for forms; sleeps that hide unresolved requests; mocks that erase the integration under test; installing a runner without permission; treating SPA components as Next server code; declaring jsdom results to be browser/E2E or security proof.

## Resources

- Load [React guidance](../react-19/SKILL.md) only for the installed React 19 profile; preserve existing SPA forms and mutations.
- For actual Query wiring/cache tests, retain [TanStack Query](../tanstack-query/SKILL.md); no SSR hydration is required in a non-SSR SPA.
- Keep the [shared security baseline](../../shared/security/baseline.md) and its capability-specific references.
- For new structure, reuse [architecture selection](../../shared/architecture-selection.md); this skill does not choose an architecture.
