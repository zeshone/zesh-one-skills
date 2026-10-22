# Architecture adoption pilot evidence

Scope: phases 0–1; existing security enrichment is the baseline, not a separate optional add-on. No later catalog additions are included.

## Reproduction inputs and limits

- BEFORE source: `52ccba7b5b05a085ac0164fadcd959a1cbb9a055`.
- Fixed prompts and loading instructions: [six response-only cases](../../tools/fixtures/architecture-adoption-pilot.json).
- BEFORE observation: `gentle-ai-explore`, task `mvjzjsoc-4-s8wv`, foreground, one additional Pi session.
- Use the same agent role/configuration and unchanged prompts for AFTER, in a fresh task. The host did not expose a resolved model/effort identity; do not claim a controlled model comparison.
- These are observed model proposals, not scaffold execution tests. Read-only tool restrictions mean this exercise cannot demonstrate that skill prose alone prevents writes, sandbox escape, secret exposure, or dependency installation.
- No application was created or modified. Runtime loading and behavior in another consumer remain unverified. The independent final check must examine both documents and observations, rather than equating text presence with compliance.

## BEFORE observations

| Case | Observed response | Interpretation for this pilot |
|---|---|---|
| `cli-small` | Proposed one console project plus tests, but reported REST-only scope and the `Console.WriteLine` prohibition as conflicts. | Simple recommendation already possible through model judgment; instructions still need an explicit CLI profile and product-output distinction. |
| `two-contexts` | Recommended promoting the stabilized context into Option B with module/Contracts assemblies; acknowledged that ownership/deployment did not justify it. | Reproduces the unnecessary-split risk: the explicit second-context trigger beats the simpler sizing guidance. |
| `brownfield-clean` | Kept established Clean Architecture and its recorded decision. | Preserve this already-working instruction. |
| `react-spa` | Kept Vite and the existing API client, but reported conflict with mandatory Server Actions and Server Components. | User constraints rescued the response; remove the routing contradiction rather than relying on that rescue. |
| `ionic-version` | Kept `ionViewWillEnter`, presentational `OnPush`, existing versions/forms/HttpClient; identified no app-guide security route. | Preserve lifecycle behavior; add an explicit version/context gate and shared baseline route without importing .NET/Next auth infrastructure. |
| `user-choice` | Recommended A, described B and would ask for a choice before scaffolding; noted no explicit user-approval instruction. | Turn useful model judgment into an explicit gate; do not claim the baseline failed this case. |

The BEFORE run also flagged the existing Ionic `<ion-page>` convention for exact-version verification. This pre-existing rule is not redesigned by the applicability pilot; no broad Ionic modernization is authorized.

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

Pending T3–T5. Do not interpret this placeholder or the acceptance checklist as execution evidence.

## Verification evidence

- Security baseline: `npm run verify` exit0; 31 skills,0 lint errors/warnings; independent verifier `mvjyq830-3-nysb` PASS;30 local references resolve.
- Security native review: `review-fdae7401c1032174`, approved and exact acknowledgement completed (authority burned).
- Pilot instruction checks, AFTER observations and final independent verification: pending.
