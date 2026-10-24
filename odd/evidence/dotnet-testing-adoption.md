# .NET testing adoption — phase 2C evidence

Scope: [S1–S4/T1–T2](../tasks/dotnet-testing-adoption.md). Baseline `b1e94f7`; guidance unit `78cfd95`. No new skill, profile unification, security rewrite, global tooling installation, CI requirement or phase3 expansion.

## Delivered surfaces

- Existing backend/testing-unit v2.1 and .NET10/testing v2.3: applicable exact/adjacent boundaries, rejection plus prohibited observable effects, and deterministic cutoff/supported cancellation guidance.
- .NET8 keeps xUnit/NSubstitute/FluentAssertions and its naming; .NET10 keeps native Assert/manual spies and its naming. Existing routes, CLI/WAF/real-engine/security/transaction/CI rules remain.
- Optional [real xUnit trial](../../tools/fixtures/dotnet-testing/README.md): one pure capacity rule, two theories/six cases, fixed isolated prerequisites. Default checks stay Node-only; execution/help never installs.

## Observed evidence

| Check | Result |
|---|---|
| Read-only map | `mvlaylwp-f-t2df` identified scoped omissions and absent evidenced project-local xUnit prerequisites;external/ignored caches not inspected |
| Permission | User explicitly chose "Autorizar requisitos NuGet para el ensayo aislado";only xunit2.9.x,runner3.x,TestSDK17.x+transitives,officialNuGet,project-local home/cache/temp |
| T1 instruction RED → GREEN | 2pass/4fail,exit1 →6/6;preservation checks already passed before implementation |
| T1 suite/assessment | Full suite33skills/0lint;verify/cleanup0. Diff initially flagged3new Markdown trailing-space breaks;fixed/rechecked. ASSESSmedium,reviewDue=false |
| Metadata/restore | Jobs6/7exit0;fixed xunit2.9.3,runner3.1.5,TestSDK17.14.1;local cache/lock generated,no SDK install |
| Real xUnit RED (`job-8`) | 6discovered,4pass/2fail,exit1;accepted-add cases failed with stub;runtime/compiler worked;cleanup removed |
| Real xUnit GREEN (`job-9`) | 6/6,exit0,0compiler warnings/errors;cleanup removed;both returned outcome and public total asserted |
| Runner/guide checks | 9/9;help/unknown/missing prerequisites exits0/2/2,exact stdout/stderr,unchanged synthetic state,no SDK or downloads |
| T2 suite | 13fixture/16logging/12pilot/9CLI/9frontend/9.NET/15readiness/1guard;33skills/0lint;verify/cleanup0 |
| Independent verifier | Fresh `mvlbkslw-g-n0ys`:S1–S3 PASS,no code blockers;both public commands once,exit0/stderr0,84suite checks/33skills0lint,6xUnit. S4 initially cleanup-only NEEDSREVISION after child-safety block;resolved below without retests |
| Authorized cleanup | Human explicitly allowed only `.dotnet-unit-trial-verifier-TcHnwA/`;parent removed it,cleanup0/absence confirmed;Git status/index/worktree byte-identical to verifier snapshots,package cache untouched |
| Native assessment/review | T1medium,reviewDue=false;T2high (`process_boundary`). Four-lens `review-df6f69699af0c2af` approved;exact acknowledgement completed,authority burned. Frozen scope18paths/406changed lines (98lock),max2concurrent reviewers;no old approval/refusal reused |

Runner-interface regression tests were added after its implementation,not claimed as a separate RED cycle. The capacity cases did observe RED before implementing the rule. The guard compares the positive amount to remaining capacity before addition,so `int.MaxValue` cannot wrap and bypass rejection.

## Prerequisites, ownership and limits

- Separately authorized restore used official nuget.org with project-local packages/home/http/plugin/temp/MSBuild extension paths. Initial restore generated the lock; ordinary execution uses locked mode,empty-source config,local packages and `test --no-restore` in an owned disposable copy. Root dependencies/CI/global settings do not change.
- Lockfile records resolved direct/transitive versions and content hashes;it is not a security/license audit or proof that downloaded MSBuild/test-runtime code is sandboxed. NuGetAudit is disabled,not silently reported passed. SDK/runtime/host environment and approved packages remain trusted.
- Runner/TRX checking requires exactly six passed cases with no failures;empty discovery cannot count as success. Temporary copy cleanup is reported separately;forced interruption/timeouts can prevent complete cleanup or process-tree termination. The approved ignored package cache remains a reusable prerequisite,not leftover test state.
- The pure rule has no persistence,clock,cancellation,network,HTTP,DB,publication or authorization seam. It checks result and in-memory total,not fabricated collaborator behavior. Time/cancellation improvements are instruction checks only;existing .NET8 library style is preserved,not executed by this .NET10 specimen.
- No WAF/real-engine/transaction/security behavior,all-platform compatibility,coverage percentage,model compliance or application-security guarantee is established. Earlier CLI/React acceptance proofs remain separate.

## Non-blocking native advisory

Native approval lists `R3-stateful-capacity` (reliability),`tools/fixtures/dotnet-testing/CapacityCounterTests.cs:12`,severity WARNING,disposition informational. No explanatory text was exposed in the closure envelope;no cause is inferred here. Native explicitly retains approval and opens no correction;this is separate later work,not grounds to reopen/re-run this candidate.

## Local delivery

T1:`78cfd95ba8a579417a7d58647a0274a7ee52555b`. T2:`1fa8180d443e25f4de37159532abedf5f70186aa`,preserving approved tree `ab6159459c1361907815d8fce87bad0fbbc5b447`. Required functional/independent/native gates and authorized cleanup are complete. This final documentary update records outcomes only;no publishing or later phase started.
