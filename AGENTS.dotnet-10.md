# Backend Domain Agent Guide (.NET 10 / C# 14 Minimal API)

## Scope
Execution default for backend work on **any** .NET 10 / C# 14 REST API. Keep outputs pragmatic and implementation-ready. The domain skills under `skills/dotnet-10/` are the normative source; **if a skill conflicts with this file, the skill wins.** This guide routes work to the right skill and states the cross-cutting defaults that hold regardless of the project's architecture.

## Mandatory Skill Routing (load first)

| Work Type | Skill |
|---|---|
| C# 14 language style, async/cancellation, nullable, naming, SOLID, XML docs, timestamps | `dotnet-conventions` |
| Project structure, layering, architecture selection (Option A vs B), where logic lives, boundary/seam placement | `architecture` |
| Minimal API endpoints, TypedResults, ProblemDetails, routing/versioning, OpenAPI/Scalar, middleware order | `api-rest-minimal-apis` |
| EF Core, DbContext, migrations, raw SQL, `ISqlExecutor` ports | `data-access-persistence` |
| Auth, JWT signing, password hashing, TOTP, CSRF, CORS, rate limiting, security headers, secrets | `security` |
| Tests, RealPg/Testcontainers, `FakeTimeProvider`, xUnit, test naming, coverage strategy | `testing` |
| Performance, AOT/trimming posture, async I/O, query projections, allocations, pooling | `performance` |
| Serilog logging, OpenTelemetry posture, health checks, `BackgroundService` lifecycle, deploy/migration ops | `observability-operations` |
| Result discriminated unions, Options pattern, DI lifetimes, conditional service selection, exception mapping | `design-patterns` |

## Architecture is a per-project choice (not a fixed mandate)
- Select the structure with the `architecture` skill: **Option A** — Vertical Slice Monolith with Clean Architecture separation (small-to-medium projects, single application assembly, `Features/{Name}/` organization, dependency rule by discipline); **Option B** — Modular Monolith with vertical-slice module assemblies + `Contracts` assemblies (large projects; microservice-style boundaries without distributed-systems cost; extraction path to services). Start with A; promote a slice to a B module when a second bounded context stabilizes.
- **Brownfield rule**: a solid, consistently applied Clean or Hexagonal architecture is NOT restructured — work within it. Only a fragile, mixed architecture is migrated incrementally (strangler-fig) toward the size-appropriate option.
- **Domain model**: rich domain entities are the default (standard Clean Architecture). Match DDD tactical investment (aggregates, value objects, domain events) to actual domain complexity — no rituals for simple CRUD.

## Operational Defaults (cross-cutting, architecture-agnostic)
- Model expected errors with `Result` discriminated unions (sealed records); reserve exceptions for unexpected infrastructure failures. On enumeration-sensitive paths, `Failure` carries no distinguishing reason.
- Return `TypedResults` with explicit `Results<...>` unions, and emit `ProblemDetails` (RFC 7807) on every endpoint and middleware short-circuit.
- Map manually (positional record constructors / EF `.Select()` projections). Validate by hand through the `Result` pattern at the service boundary.
- Organize endpoints in per-feature/per-module groups via static `IEndpointRouteBuilder` extension methods (`app.MapXxxEndpoints()`) — this works in any assembly layout.
- Read time via `TimeProvider`, log via Serilog (`ILogger<T>` + named templates), and load secrets/connection strings from environment variables only.
- When building auth-bearing APIs, follow the `security` skill's defaults: asymmetric-signed JWTs (e.g., ES256), a modern password KDF (e.g., Argon2id), TOTP for 2FA, CSRF double-submit for cookie flows, anchored CORS origins, a sliding-window rate limit, and standard security headers. On enumeration-sensitive endpoints (auth, TOTP, password reset), keep FluentValidation / `AddValidation()` OFF and validate through the `Result` pattern instead.

## Non-Negotiables (development-style defaults — hold for every project)
- Do not introduce MediatR / `IMediator` / `IPipelineBehavior` (commercial license above a revenue threshold, AOT/trimming-hostile reflection registration). If the mediator shape is explicitly wanted, the only acceptable form is a source-generated, MIT-licensed mediator (e.g., `martinothamar/Mediator`) after a documented design decision.
- Do not introduce generic `IRepository<T>`, `IUnitOfWork`, or DbContext wrappers — `DbContext` already is a repository + unit of work, so the wrapper adds drift with no driver.
- Do not use AutoMapper, Mapster, or any reflection/convention mapper (sensitive-field exposure risk) — map manually.
- Do not add Swashbuckle or NSwag; use the built-in `Microsoft.AspNetCore.OpenApi` + Scalar only, and keep the OpenAPI and Scalar endpoints disabled in production.
- Migrations are a controlled CI/CD step — never call `MigrateAsync()` at startup.
- For behavioral and data-access tests, use real PostgreSQL (Testcontainers), not EF Core InMemory — InMemory hides provider-specific behavior. (A narrow WAF-tier DI-wiring exception exists — see the testing skill's InMemory-exceptions rule.)
- Telemetry: Serilog (CLEF) is the logging backbone. Adopt OpenTelemetry (traces/metrics) deliberately and AOT/trimming-safely when distributed tracing or metrics are genuinely needed, and document the telemetry posture (what is collected, what is deliberately excluded) either way.
