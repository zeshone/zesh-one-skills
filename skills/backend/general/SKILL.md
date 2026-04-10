---
name: net8-apirest-general
description: >
  General architecture and coding conventions for ASP.NET Core 8 REST APIs using Clean Architecture with Vertical Slices.
  Trigger: When creating, scaffolding, or structuring any .NET 8 REST API project, controllers, services, repositories, or feature folders.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.7"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Creating a new .NET 8 REST API project or feature
- Scaffolding controllers, services, repositories, or DTOs
- Reviewing project structure and folder organization
- Applying naming conventions or deciding code style
- Registering dependencies in `Program.cs`

---

## Critical Patterns

### Project Structure — Vertical Slices + Clean Architecture

```
Features/
  {Feature}/
    Controllers/   ← HTTP endpoints, thin layer
    Services/      ← Business logic (interface + implementation)
    Repositories/  ← Data access (interface + implementation)
    DTOs/          ← Request/Response transfer objects
    Models/        ← Domain entities (EF Core mapped)
    Mappings/      ← Extension methods or AutoMapper profiles
    Validators/    ← FluentValidation rules
    Exceptions/    ← Domain-specific custom exceptions

Shared/
  Models/          ← ResponseDTO<T>, PagedResult<T>, shared models
  Models/BaseEntity.cs  ← Id, CreatedAt, UpdatedAt (see dataaccess/SKILL.md)
  Exceptions/      ← Shared custom exceptions
    NotFoundException.cs
    ConflictException.cs
  Helpers/         ← Global utilities (e.g., PasswordHasher)
  Middlewares/     ← Global pipeline middlewares
  Extensions/      ← Extension methods
  Interceptors/    ← EF Core interceptors (e.g., AuditInterceptor)

Database/
  Context/         ← EF Core DbContext definition

tests/             ← Unit & integration tests (see testing-unit/SKILL.md)
```

### Naming Conventions

| Element | Convention | Example |
|---|---|---|
| Private fields | `_camelCase` | `_logger`, `_repository` |
| Source code | English | All identifiers |
| Comments & docs | English | Code-level comments |
| Resource names in URLs | Plural nouns | `/users`, `/orders` |

### Dependency Injection Lifetimes

> **Rule**: Stateless helpers ALWAYS as `Singleton`. Never register a stateless helper as `Scoped`.

### Configuration — External File via Environment Variable (Mandatory)

Secrets and environment-specific settings **must never** be hardcoded or committed to source control inside `appsettings.json`. The mandatory pattern:

1. **All sensitive config** (connection strings, JWT secrets, API keys, external URLs) lives in a file **outside the repository**, on the server.
2. An environment variable (`{APP_NAME}_CONFIG_PATH` or equivalent) points to that file's absolute path.
3. `appsettings.json` in the repo contains only non-sensitive defaults (log levels, feature flags with safe defaults).

```csharp
// Program.cs — load external config file via env var
var configPath = Environment.GetEnvironmentVariable("APP_CONFIG_PATH")
    ?? throw new InvalidOperationException(
        "APP_CONFIG_PATH environment variable is not set. " +
        "Set it to the absolute path of the external configuration file.");

builder.Configuration.AddJsonFile(configPath, optional: false, reloadOnChange: true);
```

> **`reloadOnChange: true`**: ASP.NET Core's built-in file watcher hot-reloads the config when the external file changes — no restart needed for config updates.

**Prohibited — will be rejected in code review:**

```csharp
// ❌ NEVER — secrets in appsettings.json (committed to git)
// appsettings.json: { "Jwt": { "Secret": "my-real-secret" } }

// ❌ NEVER — hardcoded connection string in code
var cs = "Server=prod-server;Database=MyDb;User=sa;Password=12345;";

// ❌ NEVER — config singleton reading from a hardcoded path
var config = MyConfig.Instance.Settings.ConnectionString;
```

**Correct `appsettings.json` (safe to commit):**
```json
{
  "Serilog": { "MinimumLevel": "Information" },
  "AllowedHosts": "*"
}
```

> **External file location convention**: Use an absolute path on the server outside the deployment directory (e.g., `"<provisioned-path>/config.json"` or `"/path/set-at-provisioning-time/config.json"`). The exact path is defined by your infrastructure tooling — never hardcode it. Never inside `wwwroot` or the app folder. The environment variable name should follow the pattern `{SERVICE_NAME}_CONFIG_PATH`.

### Middleware Pipeline Order

```csharp
app.UseMiddleware<CorrelationIdMiddleware>(); // 0. FIRST — ensures correlationId is present in ALL logs
app.UseAuthentication();    // 1. Verify identity
app.UseMiddleware<UserLogContextMiddleware>(); // 1.5 (optional) only when per-user log partitioning is active; see ../logging/SKILL.md
app.UseRateLimiter();       // 2. Protect from abuse — requires builder.Services.AddRateLimiter(...)
app.UseCors();              // 3. Cross-origin policy
app.UseMiddleware<ExceptionHandlingMiddleware>(); // 4. Global exception handler — after auth so UserId is available
app.UseAuthorization();     // 5. Check permissions
app.MapControllers();       // 6. Route to controllers
```

