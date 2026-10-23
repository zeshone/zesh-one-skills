# Phase 3 — bounded proposal evidence
These synthetic response-only observations assess proposals, not execution, accessibility compliance or enforced agent behavior.
Prompts and per-case base loads: [frozen fixture](../../tools/fixtures/debugging-accessibility-proposals.json); scope: [task document](../tasks/debugging-accessibility-adoption.md). Baseline: `e50ea92`.

## Method and rubric
One BEFORE session answered three cases before either new skill existed. AFTER uses the identical prompts/base loads plus the applicable new skill through its conditional route; observers do not read previous answers or this rubric. Cases share a session, so they are not statistically independent. Model identity was unavailable; no causal improvement claim is made.
Evaluate debugging proposals for reproduction, observations versus hypotheses, a discriminating experiment, smallest correction and result/state regression checks. Evaluate UI proposals for semantic controls/names, keyboard/focus, non-color error feedback, retained runtime flow and explicit browser/assistive-technology evidence limits. These criteria are parent interpretations of the requested gaps, not additional user requirements.

## BEFORE — `mvljbl33-2-y9lm`
Tool reported: `functions.read`; model unknown. No application commands, edits or installations.

| Case | Observed proposal | Limits |
|---|---|---|
| rejected-capacity | "Mutation before rejection is a hypothesis, not an established cause"; request source, assert false and Total=8, independently exercise amounts 1/2/3; smallest evidence-based correction | Strong baseline; no explicit predicted discriminating experiment; no source inspected or bug executed |
| react-dialog | Native/existing-library button, associated input label, supported dialog wiring, initial/contained/restored focus and Escape; separate actual-browser verification | Strong baseline already; role-based jsdom pass explicitly not browser focus proof |
| ionic-focus | Inspect retained page/lifecycle/return path; meaningful return target after view readiness; associated non-color error text; keyboard/assistive-technology checks | Strong baseline already; cause and runtime outcomes unverified |

## T1 — debugging instructions
Observed RED: `node tools/test-debugging-instructions.js` exited 1, 0/7 passed because the skill, routes and script were absent. GREEN: `npm run test:debugging-instructions`, exit 0, 7/7. Full `npm run verify`, exit 0: 91 assertions, 34 skills, zero lint errors/warnings; `git diff --check` passed.

AFTER `mvljfdil-3-ew1n`, `functions.read`, model unknown, loaded the same rejected-capacity guides plus `diagnosing-bugs`. It proposed a predicted trace: "a write to 11 before returning false", distinguishing mutation-before-validation from another writer/accessor; retained fresh-instance amounts 1/2/3, false/unchanged-state regression and smallest supported correction. Baseline already covered most criteria; AFTER made the experiment/prediction explicit, not a proven causal gain. It inadvertently read the task Log (method rationale), but not BEFORE answers or this rubric: limited blinding, disclosed rather than treated as independent proof. No bug was executed.
Independent verifier `mvljhrve-4-zuez`: T1/S1–S3 PASS, full verify91/34skills/0lint and both diff checks exit0; status/index/HEAD unchanged. Initial ASSESS could not assess undeclared untracked paths and conservatively required this verifier. After explicit staging, native ASSESS: medium, 15paths/179changed lines, reviewDue=false/under_budget; no native START or consent requested.

Rollback boundary: new debugging skill, its four route rows, focused script and test, discovery/provenance additions and this unit's fixture/evidence. Existing tests and profile/security owners remain intact. No consumer runtime boundary was changed; functional bug execution is N/A for this instruction-only unit.

## T2 — accessibility instructions
Observed RED: `node tools/test-accessibility-instructions.js`, exit1, 0/8. GREEN: `npm run test:accessibility-instructions`, exit0, 8/8. Full `npm run verify`, exit0: 99 assertions, 35 skills, zero lint errors/warnings. No browser, screen reader, native runtime or application audit has been performed.

AFTER `mvljn2c0-5-intq`, `functions.read`, model unknown, loaded each UI case's base guides plus `accessibility`. React retained the library, semantic button/label, modal focus/close and separate browser/assistive-technology checks. Ionic retained its shell/forms/client/navigation and added explicit transition-ready restoration, non-color associated feedback, browser versus WebView/device evidence and no page-wide focus trap. Both BEFORE proposals already met most criteria; no measurable causal gain is established. The observer also inadvertently read Log despite instructions, limiting blinding; no re-run or independent-proof claim.

Declared-candidate native ASSESS: medium, 12paths/132changed lines, reviewDue=false/under_budget, large writer; self-verification sufficient, no independent verifier or native START required for T2. Both staged/unstaged diff checks passed. No existing consumed contract was broken: routes/checks are additive and profile/security owners are retained.

Rollback boundary: accessibility skill, frontend/apps route rows, focused script/test and related discovery/provenance docs; T1 remains independently usable. Functional browser execution is N/A for this instruction-only change; actual consumer verification is explicitly left to consumer adoption, not marked passed.
