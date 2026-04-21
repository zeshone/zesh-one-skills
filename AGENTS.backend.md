# Backend Domain Agent Guide (.NET 8 REST)

## Scope
Use this file as the execution default for backend work in this repo. Keep outputs pragmatic and implementation-ready. If a domain skill conflicts with this file, the skill wins.

## Mandatory Skill Routing (load first)

| Work Type | Skill |
|---|---|
| Architecture, project structure, feature scaffolding | `net8-apirest-general` |
| EF Core, DbContext, repositories, migrations, query optimization | `net8-apirest-dataaccess` |
| JWT/OAuth, authorization, OWASP controls, rate limiting, BOLA mitigation | `net8-apirest-security` |
| Request DTOs, model binding, query params, file uploads | `net8-apirest-requests` |
| Validation rules and validators | `net8-apirest-validations` |
| API response contracts and status mapping | `net8-apirest-responses` |
| Mapping strategy (AutoMapper/manual) and profiles | `net8-apirest-mapping` |
| Serilog, correlation, exception logging and traceability | `net8-apirest-logging` |
| Resilience and performance (Polly, limits, parallelism) | `net8-apirest-performance` |
| Unit tests (xUnit, NSubstitute, FluentAssertions) | `net8-apirest-testing-unit` |
| PR preparation/review | `github-pr` |

## Operational Defaults
- Build features with Clean Architecture + Vertical Slice conventions from `net8-apirest-general`.
- Prefer `Result<T>` contracts; treat `ResponseDTO<T>` as legacy compatibility only.
- Validate external input at the edge with FluentValidation.
- Keep API boundaries explicit: request DTO in, response contract out.
- Keep security and observability as first-class requirements, not post-work additions.

## Non-Negotiables
- Do not skip skill loading for domain-specific changes.
- Do not ship endpoints without validation, auth checks (when required), and structured error mapping.
- Do not bypass logging/correlation middleware for production flows.
- Do not optimize performance blindly; apply resilience/perf patterns from the dedicated skill.
