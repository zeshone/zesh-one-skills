# Backend Domain — .NET 8 REST API Development Guidelines

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
| `net8-apirest-general` | Clean Architecture, Vertical Slices, project structure, naming conventions | [SKILL.md](skills/backend/general/SKILL.md) |
| `net8-apirest-dataaccess` | Entity Framework Core, DbContext pooling, repositories, migrations, audit patterns | [SKILL.md](skills/backend/dataaccess/SKILL.md) |
| `net8-apirest-security` | OWASP API Top 10, JWT, authorization, BOLA prevention, secure credential handling | [SKILL.md](skills/backend/security/SKILL.md) |
| `net8-apirest-validations` | FluentValidation, input validation, decision trees, error handling | [SKILL.md](skills/backend/validations/SKILL.md) |
| `net8-apirest-requests` | Request DTOs, model binding, query parameters, file uploads | [SKILL.md](skills/backend/requests/SKILL.md) |
| `net8-apirest-responses` | ResponseDTO<T>, HTTP status codes, error formatting, inter-layer communication | [SKILL.md](skills/backend/responses/SKILL.md) |
| `net8-apirest-mapping` | AutoMapper vs manual mapping, mapping profiles, extension patterns | [SKILL.md](skills/backend/mapping/SKILL.md) |
| `net8-apirest-logging` | Serilog, structured logging, Correlation IDs, exception handling | [SKILL.md](skills/backend/logging/SKILL.md) |
| `net8-apirest-performance` | EF Core optimization, async patterns, Polly resilience, DbContextPool | [SKILL.md](skills/backend/performance/SKILL.md) |
| `net8-apirest-testing-unit` | Unit testing with xUnit, NSubstitute, FluentAssertions, test data builders | [SKILL.md](skills/backend/testing-unit/SKILL.md) |

### Shared Skills
| Skill | Description | URL |
|-------|-------------|-----|
| `github-pr` | Pull Request conventions, branch workflow, PR checklist | [SKILL.md](skills/shared/github-pr/SKILL.md) |

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
| Creating or reviewing unit tests for services, validators, or mappings | `net8-apirest-testing-unit` |
| Deciding between AutoMapper and manual mapping | `net8-apirest-mapping` |
| Implementing rate limiting or request size limits | `net8-apirest-security` |
| Setting up async patterns or parallelism | `net8-apirest-performance` |
| Handling file uploads or multipart requests | `net8-apirest-requests` |
| Reviewing OWASP compliance or security coverage | `net8-apirest-security` |
| Creating or reviewing a Pull Request | `github-pr` |
