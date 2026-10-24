# Apps Domain Agent Guide (Angular 19 + Ionic 8 + Capacitor 6)

## Scope
Inspect installed versions of Angular, Ionic, Capacitor and plugins, plus existing routing/forms/data-access conventions. The named stack is a baseline, not an upgrade instruction. Keep valid Reactive Forms and HttpClient flows; do not adopt newer APIs without checking support. If a domain skill conflicts with this file, the skill wins within its installed-version scope.

## Always-on security baseline

Load [the shared security baseline](skills/shared/security/baseline.md) for every task. On changed inputs, credentials or authorization, load matching [application checks](skills/shared/security/application-checks.md); on dependencies/native config, load [delivery checks](skills/shared/security/delivery-checks.md). Preserve the existing authenticated client, server-side authorization and sensitive-data safeguards. Do not import Next.js middleware/cookie recipes or .NET server secret-storage rules into a mobile client. Client bundles and Angular environment files cannot keep secrets private.

Include `skills/shared/security/` with the app skills and adjust paths when installed elsewhere. The runtime must load the references; links alone do not enforce them.

## Mandatory Skill Routing (load first)

| Work Type | Skill |
|---|---|
| Ionic page/component architecture, routing, lifecycle, forms, signals | `ionic-angular` |
| Bug reproduction, diagnosis, regression investigation | [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md) |
| UI accessibility, keyboard/focus, labels and feedback | [accessibility](skills/shared/accessibility/SKILL.md) |
| Native runtime, plugins, platform behavior, Capacitor config/build | `capacitor` |

## Operational Defaults
- Stack baseline: Angular 19 (standalone by default), Ionic 8, Capacitor 6, TypeScript strict. Preserve valid NgModule and RxJS ownership; no migration is required.
- API access flow is fixed: component -> service (`HttpClient`) -> interceptor-authenticated request.
- Prefer Angular signals where appropriate (`signal`, `computed`); do not replace valid RxJS state merely for consistency.
- Keep smart/presentational separation where possible for maintainability.

## Non-Negotiables
- Preserve the valid page shell of the installed Angular/Ionic integration; do not introduce wrappers copied from a different framework or rewrite valid markup for an unrelated refresh task.
- Use `ionViewWillEnter` for refreshable or stale data loading, not `ngOnInit` for revisits; skip fetching when cached data is still valid. Stop hidden-page work on leave, not only on destruction.
- Use standalone Ionic imports from `@ionic/angular/standalone` and `loadComponent` for standalone routes; preserve valid NgModule imports and `loadChildren` routes.
- Prefer `NavController` (`navigateForward`, `navigateBack`, `navigateRoot`) for explicit Ionic direction/stack animation; Angular `Router` and `routerLink` with `routerDirection` remain valid.
- Each tab keeps its own stack: do not push another tab's child pages into the current tab's history. Intentional tab-root selection from a notification or deep link is valid.
- Run `ionic cap sync` after adding/updating Capacitor plugins.
- Use `Capacitor.getPlatform()` for platform checks; never user-agent sniffing.
- Wrap Capacitor plugin calls with `try/catch` and handle cancellation separately.
- Use `ChangeDetectionStrategy.OnPush` on presentational components. Do not generalize this to every routed page or outlet ancestor; inspect state/update mechanics and installed-version lifecycle guidance first.
- Import Ionicons individually; do not load the full icon set.
- Enforce HTTPS for API traffic in production Capacitor config.
