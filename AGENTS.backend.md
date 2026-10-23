# Backend Domain Agent Guide (.NET 8 REST)

## Scope
Use this file for .NET 8 REST APIs after checking application type, installed versions and existing decisions. For a CLI or worker without HTTP, use the shared architecture-selection procedure and security baseline, not HTTP-specific controller/middleware/OpenAPI recipes. Do not add HTTP or auth infrastructure just to activate this profile. Keep outputs pragmatic; if a domain skill conflicts with this file, the skill wins.

## Always-on security baseline

At the start of every development task, load [the shared security baseline](skills/shared/security/baseline.md) and keep its constraints active. It routes specialized checks by the changed boundary, not by a mandatory audit. Load `skills/backend/security/SKILL.md` when a listed security boundary changes; retain the .NET 8 conventions.

When adopting this guide, include `skills/shared/security/` alongside `skills/backend/` and adjust these paths if the suite lives elsewhere. A link or tool declaration is not automatic loading or technical isolation; the host must actually load the guide and referenced baseline.

## Mandatory Skill Routing (load first)

| Work Type | Skill |
|---|---|
| Architecture, project structure, feature scaffolding | `net8-apirest-general` |
| Bug reproduction, diagnosis, regression investigation | [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md) |
| CLI commands, flags, help, exit codes, stdin/stdout/stderr, non-interactive use | [`cli-design`](skills/shared/cli-design/SKILL.md) |
| EF Core, DbContext, repositories, migrations, query optimization | `net8-apirest-dataaccess` |
| JWT/OAuth, authorization, OWASP controls, rate limiting, BOLA, secrets, dependency/infra/CI security | [`security`](skills/backend/security/SKILL.md) (.NET 8) |
| Request DTOs, model binding, query params, file uploads | `net8-apirest-requests` |
| Validation rules and validators | `net8-apirest-validations` |
| API response contracts and status mapping | `net8-apirest-responses` |
| Mapping strategy (AutoMapper/manual) and profiles | `net8-apirest-mapping` |
| Serilog, correlation, exception logging and traceability | `net8-apirest-logging` |
| Resilience and performance (Polly, limits, parallelism) | `net8-apirest-performance` |
| Unit tests (xUnit, NSubstitute, FluentAssertions) | `net8-apirest-testing-unit` |
| PR preparation/review | `github-pr` |

## Operational Defaults
- Use `net8-apirest-general` and its [shared architecture-selection procedure](skills/shared/architecture-selection.md) for new structural choices; recommend the minimum with alternatives and wait for the user's choice before scaffolding. Preserve a sound existing architecture. Within an approved feature-slice layout, keep the .NET 8 Clean Architecture + Vertical Slice conventions.
- Prefer `Result<T>` contracts; treat `ResponseDTO<T>` as legacy compatibility only.
- Validate external input at the edge with FluentValidation.
- Keep API boundaries explicit: request DTO in, response contract out.
- Keep security and observability as first-class requirements, not post-work additions.

## Non-Negotiables
- Do not skip skill loading for domain-specific changes.
- Do not ship endpoints without validation, auth checks (when required), and structured error mapping.
- Do not bypass logging/correlation middleware for production flows.
- Do not optimize performance blindly; apply resilience/perf patterns from the dedicated skill.
