---
name: net8-apirest-general
description: >
  General architecture and coding conventions for ASP.NET Core 8 REST APIs using Clean Architecture with Vertical Slices.
  Trigger: When creating, scaffolding, or structuring any .NET 8 REST API project, controllers, services, repositories, or feature folders.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
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

Each feature is a self-contained module:

```
Features/
  {Feature}/
    Controllers/   ← HTTP endpoints, thin layer
    Services/      ← Business logic (interface + implementation)
    Repositories/  ← Data access (interface + implementation)
    DTOs/          ← Request/Response transfer objects
    Models/        ← Domain entities (EF Core mapped)
    Mappings/      ← AutoMapper profiles or manual extension methods
    Validators/    ← FluentValidation rules
    Exceptions/    ← Domain-specific custom exceptions

Shared/
  Models/          ← ResponseDTO<T> and shared models
  Models/BaseEntity.cs  ← Id, CreatedAt, UpdatedAt (ver dataaccess/SKILL.md)
  Helpers/         ← Global utilities (e.g., PasswordHasher)
  Middlewares/     ← Global pipeline middlewares
  Extensions/      ← Extension methods
  Interceptors/    ← EF Core interceptors (e.g., AuditInterceptor — ver dataaccess/SKILL.md)

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
| Comments & docs | Spanish | Code-level comments |
| Resource names in URLs | Plural nouns | `/users`, `/orders` |

### Dependency Injection Lifetimes

| Lifetime | Use case |
|---|---|
| `Singleton` | Stateless helpers, configuration access |
| `Scoped` | DbContext, per-request services |
| `Transient` | Lightweight, stateless services |

> **Rule**: Helpers must be registered as `Singleton`. Never register a stateless helper as `Scoped`.

### Interface-First Principle

Every injectable service **must** have a corresponding interface:

```csharp
public interface IUserService
{
    Task<UserDto> GetByIdAsync(Guid id);
}

public class UserService : IUserService
{
    private readonly IUserRepository _repository;

    public UserService(IUserRepository repository)
    {
        _repository = repository;
    }

    public async Task<UserDto> GetByIdAsync(Guid id)
    {
        // ...
    }
}
```

### HTTP Verb Semantics

| Verb | Action | Success Code |
|---|---|---|
| `GET` | Read resource or collection | `200 OK` |
| `POST` | Create new resource | `201 Created` |
| `PUT` | Full update of resource | `200 OK` or `204 No Content` |
| `DELETE` | Remove resource | `204 No Content` |

### Controller Design

- Controllers must be **thin**: receive request → delegate to service → return response
- Never place business logic inside controllers
- Use `[ApiController]` and `[Route("api/[controller]")]` on every controller
- Inherit from `ControllerBase` (not `Controller`)

```csharp
[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;

    public UsersController(IUserService userService)
    {
        _userService = userService;
    }

    [HttpGet("{id:guid}")]
    public async Task<IActionResult> GetById(Guid id)
    {
        var result = await _userService.GetByIdAsync(id);
        return Ok(result);
    }
}
```

### Async/Await — Mandatory for All I/O

All database calls, HTTP calls, file system operations, and external service calls **must** be async:

```csharp
// CORRECT
public async Task<User> GetUserAsync(Guid id)
{
    return await _context.Users.FindAsync(id);
}

// WRONG — never block async code
public User GetUser(Guid id)
{
    return _context.Users.FindAsync(id).Result; // deadlock risk
}
```

### LINQ over Imperative Loops

Prefer declarative LINQ for collection operations:

```csharp
// CORRECT
var activeUsers = users.Where(u => u.IsActive).Select(u => u.Name).ToList();

// AVOID for simple operations
var activeUsers = new List<string>();
foreach (var u in users)
{
    if (u.IsActive) activeUsers.Add(u.Name);
}
```

### SOLID Principles

- **SRP**: One class, one responsibility. Keep methods short with a single purpose.
- **OCP**: Extend via interfaces, not by modifying existing classes.
- **LSP**: Subtypes must be substitutable for their base types.
- **ISP**: Prefer small, focused interfaces.
- **DIP**: Depend on abstractions (interfaces), never on concrete implementations.

### Middleware Pipeline Order

```csharp
app.UseMiddleware<CorrelationIdMiddleware>(); // 0. Correlation ID — PRIMERO: garantiza correlationId en TODO log
app.UseAuthentication();    // 1. Verify identity
app.UseRateLimiter();       // 2. Protect from abuse — requiere builder.Services.AddRateLimiter(...), ver ../security/SKILL.md
app.UseCors();              // 3. Cross-origin policy
app.UseMiddleware<ExceptionHandlingMiddleware>(); // 4. Global exception handler — después de auth para tener UserId disponible
app.UseAuthorization();     // 5. Check permissions
app.MapControllers();       // 6. Route to controllers
```

> **Contrato canónico de errores**: El mapeo de excepciones (`ValidationException`→400, `NotFoundException`→404, etc.), el formato `ProblemDetails` con `correlationId` y los 9 campos mínimos de trazabilidad viven en [`logging/SKILL.md`](../logging/SKILL.md). No duplicar esas tablas aquí.

---

## Code Examples

### Program.cs — Service Registration Pattern

```csharp
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddDbContextPool<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")));
// → ver ../dataaccess/SKILL.md para configuración robusta (poolSize, null-check, SqlConnectionStringBuilder)

