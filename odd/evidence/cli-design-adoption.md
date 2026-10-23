# CLI design adoption — phase 2A evidence

Scope: [S1–S4 and T1–T2](../tasks/cli-design-adoption.md). Baseline `81ff3cf`; T1 `d6ac394`. No phase2B/2C, publishing or global configuration changes.

## Delivered surfaces

- Project-authored [CLI interaction skill](../../skills/shared/cli-design/SKILL.md), conditional .NET8/.NET10 guide routes, and Node-only instruction/compatibility checks.
- Opt-in [line-tally acceptance specimen](../../skills/shared/cli-design/assets/line-tally/README.md), [fixed public cases](../../skills/shared/cli-design/assets/line-tally/contract.json) and [Node runner](../../tools/test-cli-contract.js). No third-party package or parser dependency.
- `npm run test:cli` requires an installed .NET10 SDK/reference packs; default `npm test`/`verify` stay Node-only. Standalone runner help/rejected options require no SDK and create no trial state.

## Observed checks

| Check | Observation |
| --- | --- |
| T1 instruction RED → GREEN | 6fail/1pass → 7/7; full suite32skills,0lint errors/warnings |
| Independent preflight | `mvkk7qcb-b-fmk9`: GO for controlled RED only; not final compliance/OS-isolation proof |
| Runner-interface RED → GREEN | 7pass/2fail,exit1 → 9/9,exit0; help/error checked without SDK PATH |
| Executable RED (`job-2`) | Build succeeded;13/13 public cases failed on the deliberately empty implementation;exit1;cleanup removed |
| Executable GREEN (`job-3`) | Same13cases passed;exit0;cleanup removed |
| Parent full verification | `npm run verify`:13fixture/16logging/12pilot/9CLI/15readiness/1guard;32skills;0lint errors/warnings;verification/cleanup exit0 |
| Post-refactor focused check | 9/9,exit0; test-runtime import deferred until executable mode |
| Final independent verifier | `mvklcr7p-c-wwf8`: PASS S1–S4 within scope; `npm test` and `npm run test:cli` each once,exit0; same suite counts and13/13 executable cases; both diff checks passed;trial/verifier cleanup removed,exit0 |
| Native assessment/review | T1medium,reviewDue=false. T2high (`process_boundary`),reviewDue=true; review pending; previous aggregate dispositions do not cover it |

Every executable case checks status, exact stdout and stderr, lack of a signal/timeout, and a shallow name/type/content-hash snapshot of its fixture directory. Synthetic input covers file/stdin, plain/JSON output, LF/CRLF, empty/blank/final lines, option-like/spaced filenames and rejected arguments/input-read failure. Stdin remains open for help, file and error cases to expose unintended prompting/consumption. There are no stored application counters.

The build and tests use a unique project-local directory with redirected named caches/temp, empty package/audit source lists, disabled inherited directory build imports and no shared build/compiler server. Cleanup targets only that invocation's directory. Node launches set `NODE_DISABLE_COMPILE_CACHE=1`. The independent verifier also redirected npm cache, read separate stdout/stderr captures, and confirmed no candidate mutations or generated leftovers. Its `npm test` stderr contained only an npm update notice (no update performed); executable-trial stderr was empty.

## Registry limitation accepted by the user

The single authorized native refresh exited0 and generated33 entries, but omitted `cli-design`; linter discovery found32 collection skills. These are different inventories, not equivalent success signals. The cause was not established. The index was retained unchanged and mirrored in Engram; the one-time grant is consumed.

User selected **"Continuar con rutas explícitas y registrar el límite"**. Use the explicit guide links; do not claim automatic registry discovery. No refresh retry, manual index repair or additional external source access occurred.

## Limits and preserved boundaries

- Acceptance of this specimen is not xUnit/unit coverage, agent-generated output evaluation, model compliance, collection loading in another consumer or an application-security guarantee.
- Test-state isolation is not an OS sandbox or hermetic execution. The inherited environment, PATH-selected installed SDK/runtime and host remain trusted; no claim is made that arbitrary ambient effects are impossible.
- Timeouts do not prove process-tree termination; forced interruption can prevent cleanup. Snapshots cover case directories, not the host.
- Other operating systems/Node versions, broken pipes, interrupt mapping, huge input and arbitrary hostile manifests were not established by this local trial. The manifest is repository-owned test data, not an external input API.
- Existing architecture/security/testing owners, .NET8/.NET10 distinctions, pilot scenarios and default Node-only CI remain applicable. The specimen's inline compiler settings are an explicit fixture-only exception to shared Directory.Build.props because the runner copies the isolated project without ancestor imports.
- External skill metadata indexing permission and native review consent are separate; no user delivery authorization extends beyond local work-unit commits.
