---
name: ionic-angular
description: >
  Operational defaults for Ionic + Angular application work.
  Trigger: When implementing or modifying Ionic Angular pages, routing, state, navigation, or lifecycle behavior.
license: Apache-2.0
allowed-tools: Read Write Edit Bash
metadata:
  author: Zesh-One
  version: "2.2"
  inspired-by: capawesome-team/ionic-angular, ionic-capacitor/angular.md
---

## When to Use

Use this skill for APP-level Ionic Angular work, not for generic Angular tutorials.

- Routed pages and tab flows.
- Navigation transitions and back-stack behavior.
- Page lifecycle data refresh and cached-view behavior.
- App state with Angular signals or existing RxJS flows.
- Standalone or NgModule Ionic architecture decisions.

Inspect installed versions of Angular, Ionic and Capacitor and the existing page/router/form setup before applying an example. Do not upgrade the stack to match latest documentation. Preserve valid Reactive Forms and HttpClient usage; signals for state do not require migrating forms or replacing the HTTP client. Verify version-sensitive APIs against the installed versions, not a newer Angular release.

Keep the [shared security baseline](../../shared/security/baseline.md) active; load matching [application checks](../../shared/security/application-checks.md) for inputs, credentials or access-control changes. This does not import Next.js/server-only authentication infrastructure.

If your task is mostly native plugin behavior, also load [the Capacitor skill](../capacitor/SKILL.md).

## Critical Patterns

Follow these as atomic rules. Each rule has intent (`Do/Don't`) and rationale (`Why`).

1. **Do** use standalone Ionic imports (`@ionic/angular/standalone`) in standalone pages/components.
   **Don't** force standalone conversion in an existing NgModule app.
   **Why:** Angular 19 makes components standalone by default; NgModule-based apps remain supported and no migration is required. Preserve valid NgModule and RxJS ownership.

2. **Do** preserve the valid page shell used by the installed Angular/Ionic integration.
   **Don't** introduce a wrapper from another framework or rewrite valid page markup for a refresh-only task; verify template support first.
   **Why:** Transitions and lifecycle depend on the actual integration, not a copied wrapper name.

3. **Do** refresh refreshable or stale data in `ionViewWillEnter` according to the data's freshness policy.
   **Don't** rely on `ngOnInit` for revisits. Do not fetch on every entry when cached data is still valid.
   **Why:** Ionic 8's `ion-router-outlet` retains visited pages in the DOM. `ngOnInit` runs once per page instance; `ngOnDestroy` runs only when the page is removed from the stack. Use `ionViewWillEnter`/`ionViewDidEnter` for revisit work and `ionViewWillLeave`/`ionViewDidLeave` for leave work (interfaces `ViewWillEnter`, `ViewDidEnter`, `ViewWillLeave`, `ViewDidLeave`). Stop work that must not run while hidden on leave, including view-scoped subscriptions, polling and listeners; destruction-only cleanup is insufficient. Restart intentionally on entry and retain final teardown on destruction.

4. **Do** prefer `NavController` (`navigateForward`, `navigateBack`, `navigateRoot`) when controlling Ionic direction and stack animation.
   **Don't** treat Angular `Router` or `routerLink` with `routerDirection` as invalid alternatives.
   **Why:** Navigation intent determines direction/stack handling.

5. **Do** prefer Angular signals (`signal`, `computed`, readonly exposure) where appropriate for app/session state.
   **Don't** spread mutable state across unrelated services and components or migrate working RxJS state just for consistency.
   **Why:** Explicit ownership matters more than changing reactive primitives.

6. **Do** keep presentational components dumb and `OnPush`.
   **Don't** make UI leaf components fetch or mutate domain state, or blindly apply that change-detection rule to every routed page/outlet ancestor.
   **Why:** Presentation, page-entry refresh and change detection are separate concerns. Inspect `ion-router-outlet`/app-shell ancestry, installed-version behavior and whether page state uses signals or `markForCheck`. Do not import newer Angular default-change-detection assumptions into Angular 19.

7. **Do** lazy-load standalone routes with `loadComponent`; preserve valid NgModule lazy routes with `loadChildren`.
   **Don't** eagerly import page components into route roots or convert module routing just to use `loadComponent`.
   **Why:** Both supported routing styles can control startup cost and mobile bundle pressure.

8. **Do** align tab button `tab` values with child route paths exactly.
   **Don't** depend on aliases or mismatched names.
   **Why:** Tab router outlets resolve by exact key/path matching.

9. **Do** keep each tab's own navigation stack. Intentionally selecting a tab root from a notification or deep link is valid.
   **Don't** push another tab's child pages into the current tab's history.
   **Why:** Each tab keeps its own stack so back navigation preserves that tab's history.

10. **Do** register only needed Ionicons explicitly in standalone pages/components.
    **Don't** assume global icon availability.
    **Why:** Standalone usage requires explicit registration and avoids accidental bundle growth.

