# .NET 8 REST API Development Guidelines

## How to Use This Guide

- Start here for cross-project norms on ASP.NET Core 8 REST API development.
- Each domain has a dedicated skill with detailed patterns (e.g., `general/SKILL.md`, `dataaccess/SKILL.md`).
- Domain skills override this file when guidance conflicts.

---

## Available Skills

Use these skills for detailed patterns on-demand:

### .NET 8 REST API Skills
| Skill | Description | URL |
|-------|-------------|-----|
| `net8-apirest-general` | Clean Architecture, Vertical Slices, project structure, naming conventions | [SKILL.md](general/SKILL.md) |
| `net8-apirest-dataaccess` | Entity Framework Core, DbContext pooling, repositories, migrations, audit patterns | [SKILL.md](dataaccess/SKILL.md) |
| `net8-apirest-security` | OWASP API Top 10, JWT, authorization, BOLA prevention, secure credential handling | [SKILL.md](security/SKILL.md) |
| `net8-apirest-validations` | FluentValidation, input validation, decision trees, error handling | [SKILL.md](validations/SKILL.md) |
| `net8-apirest-requests` | Request DTOs, model binding, query parameters, file uploads | [SKILL.md](requests/SKILL.md) |
| `net8-apirest-responses` | ResponseDTO<T>, HTTP status codes, error formatting, inter-layer communication | [SKILL.md](responses/SKILL.md) |
| `net8-apirest-mapping` | AutoMapper vs manual mapping, mapping profiles, extension patterns | [SKILL.md](mapping/SKILL.md) |
| `net8-apirest-logging` | Serilog, structured logging, Correlation IDs, exception handling | [SKILL.md](logging/SKILL.md) |
| `net8-apirest-performance` | EF Core optimization, async patterns, Polly resilience, DbContextPool | [SKILL.md](performance/SKILL.md) |

---

## Auto-invoke Skills

When performing these actions, **ALWAYS** invoke the corresponding skill FIRST:

| Action | Skill |
|--------|-------|
| Creating a new .NET 8 API project or feature | `net8-apirest-general` |
| Scaffolding controllers, services, repositories | `net8-apirest-general` |
| Creating or modifying DbContext, repositories, or migrations | `net8-apirest-dataaccess` |
| Implementing JWT, OAuth, authorization policies | `net8-apirest-security` |
| Adding security controls (CORS, rate limiting, BOLA prevention) | `net8-apirest-security` |
| Designing request DTOs or configuring model binding | `net8-apirest-requests` |
| Building API responses or response DTOs | `net8-apirest-responses` |
| Creating FluentValidation validators or validation rules | `net8-apirest-validations` |
| Creating AutoMapper profiles or mapping logic | `net8-apirest-mapping` |
| Configuring Serilog, adding logging, or exception handling | `net8-apirest-logging` |
| Optimizing EF Core queries or implementing resilience patterns | `net8-apirest-performance` |
| Reviewing or refactoring data access code | `net8-apirest-dataaccess` |
| Deciding between AutoMapper and manual mapping | `net8-apirest-mapping` |
| Implementing rate limiting or request size limits | `net8-apirest-security` & `net8-apirest-performance` |
| Setting up async patterns or parallelism | `net8-apirest-performance` |
| Handling file uploads or multipart requests | `net8-apirest-requests` |
| Reviewing OWASP compliance or security coverage | `net8-apirest-security` |

---

## Project Overview

### Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Framework | ASP.NET Core | 8.0+ |
| Language | C# | 12.0+ |
| Database ORM | Entity Framework Core | 8.0+ |
| Validation | FluentValidation | 11.0+ |
| Mapping | AutoMapper | 13.0+ |
| Logging | Serilog | 7.0+ |
| Resilience | Polly | 8.0+ |
| Authentication | JWT Bearer | Built-in |
| Package Manager | NuGet or dotnet CLI | Latest |

