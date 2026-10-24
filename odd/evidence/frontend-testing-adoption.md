# Frontend testing adoption — phase 2B evidence

Scope: [S1–S4/T1–T2](../tasks/frontend-testing-adoption.md). Baseline `a19a3d2`; guidance unit `4da572b`. No .NET testing enrichment, Next/Ionic rewrite, automatic registry refresh or publishing.

## Delivered surfaces

- Project-authored [React SPA testing guidance](../../skills/frontend/testing/SKILL.md), conditional frontend route and instruction/compatibility checks.
- Optional [rendered-component trial](../../tools/fixtures/frontend-dom/README.md): API-backed list/form, controlled injected client, semantic queries and deferred async responses. This is not production Query integration or actual-browser execution.
- Isolated package manifest/lockfile: React/react-dom19.3.1,TestingLibraryReact16.3.3,DOM10.4.2,jsdom26.1.0. Root dependencies/CI remain unchanged; default verification does not install/load the trial packages.

## Observed evidence

| Check | Result |
|---|---|
| Read-only map | `mvl8y53c-d-0iyu`: no React/DOM tools in initial collection; prerequisite absence disclosed before installation |
| Explicit permission | User chose "Autorizar dependencias solo para el ensayo aislado"; npm packages/cache/project remain local, scripts disabled, no global installs |
| T1 instruction RED → GREEN | 0pass/6fail,exit1 →6/6; first full suite caught old CLI pipeline expectation, corrected without removing old checks |
| T1 corrected suite | 13fixture/16logging/12pilot/9CLI/6frontend/15readiness/1guard;33skills/0lint;verify/cleanup exit0 |
| Metadata/install | Jobs4/5exit0; metadata scratch removed;58packages installed in isolated fixture,scripts off,optional native peers omitted |
| DOM RED | Empty component:0pass/6fail,exit1;runtime loaded,cleanup removed. Five visible-behavior failures; list-error case additionally rejected an unconsumed promise, fixed by asserting visible loading before rejection |
| DOM GREEN | Same six rule cases pass through rendered UI:6/6,exit0,cleanup removed; no expected output/effect assertion weakened |
| Runner/guide checks | 9/9: six instruction checks plus exact help/unknown/missing-prerequisite process status/stdout/stderr and unchanged synthetic state, with no package install |
| T2 full suite | 13fixture/16logging/12pilot/9CLI/9frontend/15readiness/1guard;33skills/0lint;verify/cleanup exit0 |
| Dependency structure/preservation | 58locked packages/766generated lockfile lines,officialregistry/version/integrity present,no install-script declarations; baseline rootlock/CI/backend/apps/security/pilot inputs unchanged |
| Independent verifier | Fresh `mvl9uh04-e-uemv`: PASS S1–S4; `npm test` and `npm run test:frontend-dom` once each,exit0/stderr0;same suite counts and6/6DOM;diff checks clean,verifier cleanup0,no unexpected changes/leftovers |
| Native assessment/review | T1medium,reviewDue=false; T2high (`process_boundary`). Four-lens `review-d6db3ef49b68c57e` approved; exact acknowledgement completed,authority burned. Frozen scope16paths/1202changed lines (766generated lockfile);max2concurrent reviewers. No old approval/refusal reused |

The three runner-interface tests were added after the runner implementation; they are observed regression checks, not a separately claimed RED/GREEN cycle. The six component rules did observe RED before implementation, with the test-harness precondition correction attributed above.

## Boundaries and limitations

- Installation was explicitly authorized, not inferred from `continua` or test execution. npm configuration/cache/temp was redirected into this repository; lifecycle scripts disabled and optional native peers omitted. Downloaded packages and caches are ignored local prerequisites, not committed assets. The manifest and lockfile are committed reproducibility metadata, not a claim that dependencies were audited.
- Node23.11.0 was already installed. Latest jsdom30.1.3 requires a different supported Node line; selected26.1.0 (`node>=18`) without upgrading Node. The install warned that transitive whatwg-encoding3.1.1 is deprecated; installation still exited0. No automatic remediation or security/license audit was performed.
- The trial controls an in-memory client and checks public API-visible records/write calls plus unchanged localStorage via the Storage interface. It does not prove real HTTP, backend auth/ownership/validation, production cache wiring or database behavior.
- jsdom is simulated DOM: no layout, navigation/focus fidelity, native browser events/form validation, accessibility audit or real-browser/E2E proof. Semantic queries alone are not accessibility compliance. No model-compliance or application-security guarantee.
- The runner creates/removes only its owned project-local working/temp directory; tests unmount React roots, check DOM/storage cleanup, close jsdom and restore globals. Host environment and downloaded modules remain trusted; this is not sandboxing or proof of no ambient effects. Forced interruption/timeouts may leave descendants or cleanup incomplete.
- No consumer framework/form/client migration is required by the skill or specimen. Existing SPA/Next/Ionic boundaries and shared architecture/security owners remain intact. Direct guide links are explicit loading routes, not proof of external registry discovery.

## Non-blocking native advisory

Native approval lists `R3-001` (reliability), `tools/fixtures/frontend-dom/NotesPanel.mjs:37`, severity WARNING, disposition informational. No explanatory text was exposed in the closure envelope; no cause is inferred here. Native explicitly retains approval, opens no correction and treats this as separate later work, not grounds to reopen/re-run this candidate.

## Local delivery

T1: `4da572bb950b1c45c2873784a58bb8d528d6a0be`. T2: `cc98088e599f25e1e276f32bad3dfedfbe6722c3`,preserving approved tree `a78cf2a31e228aa7da0981a662cf09a5b748ce56`. All required functional/independent/native checks complete. This final documentary update records outcomes only. No publishing or later phase started.
