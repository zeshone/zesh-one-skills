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
| Comments & docs | Spanish | Code-level comments |
| Resource names in URLs | Plural nouns | `/users`, `/orders` |

### Dependency Injection Lifetimes

| Lifetime | Use case |
|---|---|
| `Singleton` | Stateless helpers, configuration access |
| `Scoped` | DbContext, per-request services |
| `Transient` | Lightweight, stateless services |

> **Regla**: Los helpers stateless SIEMPRE como `Singleton`. Nunca registrar un helper sin estado como `Scoped`.

### Middleware Pipeline Order

```csharp
app.UseMiddleware<CorrelationIdMiddleware>(); // 0. PRIMERO — garantiza correlationId en TODO log
app.UseAuthentication();    // 1. Verify identity
app.UseRateLimiter();       // 2. Protect from abuse — requiere builder.Services.AddRateLimiter(...)
app.UseCors();              // 3. Cross-origin policy
app.UseMiddleware<ExceptionHandlingMiddleware>(); // 4. Global exception handler — después de auth para tener UserId
app.UseAuthorization();     // 5. Check permissions
app.MapControllers();       // 6. Route to controllers
```

> **Contrato canónico de errores**: El mapeo de excepciones, el formato `ProblemDetails` y los 9 campos de trazabilidad viven en [`logging/SKILL.md`](../logging/SKILL.md).

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
// → ver ../dataaccess/SKILL.md para configuración robusta

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

## Skill Family — Cuándo Cargar Cada Una

| Dominio | Skill | Cuándo Cargar |
|---------|-------|---------------|
| Acceso a datos, EF Core, DbContext, repositorios, migraciones | [`dataaccess`](../dataaccess/SKILL.md) | Configurando DbContext, escribiendo repositorios |
| Tests unitarios, mocks, assertions, test data builders | [`testing-unit`](../testing-unit/SKILL.md) | Creando o revisando unit tests |
| Responses, DTOs de respuesta, manejo de errores HTTP | [`responses`](../responses/SKILL.md) | Diseñando contratos de respuesta, exception handling |
| Validaciones, FluentValidation | [`validations`](../validations/SKILL.md) | Escribiendo o revisando reglas de validación |
| Mapping, AutoMapper, extensiones de mapeo | [`mapping`](../mapping/SKILL.md) | Configurando perfiles de AutoMapper o extension methods |
| Seguridad, autenticación, autorización, rate limiting | [`security`](../security/SKILL.md) | Implementando auth, JWT, policies, rate limiter |
| Performance, caching, optimización | [`performance`](../performance/SKILL.md) | Optimizando queries, caching, response compression |
| Requests, DTOs de entrada | [`requests`](../requests/SKILL.md) | Diseñando contratos de request |
| Logging, observabilidad | [`logging`](../logging/SKILL.md) | Configurando Serilog, structured logging |

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
- **Removed**: SOLID principles section (agente ya lo conoce)
- **Removed**: Async/await tutorial, LINQ over loops, HTTP verb semantics table (estándar del lenguaje)
- **Removed**: Controller design boilerplate (thin controller es conocimiento base)
- **Kept**: Project structure, naming conventions, DI lifetimes rule, pipeline order, skill family table

### v1.2 — 2026-03-25
- Pipeline order — CorrelationIdMiddleware movido a posición 0

### v1.1 — 2026-03-24
- Skill Family table agregada; DI lifetimes, naming conventions
