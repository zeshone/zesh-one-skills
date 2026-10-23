---
name: capacitor
description: >
  Operational defaults for Capacitor native-layer work in app projects.
  Trigger: When implementing or modifying Capacitor plugins, platform behavior, native listeners, or mobile runtime integration.
license: Apache-2.0
allowed-tools: Read Write Edit Bash
metadata:
  author: Zesh-One
  version: "2.2"
  inspired-by: capawesome-team/ionic-capacitor
---

## When to Use

Use this skill when the task crosses from web app code into native runtime behavior.

- Installing and wiring Capacitor plugins.
- Platform checks (`web`, `ios`, `android`) and runtime branching.
- Native API calls (camera, push, filesystem, network, app lifecycle).
- Native listener lifecycle and cleanup.
- Syncing web/native layers after dependency changes.

Confirm installed Capacitor core/platform/plugin versions and supported APIs before changing native integration; the guide's Capacitor 6 baseline does not authorize an upgrade. Do not install a plugin or run native sync for a UI-only task. Install/update only for an authorized native requirement; then retain the sync and device-verification rules below.

Keep [shared security](../../shared/security/baseline.md) active, and load [delivery checks](../../shared/security/delivery-checks.md) when dependencies or native configuration change.

For Ionic page architecture and lifecycle integration, pair with [the Ionic Angular skill](../ionic-angular/SKILL.md).

## Critical Patterns

Apply these as rules, not suggestions.

1. **Do** run `ionic cap sync` after every plugin install/change.
   **Don't** assume `npm install` updates native projects.
   **Why:** Native iOS/Android projects are not auto-synced by JS dependency install.

2. **Do** use `Capacitor.getPlatform()` and `Capacitor.isNativePlatform()`.
   **Don't** use user-agent sniffing (`navigator.userAgent`).
   **Why:** UA checks are brittle and often wrong in embedded webviews.

3. **Do** wrap every native plugin call in `try/catch`.
   **Don't** collapse cancellation and operational errors into one path.
   **Why:** User cancellation is expected behavior, not a crash signal.

4. **Do** normalize plugin error handling at service boundaries: distinguish denied/restricted permission, unsupported capability, cancellation and operational failures.
   **Don't** duplicate ad-hoc string matching in multiple components or report every outcome as an unexpected error.
   **Why:** Recovery differs by outcome; error codes and restricted-state handling are plugin/platform-specific, so verify against the installed version.

5. **Do** register long-lived native listeners intentionally and store handles.
   **Don't** leak listeners across view re-entry or logout/login cycles.
   **Why:** Orphan listeners cause duplicate events, memory growth, and phantom side effects.

6. **Do** inspect the installed plugin's platform support and check runtime/plugin availability where applicable with `Capacitor.isNativePlatform()`, `Capacitor.getPlatform()` and `Capacitor.isPluginAvailable('<Name>')` before execution.
   **Don't** equate native-platform detection with plugin availability or permission. Do not blanket-disable web: many official plugins ship web implementations.
   **Why:** Platform, plugin availability and permission are separate checks; provide a supported web path or an explicit unsupported-capability fallback.

7. **Do** initialize push registration when user context is available.
   **Don't** register blindly before auth/session restoration.
   **Why:** Token association and deep-link routing usually depend on user identity.

8. **Do** treat filesystem persistence as data-domain design (path, retention, cleanup).
   **Don't** use Preferences for large payloads or file-like data.
   **Why:** Preferences are key-value settings storage, not document storage.

9. **Do** test plugin flows on real devices for release-critical features.
   **Don't** trust browser-only validation for native behavior.
   **Why:** Permissions, OS dialogs, and hardware quirks differ from emulators/browsers.

10. **Do** keep config changes explicit in `capacitor.config.*` and review security-impacting flags.
    **Don't** enable permissive flags (cleartext/mixed content) outside controlled dev scenarios.
    **Why:** Misconfigured transport/security flags can silently ship high-risk builds.

11. **Do** follow each permission-based plugin's own docs: use `checkPermissions()`/`requestPermissions()` where exposed, with states such as `granted`, `denied`, `prompt`, `prompt-with-rationale`.
    **Don't** assume all plugins share permission aliases, request timing or recovery paths; verify against the installed version. Request only permissions needed for the user action and handle denial without repeated prompt loops.
    **Why:** Runtime permission handling does not replace native configuration: add required iOS Info.plist usage descriptions and Android manifest permissions according to that plugin's own docs. Explain any settings-based recovery only when supported by the target platform/plugin.

## Constraints & Tradeoffs

- **Constraint:** `ionic cap sync` after plugin changes is mandatory.
  **Tradeoff:** Adds iteration cost, but prevents mismatched native/JS state.

- **Constraint:** Platform detection must use Capacitor APIs; execution also requires plugin support and applicable permissions.
  **Tradeoff:** More explicit checks, but valid web implementations remain usable rather than being rejected by a native-only gate.

- **Constraint:** Native calls require `try/catch` and cancel/error separation.
  **Tradeoff:** More boilerplate; dramatically better UX and observability.

- **Constraint:** Listener cleanup is required.
  **Tradeoff:** Requires lifecycle discipline and handle management.

- **Constraint:** Skill remains concise operational guidance.
  **Tradeoff:** Deep plugin API coverage belongs in official docs, not here.

## Anti-Patterns

- Installing plugin packages without running `ionic cap sync`.
- UA sniffing for iOS/Android detection.
- Treating native detection or plugin availability as proof that permission is granted.
- Disabling a supported web implementation solely because the runtime is not native.
- Single generic catch block showing “unexpected error” for user cancellation.
- Creating listeners on every page enter without removing previous handles.
- Plugin calls from UI components without service-level guards and mapping.
- Using Preferences as a pseudo-database for large JSON blobs.
- Enabling insecure config flags and forgetting to revert for release.
- Assuming emulator success equals production readiness.
- Embedding long plugin tutorials in this skill file.

## Progressive Disclosure

Implement only the layer you need.

1. **Baseline:** confirm installed versions, plugin/platform support, availability and applicable permission workflow; install/change a plugin only when required and authorized, then sync native projects.
2. **Error layer:** add `try/catch` and distinct permission-denied/restricted, unsupported, cancellation and operational failure paths.
3. **Lifecycle layer:** register listeners once, keep handles, remove on teardown.
4. **Security/config layer:** review `capacitor.config.*` flags for target environment.
5. **Production hardening:** verify on real devices with permission-denied/cancel/network-off scenarios, plus first request, granted, denied/restricted and supported-web/unsupported fallbacks as applicable; verify against the installed version.

Short non-obvious reference snippet:

```typescript
import { Capacitor } from '@capacitor/core';

const platform = Capacitor.getPlatform(); // 'ios' | 'android' | 'web'
const isNative = Capacitor.isNativePlatform();
```

## Resources

- Shared security: [baseline](../../shared/security/baseline.md) and [delivery checks](../../shared/security/delivery-checks.md)
- Capacitor docs: https://capacitorjs.com/docs
- Capacitor config reference: https://capacitorjs.com/docs/config
- Capacitor plugin APIs index: https://capacitorjs.com/docs/apis
- Companion Ionic app guidance: ../ionic-angular/SKILL.md

## Changelog

### v2.2 — 2026-10-23
- Separated platform/plugin capability, supported web paths, plugin-specific permissions/configuration and recoverable outcomes.

### v2.1 — 2026-10-22
- Added installed-version and native-change gates, avoiding plugin installation/sync for UI-only work while preserving native protections.

### v2.0 — 2026-04-21
- Rewritten into concise operational guidance with explicit sync, platform, error, and lifecycle defaults.
