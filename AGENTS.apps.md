# Apps Domain Agent Guide (Angular 19 + Ionic 8 + Capacitor 6)

## Scope
Use this file as the execution default for mobile app work in this repo. Keep outputs practical, lifecycle-safe, and production-oriented. If a domain skill conflicts with this file, the skill wins.

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
- Every routed page must include `<ion-page>` wrapper.
- Use `ionViewWillEnter` for refreshable data loading, not `ngOnInit`.
- Use standalone Ionic imports from `@ionic/angular/standalone`.
- Use `NavController` for page navigation transitions.
- Do not programmatically jump between tab roots.
- Run `ionic cap sync` after adding/updating Capacitor plugins.
- Use `Capacitor.getPlatform()` for platform checks; never user-agent sniffing.
- Wrap Capacitor plugin calls with `try/catch` and handle cancellation separately.
- Use `ChangeDetectionStrategy.OnPush` on presentational components.
- Import Ionicons individually; do not load the full icon set.
- Enforce HTTPS for API traffic in production Capacitor config.
