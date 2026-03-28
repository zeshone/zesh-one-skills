---
name: net8-apirest-general
description: >
  General architecture and coding conventions for ASP.NET Core 8 REST APIs using Clean Architecture with Vertical Slices.
  Trigger: When creating, scaffolding, or structuring any .NET 8 REST API project, controllers, services, repositories, or feature folders.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.3"
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
  Models/BaseEntity.cs  ← Id, CreatedAt, UpdatedAt (ver dataaccess/SKILL.md)
  Helpers/         ← Global utilities (e.g., PasswordHasher)
  Middlewares/     ← Global pipeline middlewares
  Extensions/      ← Extension methods
  Interceptors/    ← EF Core interceptors (e.g., AuditInterceptor)

Database/
  Context/         ← EF Core DbContext definition

tests/             ← Unit & integration tests (ver testing-unit/SKILL.md)
```

### Naming Conventions

| Element | Convention | Example |
|---|---|---|
| Classes, Interfaces, Methods, Properties | PascalCase | `UserService`, `IUserRepository` |
| Private fields | `_camelCase` | `_logger`, `_repository` |
| Local variables, parameters | camelCase | `userId`, `inputString` |
| Source code | English | All identifiers |
| Comments & docs | English | Code-level comments |
| Resource names in URLs | Plural nouns | `/users`, `/orders` |

### Dependency Injection Lifetimes

| Lifetime | Use case |
|---|---|
| `Singleton` | Stateless helpers, configuration access |
| `Scoped` | DbContext, per-request services |
| `Transient` | Lightweight, stateless services |

> **Rule**: Stateless helpers ALWAYS as `Singleton`. Never register a stateless helper as `Scoped`.

### Middleware Pipeline Order

```csharp
app.UseMiddleware<CorrelationIdMiddleware>(); // 0. FIRST — ensures correlationId is present in ALL logs
app.UseAuthentication();    // 1. Verify identity
app.UseRateLimiter();       // 2. Protect from abuse — requires builder.Services.AddRateLimiter(...)
app.UseCors();              // 3. Cross-origin policy
app.UseMiddleware<ExceptionHandlingMiddleware>(); // 4. Global exception handler — after auth so UserId is available
app.UseAuthorization();     // 5. Check permissions
app.MapControllers();       // 6. Route to controllers
```

> **Canonical error contract**: Exception mapping, `ProblemDetails` format, and the 9 traceability fields live in [`logging/SKILL.md`](../logging/SKILL.md).

### Interface-First — Toda dependencia inyectable tiene interface

```csharp
public interface IUserService { Task<UserDto> GetByIdAsync(Guid id); }
public class UserService : IUserService { ... }
```

### Custom Exception — Domain Error Modeling

```csharp
public class NotFoundException : Exception
{
    public NotFoundException(string resource, Guid id)
        : base($"{resource} with id '{id}' was not found.") { }
}
```

---

## Code Examples

### Program.cs — Service Registration Pattern

```csharp
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddDbContextPool<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")));
// → see ../dataaccess/SKILL.md for robust configuration

builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddSingleton<IPasswordHasher, PasswordHasher>();
builder.Services.AddValidatorsFromAssemblyContaining<Program>();
builder.Services.AddAutoMapper(typeof(Program));

var app = builder.Build();

app.UseMiddleware<CorrelationIdMiddleware>();
app.UseAuthentication();
app.UseRateLimiter();
app.UseCors();
app.UseMiddleware<ExceptionHandlingMiddleware>();
app.UseAuthorization();
app.MapControllers();
app.Run();
```

---

## Skill Family — When to Load Each One

| Domain | Skill | When to Load |
|---------|-------|---------------|
| Data access, EF Core, DbContext, repositories, migrations | [`dataaccess`](../dataaccess/SKILL.md) | Configuring DbContext, writing repositories |
| Unit tests, mocks, assertions, test data builders | [`testing-unit`](../testing-unit/SKILL.md) | Creating or reviewing unit tests |
| Responses, response DTOs, HTTP error handling | [`responses`](../responses/SKILL.md) | Designing response contracts, exception handling |
| Validations, FluentValidation | [`validations`](../validations/SKILL.md) | Writing or reviewing validation rules |
| Mapping, AutoMapper, mapping extensions | [`mapping`](../mapping/SKILL.md) | Configuring AutoMapper profiles or extension methods |
| Security, authentication, authorization, rate limiting | [`security`](../security/SKILL.md) | Implementing auth, JWT, policies, rate limiter |
| Performance, caching, optimization | [`performance`](../performance/SKILL.md) | Optimizing queries, caching, response compression |
| Requests, input DTOs | [`requests`](../requests/SKILL.md) | Designing request contracts |
| Logging, observability | [`logging`](../logging/SKILL.md) | Configuring Serilog, structured logging |

---

## Commands

```bash
dotnet new webapi -n MyApi --framework net8.0
dotnet build MySolution.sln
dotnet run --project src/MyApi
dotnet ef migrations add InitialCreate --project src/MyApi
dotnet ef database update --project src/MyApi
```

---

## Resources

- **Data Access**: See [../dataaccess/SKILL.md](../dataaccess/SKILL.md)
- **Testing**: See [../testing-unit/SKILL.md](../testing-unit/SKILL.md)
- **Responses**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Validations**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md)
- **Security**: See [../security/SKILL.md](../security/SKILL.md)
- **Performance**: See [../performance/SKILL.md](../performance/SKILL.md)
- **Requests**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Logging**: See [../logging/SKILL.md](../logging/SKILL.md)

---

## Changelog

### v1.3 — 2026-03-28
- **Removed**: SOLID principles section (agent already knows this)
- **Removed**: Async/await tutorial, LINQ over loops, HTTP verb semantics table (language standard)
- **Removed**: Controller design boilerplate (thin controller is baseline knowledge)
- **Kept**: Project structure, naming conventions, DI lifetimes rule, pipeline order, skill family table

### v1.2 — 2026-03-25
- Pipeline order — CorrelationIdMiddleware moved to position 0

### v1.1 — 2026-03-24
- Skill Family table added; DI lifetimes, naming conventions
