# Architecture adoption pilot evidence

Scope: phases 0–1; existing security enrichment is the baseline, not a separate optional add-on. No later catalog additions are included.

## Reproduction inputs and limits

- BEFORE source: `52ccba7b5b05a085ac0164fadcd959a1cbb9a055`.
- Fixed prompts and loading instructions: [six response-only cases](../../tools/fixtures/architecture-adoption-pilot.json).
- BEFORE observation: `gentle-ai-explore`, task `mvjzjsoc-4-s8wv`, foreground, one additional Pi session.
- AFTER source: `43c3d4de1886c998c358e4397a02c1e2cd290453`; fresh `gentle-ai-explore` task `mvk4ew8o-8-0vvw`, foreground, one additional Pi session.
- Both phases used the same agent role and requested configuration, with unchanged fixture prompts. Each phase answered six cases in one batch, not six isolated trials. The host did not expose resolved model/effort identities; this is not a controlled model comparison or causal benchmark.
- These are observed model proposals, not scaffold execution tests. Read-only tool restrictions mean this exercise cannot demonstrate that skill prose alone prevents writes, sandbox escape, secret exposure, or dependency installation.
- No application was created or modified. Runtime loading and behavior in another consumer remain unverified. Final verification examines documents and these observations separately; text presence does not establish model compliance.

## BEFORE observations

| Case | Observed response | Interpretation for this pilot |
|---|---|---|
| `cli-small` | Proposed one console project plus tests, but reported REST-only scope and the `Console.WriteLine` prohibition as conflicts. | Simple recommendation already possible through model judgment; instructions still need an explicit CLI profile and product-output distinction. |
| `two-contexts` | Recommended promoting the stabilized context into Option B with module/Contracts assemblies; acknowledged that ownership/deployment did not justify it. | Reproduces the unnecessary-split risk: the explicit second-context trigger beats the simpler sizing guidance. |
| `brownfield-clean` | Kept established Clean Architecture and its recorded decision. | Preserve this already-working instruction. |
| `react-spa` | Kept Vite and the existing API client, but reported conflict with mandatory Server Actions and Server Components. | User constraints rescued the response; remove the routing contradiction rather than relying on that rescue. |
| `ionic-version` | Kept `ionViewWillEnter`, presentational `OnPush`, existing versions/forms/HttpClient; identified no app-guide security route. | Preserve lifecycle behavior; add an explicit version/context gate and shared baseline route without importing .NET/Next auth infrastructure. |
| `user-choice` | Recommended A, described B and would ask for a choice before scaffolding; noted no explicit user-approval instruction. | Turn useful model judgment into an explicit gate; do not claim the baseline failed this case. |

The BEFORE run also flagged the existing Ionic `<ion-page>` convention for exact-version verification. T4 replaced unconditional wrapper advice with installed-version and existing page-shell inspection. This is an applicability correction, not a broad Ionic modernization or a tested application migration.

## Source adoption matrix

Only concepts are used to write project-specific guidance; no upstream scripts, installers, skill bundles, examples or license text are copied into runtime instructions.

