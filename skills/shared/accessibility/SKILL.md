---
name: accessibility
description: "Trigger: accessibility, a11y, keyboard, focus, labels, UI feedback. Improve changed interfaces and report verification limits."
license: MIT
allowed-tools: Read Write Edit Bash Grep Glob
metadata:
  author: Zesh-One
  version: "1.0"
---

# Accessible interactions

## When to Use

Load for changed web or hybrid UI interactions in React SPA, Next or Angular/Ionic, within installed versions. This is not a full audit or a framework migration. Inspect the existing UI library and preserve valid forms, navigation, lifecycle and authenticated-client flows.

## Critical Patterns

1. **Semantics first.** Prefer native controls, headings and landmarks. Give each control an accessible name and inputs associated labels; placeholders alone are not labels. Reuse supported library primitives; ARIA augments semantics only when the required interaction behavior is implemented, not as a replacement for a working button.
2. **Keyboard and focus.** Verify keyboard activation, logical order and visible focus without positive tabindex shortcuts. For a modal dialog, use supported initial focus, containment, Escape/close and restore focus to its trigger or a meaningful surviving target. Do not trap focus in ordinary pages. On navigation/revisit, inspect framework transition readiness before moving focus; avoid arbitrary delays.
3. **Understandable feedback.** Associate error text and invalid state with the affected input; use text or another distinguishable cue, not color alone. Make relevant asynchronous status announcements discoverable without stealing focus or announcing sensitive content. Preserve recovery, loading and disabled states using the existing form flow.
4. **Perceivable presentation.** Check text/control contrast, zoom, reflow and reduced motion in the actual layout; do not rely on snapshots or CSS class names. Identify the target standard/version/level if conformance is requested; a local check is not a compliance guarantee.
5. **Evidence before claims.** Use existing tooling for role/name/state and interaction regressions. Validate keyboard, focus and layout in an actual browser; use relevant assistive technology for names, reading order and announcements. jsdom and automated scans cover only part of accessibility, not full conformance or native-device behavior. Report unperformed checks instead of treating them as passes.

## Constraints & Tradeoffs

Keep changes proportional to the affected journey. Do not require shadcn or a new UI/test stack; no native plugin is needed merely to inspect UI semantics. Missing browser/device/assistive-technology tools require disclosure and explicit approval before installation. Keep the security baseline active.

## Anti-Patterns

Clickable divs without equivalent interaction; ARIA-only repairs; hidden focus; color-only errors; broad markup rewrites; declaring accessibility from a green role query or scanner score.

## Progressive Disclosure

Report the affected journey and observed barrier, smallest correction, exact checks/environment and remaining manual/browser/device checks. Distinguish proposed checks from executed evidence; do not claim verified accessibility for an untested application.

## Resources

- Retain the [security baseline](../security/baseline.md) and [evidence-led diagnosis](../diagnosing-bugs/SKILL.md) for reported regressions.
- Reuse [React SPA testing](../../frontend/testing/SKILL.md) only where applicable; Angular/Ionic and Next keep their installed testing flow.
