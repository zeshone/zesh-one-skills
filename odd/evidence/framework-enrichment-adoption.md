# Phase 4 — framework enrichment evidence
These checks validate instruction text and routes, not application, device, browser or model behavior. Scope: [task document](../tasks/framework-enrichment-adoption.md). Baseline: `e75f575`.

## Method
Two read-only Pi explorers mapped gaps in the existing React-side and Ionic/Capacitor skills; the parent spot-checked quoted lines and verified disputed API claims upstream (context7: Zustand v5 migration guide, Zod 4 params API). Each unit then ran test-first: a new Node instruction-contract script observed RED before the skill edits and GREEN after, followed by the full `npm run verify`. No BEFORE/AFTER response observers were run in this phase: the units correct or scope existing rules rather than fill a missing capability, so no behavioral-gain claim is made.

## T2 — apps instructions
Writer: headless Pi (`openai-codex/gpt-6.1-sol`, lean flags), edit surfaces `skills/apps/{ionic-angular,capacitor}/SKILL.md`, `AGENTS.apps.md`, `tools/test-apps-instructions.js`.
Observed RED: `node tools/test-apps-instructions.js`, exit 1, 6 failed / 1 preservation passed. GREEN: exit 0, 7/7. Writer full `npm run verify`, exit 0: 99 assertions, 35 skills, zero lint errors/warnings; `git diff --check` passed.
Parent review: diff confined to the four surfaces; changelog dates corrected to the commit date; script registered in `test` and `verify`; ADOPTING and CHANGELOG updated. Removed unconditional mandates ("Standalone-first is mandatory", "refresh is mandatory on routed pages", blanket NavController/tab-root rules); added Capacitor platform/plugin/permission separation and a routine behavior checklist using installed tooling. No device, emulator or Angular build was run.