11. **Do** place page-entry fetch logic behind idempotent methods (safe on every enter).
    **Don't** couple data load to constructor side effects.
    **Why:** Entry hooks may run many times; deterministic loaders avoid stale/multi-subscribe bugs.

12. **Do** follow supported project template conventions; built-in control flow (`@if`, `@for`, `@switch`) is stable since Angular 17 and `@for` requires `track` by stable entity identity.
    **Don't** migrate working structural directives merely for consistency or track mutable lists by unstable positions.
    **Why:** Existing structural directives still work; stable identity preserves row state during updates. For other version-sensitive template APIs, verify against the installed version.

13. **Do** use `inject()` only in an injection context (constructor, field initializer, provider factory) or through `runInInjectionContext`.
    **Don't** call it directly from lifecycle/event methods or replace valid constructor injection merely for consistency.
    **Why:** Dependency resolution requires an injection context; for version-specific DI APIs, verify against the installed version.

## Constraints & Tradeoffs

- **Constraint:** Match the existing standalone/NgModule architecture.
  **Tradeoff:** Standalone examples need adaptation for module-owned imports and routes, not an architectural migration.

- **Constraint:** Entry refresh follows freshness policy; hidden-page work stops on leave.
  **Tradeoff:** Cached pages need separate entry/leave cleanup and final destruction handling, not unconditional fetching.

- **Constraint:** Prefer NavController for explicit Ionic direction/stack control.
  **Tradeoff:** Router APIs remain valid; choose based on navigation intent and preserve independent tab stacks.

- **Constraint:** Prefer signals where appropriate without replacing valid RxJS ownership.
  **Tradeoff:** Interop with RxJS-heavy services may require bridge points (`toSignal`, explicit effect boundaries); verify against the installed version.

- **Constraint:** Skill content must remain operational, not encyclopedic.
  **Tradeoff:** Keep only decision-driving guidance; move deep explanation to external docs.

## Anti-Patterns

- Copying page-shell markup without checking support in the installed Angular/Ionic integration.
- Data load in constructor/`ngOnInit` expecting automatic revisit refresh.
- Ignoring Ionic direction/stack intent when choosing navigation APIs.
- Two-way mutable shared state without clear ownership boundaries.
- Destruction-only cleanup for work that must stop while a cached page is hidden.
- Presentational component importing HTTP services directly.
- Route eager imports that bloat initial chunk.
- Tab buttons whose `tab` key does not match route path.
- Copy-pasted NgModule setup into standalone app root.
- Shipping long onboarding prose inside this skill file.

## Progressive Disclosure

Start minimal, then expand only if the task demands it.

1. **Baseline:** installed-version gate, existing valid page shell, freshness-aware entry refresh (`ionViewWillEnter`), leave cleanup and navigation intent.
2. **State step:** clarify page/app state ownership and derived values; use signals where appropriate without migrating valid RxJS flows.
3. **Architecture step:** split smart/presentational components and enforce OnPush at UI edges.
4. **Performance step:** audit lazy routes, icon registration scope, and repeated entry fetch cost.
5. **Advanced only when needed:** tabs edge cases, modal routing nuances, deep-link restoration.

**Routine behavior checks:** use installed testing tooling for changed routes/data flows; do not add a testing stack merely to follow this skill.
- Check initial entry and revisit, including fresh-cache reuse and leave/re-entry cleanup without duplicate work.
- Check loading/empty/error/retry states and stale request races; an older response must not overwrite newer route/data state.
- Check back/tab behavior, independent stacks and intentional deep-link tab selection.
- Check retained form state across cached-page revisits and explicit reset/save behavior.
- Keep native/device checks in the Capacitor skill; report executed evidence and unperformed checks separately.

Short non-obvious reference snippet (the freshness policy belongs to the app):

```typescript
import { ViewWillEnter } from '@ionic/angular/standalone';

export class OrdersPage implements ViewWillEnter {
  ionViewWillEnter(): void {
    this.refreshOrdersIfStale();
  }
}
```

## Resources

- Shared security: [baseline](../../shared/security/baseline.md) and [application checks](../../shared/security/application-checks.md)
- Ionic Angular lifecycle: https://ionicframework.com/docs/angular/lifecycle
- Ionic Angular navigation: https://ionicframework.com/docs/angular/navigation
- Angular signals guide: https://angular.dev/guide/signals
- Companion native guidance: ../capacitor/SKILL.md

## Changelog

### v2.2 — 2026-10-23
- Preserved NgModule/RxJS ownership; bounded cached-page work, navigation and Angular idioms; added routine behavior checks.

### v2.1 — 2026-10-22
- Added installed-version/forms/HttpClient gates and shared security routing; preserved revisit hooks and native navigation conventions.
- Kept page-shell and change-detection decisions specific to the installed Angular/Ionic integration rather than transplanting other-framework/newer-version rules.

### v2.0 — 2026-04-21
- Rewritten into concise operational guidance for Ionic Angular application work.
- Enforced standalone imports, `ionViewWillEnter` refresh, NavController transitions, and signals-first state defaults.