### Architecture

- **Clean Architecture** with **Vertical Slices** per feature
- Cross-cutting concerns (logging, validation, exception handling) managed globally
- Repository Pattern for data access
- Service Layer for business logic
- OWASP API Security Top 10 compliance mandatory
- Structured logging via Serilog with Correlation IDs

### Project Structure

```
src/
  YourApi.csproj
  Program.cs
  
Features/
  {Feature}/
    Controllers/          ← HTTP endpoints (thin)
    Services/             ← Business logic
    Repositories/         ← Data access
    DTOs/                 ← Request/Response transfer objects
    Models/               ← Domain entities
    Validators/           ← FluentValidation rules
    Mappings/             ← AutoMapper profiles
    Exceptions/           ← Domain-specific exceptions

Shared/
  Models/                 ← ResponseDTO<T>, shared dtos
  Helpers/                ← Global utilities
  Middlewares/            ← Global pipeline middlewares
  Extensions/             ← Extension methods

Database/
  Context/                ← DbContext
  Migrations/             ← EF Core migrations
```

---

## Development Setup

### Prerequisites

- **.NET 8 SDK** (download from https://dotnet.microsoft.com/download)
- **Visual Studio 2022** or **VS Code** with C# extension
- **SQL Server** (LocalDB, Express, or Azure SQL)
- **Git** for version control

### Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd <api-repo>

# Restore dependencies
dotnet restore

# Build the solution
dotnet build

# Apply migrations (if database exists)
dotnet ef database update --project Database

# Run the API
dotnet run --project src/YourApi.csproj
```

### Database Migrations

```bash
# Create a new migration
dotnet ef migrations add MigrationName --project Database --startup-project src/YourApi.csproj

# Apply migrations
dotnet ef database update --project Database --startup-project src/YourApi.csproj

# Revert last migration
dotnet ef migrations remove --project Database --startup-project src/YourApi.csproj
```

---

## Code Quality & Standards

### Naming Conventions

| Item | Convention | Example |
|------|-----------|---------|
| Namespaces | `CompanyName.ProjectName.Feature` | `ClickSeguros.Api.Users` |
| Classes | PascalCase | `UserService`, `GetUserByIdRequest` |
| Methods | PascalCase | `GetByIdAsync()`, `CreateAsync()` |
| Properties | PascalCase | `FirstName`, `UserId` |
| Private fields | `_camelCase` | `_repository`, `_logger` |
| Local variables | `camelCase` | `userId`, `newUser` |
| Constants | `UPPER_CASE` | `MAX_PAGE_SIZE`, `DEFAULT_TIMEOUT` |
| Request DTOs | `{Action}{Resource}Request` | `CreateUserRequest`, `UpdateOrderRequest` |
| Response DTOs | `{Resource}Dto` or `{Resource}Response` | `UserDto`, `OrderResponse` |

### Async/Await Standards

- **Always** use async methods in repositories and services (`async Task`, `async Task<T>`)
- **Always** use `.ConfigureAwait(false)` in library code
- **Never** use `.Result` or `.Wait()` (deadlock risk)
- Use `IAsyncEnumerable<T>` for large dataset streaming

### HTTP Status Codes

| Status | Use Case | Example |
|--------|----------|---------|
| 200 | Successful GET, PUT | Fetched user successfully |
| 201 | Successful POST (resource created) | User created |
| 202 | Accepted (async operation) | Background job queued |
| 204 | Successful DELETE, no content | Resource deleted |
| 400 | Bad Request (validation failed) | Missing required field |
| 401 | Unauthorized (no auth provided) | Missing JWT token |
| 403 | Forbidden (insufficient permissions) | User lacks role |
| 404 | Not Found | Resource does not exist |
| 409 | Conflict (e.g., duplicate email) | Email already registered |
| 500 | Internal Server Error | Unhandled exception |

---

## Commit & Pull Request Guidelines

### Conventional Commits

Follow the format: `<type>[scope]: <description>`

**Types:**
- `feat` — New feature or endpoint
- `fix` — Bug fix
- `docs` — Documentation changes
- `chore` — Build, dependencies, utilities (no production code change)
- `refactor` — Code structure change without feature/fix
- `perf` — Performance improvement
- `test` — Test-related changes
- `style` — Code style, formatting (no logic change)

**Examples:**
```
feat(users): add endpoint to list users with pagination
fix(auth): prevent BOLA on user profile endpoint
docs: update API documentation for auth flow
chore(deps): upgrade Entity Framework Core to 8.0.1
refactor(database): optimize repository query patterns
perf: add DbContextPool to reduce allocation overhead
test: add unit tests for UserService validation
```

### Pull Request Checklist

Before submitting a PR, verify:

- [ ] Code follows naming conventions and architecture patterns
- [ ] FluentValidation used for input validation
- [ ] Security checks implemented (BOLA, authorization, CORS)
- [ ] Logging added for important business operations and errors
- [ ] Database queries optimized (includes, projections, no N+1 problems)
- [ ] Async/await used correctly (no `.Result` or `.Wait()`)
- [ ] Unit tests added for new business logic
- [ ] Integration tests added for database/external service calls
- [ ] All tests pass: `dotnet test`
- [ ] Code builds without warnings: `dotnet build /p:TreatWarningsAsErrors=true`
- [ ] Database migrations included (if schema changes)
- [ ] Documentation updated (README, API docs, architecture decisions)

---

## Best Practices Summary

### General

| Practice | Rationale |
|----------|-----------|
| Use Vertical Slices (features are self-contained modules) | Easier to locate, test, and maintain code |
| Dependency Injection via constructor | Testability, loose coupling |
| Use DTOs at API boundaries (never expose entities) | Security, backward compatibility, flexibility |
| Immutable DTOs (readonly properties or records) | Safer concurrent access, less bugs |

### Data Access

| Practice | Rationale |
|----------|-----------|
| Always use `AddDbContextPool` (never `AddDbContext`) | Reduces allocation pressure, reuses context instances |
| Use `.AsNoTracking()` for read-only queries | Performance improvement, lower memory overhead |
| Implement soft delete via global query filters | Preserve data history, audit compliance |
| Use stored procedures only for complex scenarios | SQL injection risk, slower to evolve |
| Parameterize all raw SQL queries | Prevents SQL injection |

### Security

| Practice | Rationale |
|----------|-----------|
| Always verify resource ownership (BOLA check) | Prevents unauthorized data access |
| Use GUID for IDs (never sequential integers) | Reduces enumeration attacks |
| Hash passwords with PBKDF2 or bcrypt (never plaintext or MD5) | Compliance, breach protection |
| Implement rate limiting per endpoint | DDoS prevention, API abuse prevention |
| Use short-lived JWT tokens (5-15 min) + refresh tokens | Limits exposure of compromised tokens |
| Validate CORS origins whitelist | Cross-origin attack prevention |

### Logging

| Practice | Rationale |
|----------|-----------|
| Use structured logging (Serilog with properties) | Easier to parse, aggregate, search logs |
| Add Correlation ID to all requests | Trace user requests across services |
| Log at appropriate levels (Info, Warning, Error, Debug) | Easier to filter and alert |
| Never log sensitive data (passwords, PII) | Compliance, data protection |
| Log business-critical operations (login, payment, deletion) | Audit trail, compliance |

---

## Additional Resources

- [Microsoft — ASP.NET Core Documentation](https://learn.microsoft.com/en-us/aspnet/core/)
- [Entity Framework Core Docs](https://learn.microsoft.com/en-us/ef/core/)
- [OWASP API Security Top 10 2023](https://owasp.org/www-project-api-security/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [RESTful API Best Practices](https://restfulapi.net/)