// Feature services
builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddScoped<IUserRepository, UserRepository>();

// Shared singletons
builder.Services.AddSingleton<IPasswordHasher, PasswordHasher>();
builder.Services.AddValidatorsFromAssemblyContaining<Program>();
builder.Services.AddAutoMapper(typeof(Program));

var app = builder.Build();

app.UseMiddleware<CorrelationIdMiddleware>(); // 0 — ver logging/SKILL.md para implementación canónica
app.UseAuthentication();
app.UseRateLimiter(); // requiere builder.Services.AddRateLimiter(...) — ver ../security/SKILL.md
app.UseCors();
app.UseMiddleware<ExceptionHandlingMiddleware>(); // ver logging/SKILL.md para contrato ProblemDetails
app.UseAuthorization();
app.MapControllers();
app.Run();
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

## Skill Family — Cuándo Cargar Cada Una

Esta skill es la **skill madre** de la familia `net8-apirest`. Define arquitectura y convenciones comunes. Para dominios específicos, cargá la skill correspondiente además de esta:

| Dominio | Skill | Cuándo Cargar |
|---------|-------|---------------|
| Acceso a datos, EF Core, DbContext, repositorios, migraciones | [`dataaccess`](../dataaccess/SKILL.md) | Configurando DbContext, escribiendo repositorios, connection strings, interceptors |
| Tests unitarios, mocks, assertions, test data builders | [`testing-unit`](../testing-unit/SKILL.md) | Creando o revisando unit tests de services, validators, mappings |
| Responses, DTOs de respuesta, manejo de errores HTTP | [`responses`](../responses/SKILL.md) | Diseñando contratos de respuesta, exception handling middleware |
| Validaciones, FluentValidation | [`validations`](../validations/SKILL.md) | Escribiendo o revisando reglas de validación |
| Mapping, AutoMapper, extensiones de mapeo | [`mapping`](../mapping/SKILL.md) | Configurando perfiles de AutoMapper o extension methods de mapeo |
| Seguridad, autenticación, autorización, rate limiting | [`security`](../security/SKILL.md) | Implementando auth, JWT, policies, rate limiter |
| Performance, caching, optimización | [`performance`](../performance/SKILL.md) | Optimizando queries, caching, response compression |
| Requests, DTOs de entrada | [`requests`](../requests/SKILL.md) | Diseñando contratos de request |
| Logging, observabilidad | [`logging`](../logging/SKILL.md) | Configurando Serilog, structured logging, health checks |

> **Regla**: Si tu tarea cae en uno de estos dominios, cargá la skill específica ADEMÁS de esta general. La general da el marco; la específica da las reglas detalladas.

---

## Commands

```bash
# Create new .NET 8 Web API project
dotnet new webapi -n MyApi --framework net8.0

# Build solution
dotnet build MySolution.sln

# Run API locally
dotnet run --project src/MyApi

# Add EF Core migration
dotnet ef migrations add InitialCreate --project src/MyApi

# Apply migration
dotnet ef database update --project src/MyApi
```

---

## Resources

- **Standards**: See [../../../../rules-to-skills/Standardized_NET_Rules.md](../../../../rules-to-skills/Standardized_NET_Rules.md)
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

### v1.2 — 2026-03-25
- **Updated**: Pipeline order — `CorrelationIdMiddleware` movido a posición 0 (antes de `Authentication`) para garantizar correlationId en todo log, incluyendo errores de auth.
- **Added**: Nota de cross-reference a `logging/SKILL.md` para contrato canónico de errores (ProblemDetails, mapeo de excepciones, 9 campos de trazabilidad).
- **Updated**: Snippet `Program.cs` alineado con nuevo pipeline order.
- **Meta**: Cambio mínimo bajo SDD formal (net8-apirest-logging, decision D-36/B).

### v1.1 — 2026-03-24
- **Fixed**: Connection string key `"Default"` → `"DefaultConnection"` (alineación con `dataaccess` v2.1)
- **Fixed**: Nota de delegación a `dataaccess` para configuración robusta de `AddDbContextPool`
- **Added**: Sección `Skill Family — Cuándo Cargar Cada Una` con tabla de navegación a 9 skills hijas
- **Added**: `Shared/Interceptors/`, `Shared/Models/BaseEntity.cs`, `tests/` en estructura de proyecto
- **Added**: Nota inline sobre prerequisito de `UseRateLimiter` (→ `security`)
- **Updated**: Sección `Resources` completa con todas las skills de la familia
- **Meta**: Primer cambio bajo SDD formal (engram, OpenCode only)
