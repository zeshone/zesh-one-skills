# Phase 4 — framework enrichment evidence
These checks validate instruction text and routes, not application, device, browser or model behavior. Scope: [task document](../tasks/framework-enrichment-adoption.md). Baseline: `e75f575`.

## Method
Two read-only Pi explorers mapped gaps in the existing React-side and Ionic/Capacitor skills; the parent spot-checked quoted lines and verified disputed API claims upstream (context7: Zustand v5 migration guide, Zod 4 params API). Each unit then ran test-first: a new Node instruction-contract script observed RED before the skill edits and GREEN after, followed by the full `npm run verify`. No BEFORE/AFTER response observers were run in this phase: the units correct or scope existing rules rather than fill a missing capability, so no behavioral-gain claim is made.

## T1 — React stack instructions
Writer: headless Pi, edit surfaces `skills/frontend/{react-19,nextjs-15,tanstack-query,zustand-5,typescript,zod-4,tailwind-4,shadcn-ui}/SKILL.md`, `AGENTS.frontend.md`, `tools/test-react-stack-instructions.js`.
Observed RED: `node tools/test-react-stack-instructions.js`, exit 1, 14 failed / 1 preservation passed. GREEN: exit 0, 15/15. Writer full `npm run verify`, exit 0 (one transient failure while the parent integrated T2, passed on retry).
Parent review found one unrequested regression: Zod rule 1 had been weakened from the correct Zod 4 preference for top-level format validators to a generic "verify syntax" rule. Parent added a RED assertion (exit 1, R03 failed), restored the rule with the deprecated-but-working chained form, GREEN 15/15. Changelog dates corrected to the commit date. Verified upstream in this phase: Zustand v5 `useShallow` paths and Zod 4 params (see task Log L3); the other version facts given to the writer (React `useFormState` deprecation, Query per-request server client and `cancelQueries` optimism, Next 15 uncached defaults, Tailwind 4 variable shorthand) are stated as established upstream behavior and the skills still tell agents to verify against the installed version. No React, Next.js or browser execution was run.

## T2 — apps instructions
Writer: headless Pi (`openai-codex/gpt-6.1-sol`, lean flags), edit surfaces `skills/apps/{ionic-angular,capacitor}/SKILL.md`, `AGENTS.apps.md`, `tools/test-apps-instructions.js`.
Observed RED: `node tools/test-apps-instructions.js`, exit 1, 6 failed / 1 preservation passed. GREEN: exit 0, 7/7. Writer full `npm run verify`, exit 0: 99 assertions, 35 skills, zero lint errors/warnings; `git diff --check` passed.
Parent review: diff confined to the four surfaces; changelog dates corrected to the commit date; script registered in `test` and `verify`; ADOPTING and CHANGELOG updated. Removed unconditional mandates ("Standalone-first is mandatory", "refresh is mandatory on routed pages", blanket NavController/tab-root rules); added Capacitor platform/plugin/permission separation and a routine behavior checklist using installed tooling. No device, emulator or Angular build was run.

## T3 — maintainer guide
Writer: headless Pi, edit surfaces `docs/MAINTAINING.md`, `README.md`, `CREDITS.md`, `tools/test-maintainer-instructions.js`.
Observed RED: `node tools/test-maintainer-instructions.js`, exit 1, 15 failed / 2 passed (the dynamic parity check also caught the not-yet-registered React script). GREEN: exit 0, 17/17. Writer full `npm run verify`, exit 0.
Parent review: guide claims spot-checked against the files (shared skills declare MIT; the trial READMEs exist; lint coverage matches `tools/lint-skills.js`). Parent added a RED assertion and guide line for the exact pipeline list pinned by `tools/test-cli-instructions.js`, removed the parity test's self-exclusion after registering the script (RED 3 failed → GREEN 17/17), and rewrote the CREDITS section in the file's Spanish without URLs the writer had derived from names but not verified. Lint code is unchanged (S6).
