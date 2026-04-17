# Apps Domain — ZeshOne Mobile Development Guidelines

## How to Use This Guide

- Start here for cross-project norms on mobile app development.
- Each skill has detailed patterns available on-demand.
- Domain skills override this file when guidance conflicts.

---

## Stack Base

| Layer | Technology |
|-------|-----------|
| Framework | Angular 19 (Standalone components) |
| UI + Mobile shell | Ionic 8 (`@ionic/angular`) |
| Native runtime | Capacitor 6 |
| Language | TypeScript 5 (strict mode) |
| State | Angular Signals + `computed()` |
| HTTP | Angular `HttpClient` + functional interceptors |
| Forms | Angular Reactive Forms |
| Icons | Ionicons 7 (import individually) |
| Local storage | `@capacitor/preferences` |

**API pattern:** All calls go via Angular `HttpClient` services with an `authInterceptor` injecting the Bearer token. Never call APIs directly from components.

---

## Available Skills

### Apps Skills (ZeshOne)

| Skill | Description | Path |
|-------|-------------|------|
| `ionic-angular` | Standalone components, page structure, lifecycle hooks, navigation (tabs/menu/modal), forms, signals, services, smart/presentational split, performance | [SKILL.md](skills/apps/ionic-angular/SKILL.md) |
| `capacitor` | Capacitor config, platform detection, native plugins (camera, push, preferences, filesystem, geolocation, haptics), theming, dev workflow, build/deploy | [SKILL.md](skills/apps/capacitor/SKILL.md) |

---

## Auto-invoke Skills

When performing these actions, **ALWAYS** invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| Creating any Ionic page or component | `ionic-angular` |
| Setting up routing, tabs, side menu, guards | `ionic-angular` |
| Data fetching in pages (services, HTTP) | `ionic-angular` |
| Managing state with signals | `ionic-angular` |
| Building forms with Ionic inputs | `ionic-angular` |
| Using any Capacitor plugin | `capacitor` |
| Platform detection (iOS vs Android vs web) | `capacitor` |
| Handling push notifications | `capacitor` |
| Accessing camera, filesystem, preferences | `capacitor` |
| Configuring `capacitor.config.ts` | `capacitor` |
| Building for iOS or Android | `capacitor` |
| Theming and dark mode | `capacitor` |

---

## Critical Cross-Skill Rules

- **Every routed page must be wrapped in `<ion-page>`** — without it, transitions break and safe areas fail.
- **Use `ionViewWillEnter` for data fetching**, never `ngOnInit` — Ionic caches pages and `ngOnInit` won't re-fire on back-navigation.
- **Architecture is Standalone** — import Ionic components from `@ionic/angular/standalone`, call `addIcons()` in constructor.
- **Signals are the state primitive** — prefer `signal()` + `computed()` over `BehaviorSubject` for new code.
- **Use `NavController`** (not Angular `Router`) for page navigation — it handles Ionic transition animations.
- **Never navigate between tabs programmatically** — tabs have independent stacks; only the tab bar switches tabs.
- **`ionic cap sync` after every plugin install** — `npm install` alone doesn't update native projects.
- **Use `Capacitor.getPlatform()`** for platform detection — never user-agent sniffing.
- **Always wrap Capacitor calls in try/catch** — distinguish user cancellation from real errors.
- **`ChangeDetectionStrategy.OnPush`** on all presentational components.
- **Icons imported individually** — `import { heart } from 'ionicons/icons'`, never the full library.
- **All API calls use HTTPS** — `allowMixedContent: false` in `capacitor.config.ts` (production).

---

## Lifecycle Reference (Quick)

| Hook | When | Use |
|------|------|-----|
| `ngOnInit` | Once | One-time setup |
| `ionViewWillEnter` | Every visit | **Fetch/refresh data** |
| `ionViewDidEnter` | Every visit (after animation) | Heavy work, analytics |
| `ionViewWillLeave` | Before leaving | Save drafts, kill subscriptions |
| `ngOnDestroy` | Page popped from stack | Final cleanup |