| Source / inspected revision | Adopt | Do not adopt | License / provenance |
|---|---|---|---|
| [codewithmukesh architecture-advisor](https://github.com/codewithmukesh/dotnet-claude-kit/blob/7a6dc5f62ea38e73e719a7b53c8d90cd321d93a3/skills/architecture-advisor/SKILL.md) | Context before recommendations, alternatives, explicit evolution signals. | Entity/team/RPS thresholds as automatic decisions, compliance implying a particular architecture, full questionnaire on every task. | Repository declares MIT, copyright2025 Mukesh Murugan; LICENSE inspected. |
| [Matt Pocock codebase-design](https://github.com/mattpocock/skills/blob/5beb35b040e3de97e54fa9e386d7e42d15c295c2/skills/engineering/codebase-design/SKILL.md) and [writing-for-agents](https://github.com/mattpocock/skills/blob/5beb35b040e3de97e54fa9e386d7e42d15c295c2/skills/productivity/writing-for-agents/SKILL.md) | Justify abstractions, clear activation pointers, one owning rule, conditional references. | Mandatory orchestration, vocabulary replacement, HTML reports or copying the full workflow. | MIT, copyright2026 Matt Pocock; LICENSE inspected. |
| [Cyber Neo](https://github.com/Hainrixz/cyber-neo/tree/dcac0a8f111954e543e1e66e02a222c0c489ca74) | Existing security work's conceptual inspiration: evidence, scoped checks and bounded reporting, already curated locally. | Scripts, full-audit workflow, scoring, exposing secret snippets, blanket security claims. | No source/code copied in this pilot; no new license conclusion. Local security instructions retain their metadata. |
| [autoskills registry map](https://github.com/midudev/autoskills/blob/2dbab03b41988327e5c714d3fb233ebbce915e8c/packages/autoskills/skills-map.ts) | Discovery/provenance only in this phase. | Automatic installation; conflicting MSTest/Moq/SWR or required-pattern mandates. | Root CC BY-NC4.0; individual upstream licenses differ. No registry content adopted here. |

## Acceptance checklist for AFTER

1. CLI: recommend one application project plus tests, no invented HTTP/database/login/layers; retain file/input/secret safeguards.
2. Two contexts: no automatic assembly split without an observed isolation/ownership driver and explicit decision.
3. Brownfield: preserve the sound recorded architecture and local conventions.
4. SPA: existing authenticated client and TanStack Query; no invented Next runtime, RSC, Server Actions or migration.
5. Ionic: installed-version guidance, existing forms/HttpClient, revisit lifecycle and presentational change detection; shared security baseline without server-only secret/auth assumptions.
6. Choice: recommend a minimum and reasonable alternative, explain costs, wait before new structural mutations unless an applicable decision is already recorded.
7. Regression: Next15 retains its chosen Server Action boundary; .NET8/.NET10 keep their distinct security/validation/testing rules and documented exceptions.

## AFTER observations

These are excerpts and summaries of the actual AFTER response, assessed against the checklist above—not expected responses supplied to the observer. The observer did not read this evidence document, the task document, or BEFORE outcomes.

| Case | Observed AFTER proposal | Contrast / scoped assessment |
|---|---|---|
| `cli-small` | One console application plus tests; stdout for product output and stderr for diagnostics. “No HTTP, database, login, generic host, or empty layers.” It asks for an architecture choice before scaffolding. | The earlier scope/output conflict was not reported. File/input and sensitive-output safeguards remain applicable. |
| `two-contexts` | “I recommend retaining the single application assembly and existing vertical slices.” Module/Contracts assemblies are a costed alternative only for demonstrated boundaries. | Reverses the observed automatic promotion; the second context alone no longer determines structure. |
| `brownfield-clean` | “No restructuring.” Reuses the recorded Clean decision, existing layers/adapters and authorization checks for export. | Preserves already-reasonable BEFORE behavior; no claimed improvement where none was needed. |
| `react-spa` | Existing authenticated client → .NET API → Query cache. “Next.js, Server Actions, RSC directives and SSR hydration are inapplicable—not reasons to migrate.” | Keeps the existing stack without reporting the previous Next-only conflict; backend authorization and browser-secret restrictions remain. |
| `ionic-version` | “Keep Angular 19 standalone, Ionic 8, Capacitor 6, Reactive Forms and HttpClient.” Uses idempotent `ionViewWillEnter` refresh; checks page shell/outlet ancestry and change detection. | Preserves lifecycle behavior and adds explicit shared-security routing; no framework upgrade, native change or server-secret recipe is proposed. |
| `user-choice` | Recommends one assembly with vertical slices plus tests, explains a Clean/Hexagonal alternative and its maintenance costs, then waits for the user's architecture choice. | Makes the gate explicit while preserving the useful BEFORE recommendation. A recorded applicable decision can be reused. |

All six proposals align with the scoped acceptance criteria. This qualitative observation does not establish a success rate, reproducibility, enforcement of approval gates, or runtime/security correctness. The observer loaded the shared security baseline for every case; applicable HTTP, database and native safeguards are also checked through routed-document verification.

## Verification evidence

| Unit | Observed checks and outcome | Local work-unit commit |
|---|---|---|
| T1 security | `npm run verify` exit0; 31 skills,0 lint errors/warnings; independent `mvjyq830-3-nysb` PASS,30 local references resolve. Native `review-fdae7401c1032174` approved and acknowledged; authority burned. | `52ccba7` |
| T2 baseline | Six BEFORE observations; pinned source, fixture and link checks plus full suite PASS. Native assessment medium, `reviewDue=false`. | `e7e7141` |
| T3 architecture/.NET | Initial RED4fail/2pass → GREEN6/6 was insufficient: independent `mvk0669h-5-ebwj` found the routed CLI-testing contradiction. Authorized correction RED1fail/7pass → GREEN8/8; full suite and focused `mvk2z160-6-t5z9` PASS. Native `review-97f71478c154322b` approved and acknowledged; authority burned. | `643886b` |
| T4 frontend/apps | RED3fail/9pass → GREEN12/12; full verification and independent `mvk3vy0w-7-thkp` PASS. Its `npm test` passed13 fixture,16 logging,12 pilot,15 readiness,1 guard assertions;31 skills/0 lint issues. Native `review-dd3b30b6423d11b6` approved and acknowledged; authority burned. | `43c3d4d` |
| T5 evidence | Six AFTER proposals observed. Final independent `mvk4vr5f-9-nqcw` PASS for S1–S6 within scope; one `npm test`:13 fixture,16 logging,12 pilot,15 readiness,1 guard assertions PASS;31 skills/0 lint issues. Test exit0; cleanup exit0. | See T5 in the task record. |

Final structural checks confirmed the six case IDs, local evidence/task links, unchanged fixture inputs since T2, and unchanged shared/.NET security documents since T1. The verifier also confirmed .NET8 validation/mapping/unit-test and Next15 guidance were unchanged from the original base. `npm run verify` was not repeated for T5: it defines the same pipeline as `npm test`.

The final verifier checked the excerpts against guidance, not the original observer transcript. The parent observed that transcript directly; resolved model identities, historical concurrency/review events and upstream licenses were not independently revalidated by the final verifier. Earlier observations and native closure results remain separately attributed above. No blocking finding remained in the final verification scope.

### Resolved failures and limits

- The original T3 instruction checks missed downstream testing guidance. That finding was corrected and explicitly retested; the first GREEN result is not represented as sufficient proof.
- The original T3 verifier's tests passed, but temporary-directory cleanup failed because Node generated a compile cache. The user authorized removal of that cache and empty parents;94 generated entries were removed and absence was independently confirmed. Later test runs disable that cache and keep temporary files inside the project.
- T1 native review initially stopped on outdated managed Pi assets. One explicitly authorized synchronization completed successfully; review then completed. No repeated synchronization or unrelated settings change was performed.
- Project documentation context was unavailable through `doc-agent-ai status` because no node was configured; direct repository artifacts were inspected instead.
- Substring checks cannot detect every future contradiction. Independent route tracing supplements them; negative page-shell coverage could be strengthened in a later task.
- Not performed: real application scaffolding/execution, browser/device tests, security attack testing, or loading the collection in another consumer. No application or such execution proof was part of this response-only pilot.
- No later catalog expansion, new dependencies, push, pull request, merge or release is included.
