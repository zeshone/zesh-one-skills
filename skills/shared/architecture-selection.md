# Contextual architecture selection

Use this decision procedure from the existing domain architecture skill when creating a project or changing structural boundaries. It is not another runtime stack or a requirement to redesign each feature.

## Inspect, recommend, then wait

1. Inspect the application type (CLI, worker, HTTP API, SPA, SSR or mobile), installed framework versions, existing architecture and recorded decisions. Identify actual domain complexity, team ownership, trust/data boundaries and deployment constraints. Ask only for missing context that changes the decision; do not run a fixed questionnaire for routine changes.
2. Recommend the minimum sufficient structure and at least one reasonable alternative. Explain the concrete benefit, added projects/abstractions, maintenance/testing cost and evidence that would justify changing later. Say when an alternative adds no present value.
3. Wait for the user to choose before scaffolding, splitting/merging projects or restructuring. A recommendation is not authorization. If an applicable choice is already recorded and the user has authorized work within it, reuse it rather than asking on every feature. Record a newly approved choice, its rationale and revisit signals in the project's existing decision location; no new ADR framework is required.
4. Preserve a sound brownfield architecture. A defect can justify proposing a bounded migration, not silently rewriting the system. Move incrementally only after the structural change is approved.

## Profiles, not size thresholds

| Profile | Useful when | Cost / boundary |
|---|---|---|
| Minimal monolith with feature/vertical-slice organization | Small CLI, worker or API; cohesive ownership and one deployment | Start with one application project plus tests. Folders and code review express boundaries; avoid empty layers or speculative interfaces. |
| Clean/Hexagonal responsibility separation | Business rules must stay independent of changing adapters, or this is the sound existing structure | Ports and separation have a maintenance cost. Logical layers do not require one assembly each; physical enforcement needs a separate justification. |
| Modular monolith with module assemblies | An observed boundary problem, distinct ownership or independently evolving contracts needs compiler-enforced isolation | More projects/contracts and explicit data ownership; still one deployment. Agree on the boundary before extracting it. |

Vertical slices organize changes; Clean/Hexagonal separates responsibilities; assemblies enforce physical boundaries. They can be combined, not treated as mandatory growth stages. DDD tactical patterns address actual domain complexity across any profile; simple CRUD does not need aggregates/events by ritual. None of these choices authorizes microservices or a stack upgrade.

A second bounded context alone, entity count, LOC, RPS, team size or hypothetical future growth are not automatic triggers for more assemblies. Use these as questions, not thresholds. Recommend keeping the current layout when there is no observed isolation/ownership driver. Revisit after concrete coupling, ownership or independent-change evidence, and obtain a new explicit decision before changing structure.

## Applicability and security

Apply transport, persistence, authentication and deployment rules only to capabilities that exist or are explicitly requested. A local CLI does not need HTTP endpoints, OpenAPI, middleware, a database or login simply to satisfy an API guide. A SPA does not become Next.js; a mobile application does not become a server.

Always load the [shared security baseline](security/baseline.md). Input/file handling, secret protection and existing controls remain applicable without HTTP. Specialize by the actual changed trust boundary; do not disable an existing protection or create infrastructure merely to activate a specialist rule. Keep each domain/version's security, validation and testing conventions separate.
