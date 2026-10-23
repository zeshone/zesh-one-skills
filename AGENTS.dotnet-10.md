# .NET 10 / C# 14 Agent Guide (API, CLI, worker)

## Scope
Identify the application type, installed framework versions and existing architecture before routing. Use the `architecture` skill for structural choices in .NET 10 APIs, CLIs and workers. The other domain skills remain capability-specific: apply HTTP, persistence and authentication guidance only when that capability exists or is requested, not to create it. Preserve applicable existing controls and load the security baseline even for a local CLI.

The domain skills under `skills/dotnet-10/` are normative; **if a skill conflicts with this file, the skill wins.** HTTP operational defaults below belong to the Minimal API profile, not to every .NET program.

## Always-on security baseline

At the start of every development task, load [the shared security baseline](skills/shared/security/baseline.md) and keep its constraints active. It routes specialized checks by the changed boundary, not by a mandatory audit. Load `skills/dotnet-10/security/SKILL.md` when a listed security boundary changes; retain the .NET 10 conventions.

When adopting this guide, include `skills/shared/security/` alongside `skills/dotnet-10/` and adjust these paths if the suite lives elsewhere. A link or tool declaration is not automatic loading or technical isolation; the host must actually load the guide and referenced baseline.

## Mandatory Skill Routing (load first)

| Work Type | Skill |
|---|---|
| C# 14 language style, async/cancellation, nullable, naming, SOLID, XML docs, timestamps | `dotnet-conventions` |
| Bug reproduction, diagnosis, regression investigation | [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md) |
| Project structure, contextual architecture choice, CLI/worker applicability, layering, boundary/seam placement | `architecture` |
| CLI commands, flags, help, exit codes, stdin/stdout/stderr, non-interactive use | [`cli-design`](skills/shared/cli-design/SKILL.md) |
| Minimal API endpoints, TypedResults, ProblemDetails, routing/versioning, OpenAPI/Scalar, middleware order | `api-rest-minimal-apis` |
| EF Core, DbContext, migrations, raw SQL, `ISqlExecutor` ports | `data-access-persistence` |
| Auth, JWT signing, password hashing, TOTP, CSRF, CORS, rate limiting, security headers, secrets, dependency/infra/CI security | [`security`](skills/dotnet-10/security/SKILL.md) (.NET 10) |
| Tests, RealPg/Testcontainers, `FakeTimeProvider`, xUnit, test naming, coverage strategy | `testing` |
| Performance, AOT/trimming posture, async I/O, query projections, allocations, pooling | `performance` |
| Serilog logging, OpenTelemetry posture, health checks, `BackgroundService` lifecycle, deploy/migration ops | `observability-operations` |
| Result discriminated unions, Options pattern, DI lifetimes, conditional service selection, exception mapping | `design-patterns` |

## Architecture is a per-project choice (not a fixed mandate)
- Load [`architecture`](skills/dotnet-10/architecture/SKILL.md) and its [shared selection procedure](skills/shared/architecture-selection.md). Recommend the smallest sufficient structure with alternatives and costs; wait for the user's choice before scaffolding or changing boundaries. Reuse an applicable recorded decision.
- A second context or speculative growth does not force module assemblies. The architecture skill owns the layouts and their tradeoffs; do not duplicate a competing sizing policy here.
- **Brownfield rule**: preserve a solid, consistently applied Clean or Hexagonal architecture. Propose a bounded incremental migration only for demonstrated defects and after explicit approval.
- **Domain model**: rich domain entities are the default (standard Clean Architecture). Match DDD tactical investment (aggregates, value objects, domain events) to actual domain complexity — no rituals for simple CRUD.

## Operational Defaults (within the applicable capability)
For a CLI/worker without HTTP or persistence, skip endpoint/middleware/database recipes and their dependencies. Use the architecture skill's CLI output/logging distinction; this does not remove validation or sensitive-data safeguards.
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
- For behavioral and data-access tests that exercise PostgreSQL, use real PostgreSQL (Testcontainers), not EF Core InMemory — InMemory hides provider-specific behavior. Do not add a database to test an application that does not use one. (A narrow WAF-tier DI-wiring exception exists — see the testing skill's InMemory-exceptions rule.)
- Telemetry for HTTP APIs and hosted services: Serilog (CLEF) is the logging backbone. Adopt OpenTelemetry (traces/metrics) deliberately and AOT/trimming-safely when distributed tracing or metrics are genuinely needed, and document the telemetry posture (what is collected, what is deliberately excluded) either way.
