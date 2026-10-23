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
| Native runtime, plugins, platform behavior, Capacitor config/build | `capacitor` |

## Operational Defaults
- Stack baseline: Angular 19 standalone, Ionic 8, Capacitor 6, TypeScript strict.
- API access flow is fixed: component -> service (`HttpClient`) -> interceptor-authenticated request.
- Use Angular Signals (`signal`, `computed`) as state primitive for new code.
- Keep smart/presentational separation where possible for maintainability.

## Non-Negotiables
- Preserve the valid page shell of the installed Angular/Ionic integration; do not introduce wrappers copied from a different framework or rewrite valid markup for an unrelated refresh task.
- Use `ionViewWillEnter` for refreshable data loading, not `ngOnInit`.
- Use standalone Ionic imports from `@ionic/angular/standalone`.
- Use `NavController` for page navigation transitions.
- Do not programmatically jump between tab roots.
- Run `ionic cap sync` after adding/updating Capacitor plugins.
- Use `Capacitor.getPlatform()` for platform checks; never user-agent sniffing.
- Wrap Capacitor plugin calls with `try/catch` and handle cancellation separately.
- Use `ChangeDetectionStrategy.OnPush` on presentational components. Do not generalize this to every routed page or outlet ancestor; inspect state/update mechanics and installed-version lifecycle guidance first.
- Import Ionicons individually; do not load the full icon set.
- Enforce HTTPS for API traffic in production Capacitor config.