> **Canonical error contract**: Exception mapping, `ProblemDetails` format, and the 9 traceability fields live in [`logging/SKILL.md`](../logging/SKILL.md).
>
> **Optional middleware slot**: `UseMiddleware<UserLogContextMiddleware>()` belongs between `UseAuthentication()` and `UseRateLimiter()` only when per-user log partitioning is active. See [`../logging/SKILL.md`](../logging/SKILL.md).
>
> **UseRateLimiter before UseCors**: CORS preflight (`OPTIONS`) requests that hit the rate limit return `429` without CORS headers, so the browser reports a CORS error. Consider exempting `OPTIONS` requests from rate limiting if this is a concern.
>
> **ExceptionHandlingMiddleware tradeoff**: It is intentionally positioned before `UseAuthorization()` so authenticated `UserId` is already available for error logging. Authorization policy handler exceptions are **not** caught here — they surface as unhandled `500` responses unless you also add a fallback middleware after `UseAuthorization()`.

### Interface-First — All Injectable Dependencies Have an Interface

```csharp
public interface IUserService { Task<UserDto> GetByIdAsync(Guid id); }
public class UserService : IUserService { ... }
```

### Custom Exception — Domain Error Modeling

```csharp
// Shared/Exceptions/NotFoundException.cs
public class NotFoundException : Exception
{
    public NotFoundException(string resource, Guid id)
        : base($"{resource} with id '{id}' was not found.") { }
}
```

```csharp
// Shared/Exceptions/ConflictException.cs
public class ConflictException : Exception
{
    public ConflictException(string message) : base(message) { }
}
```


## Skill Family — When to Load Each One

| Domain | Skill | When to Load |
|---------|-------|---------------|
| Data access, EF Core, DbContext, repositories, migrations | [`dataaccess`](../dataaccess/SKILL.md) | Configuring DbContext, writing repositories |
| Unit tests, mocks, assertions, test data builders | [`testing-unit`](../testing-unit/SKILL.md) | Creating or reviewing unit tests |
| Responses, response DTOs, HTTP error handling | [`responses`](../responses/SKILL.md) | Designing response contracts, exception handling |
| Validations, FluentValidation | [`validations`](../validations/SKILL.md) | Writing or reviewing validation rules |
| Mapping, AutoMapper, mapping extensions | [`mapping`](../mapping/SKILL.md) | Configuring AutoMapper profiles or extension methods |
| Security, authentication, authorization, rate limiting | [`security`](../security/SKILL.md) | Implementing auth, JWT, policies, rate limiter |
| Performance, caching, optimization | [`performance`](../performance/SKILL.md) | Resilience patterns (Polly, circuit breaker, retry), Kestrel limits, response caching, async parallelism |
| Requests, input DTOs | [`requests`](../requests/SKILL.md) | Designing request contracts |
| Logging, observability | [`logging`](../logging/SKILL.md) | Configuring Serilog, structured logging |

---

## Changelog

### v1.7 — 2026-04-09
- **Fixed (Round 4)**: Replaced the concrete `/etc/myapp/config.json` example with provisioned-path placeholders and an explicit note that infrastructure tooling defines the exact absolute path.
- **Fixed (Round 4)**: Added the pipeline-order tradeoff note explaining that `ExceptionHandlingMiddleware` runs before `UseAuthorization()` for `UserId` logging, so authorization handler exceptions are not caught unless a fallback middleware is added later in the pipeline.

### v1.6 — 2026-04-09
- **Fixed (Round 3)**: Added the optional `UserLogContextMiddleware` slot between `UseAuthentication()` and `UseRateLimiter()` in the canonical pipeline, with a cross-ref to `logging/SKILL.md`.
- **Fixed (Round 3)**: Added a note explaining the `UseRateLimiter()` before `UseCors()` side-effect: rate-limited preflight `OPTIONS` requests return `429` without CORS headers, which browsers surface as CORS errors. Documented the `OPTIONS` exemption tradeoff.

### v1.5 — 2026-04-09
- **Added**: External configuration file pattern via environment variable — mandatory for all sensitive settings. Explicit prohibition on secrets in `appsettings.json`, hardcoded strings, and config singletons with static paths. Pattern extracted from production microservices audit.

### v1.4 — 2026-04-09
- **Fixed (W-13)**: Removed PascalCase and camelCase rows from Naming Conventions table — standard C# conventions the agent already knows. Kept only the ZeshOne-specific decisions: `_camelCase` for private fields, English for all code, plural nouns for URLs.

### v1.3 — 2026-03-28
- **Removed**: SOLID principles section (agent already knows this)
- **Removed**: Async/await tutorial, LINQ over loops, HTTP verb semantics table (language standard)
- **Removed**: Controller design boilerplate (thin controller is baseline knowledge)
- **Kept**: Project structure, naming conventions, DI lifetimes rule, pipeline order, skill family table

### v1.2 — 2026-03-25
- Pipeline order — CorrelationIdMiddleware moved to position 0

### v1.1 — 2026-03-24
- Skill Family table added; DI lifetimes, naming conventions
