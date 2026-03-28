---
name: net8-apirest-requests
description: >
  Best practices for handling HTTP requests in ASP.NET Core 8 REST APIs covering model binding, DTO design, input binding sources, and request pipeline conventions.
  Trigger: When designing request DTOs, configuring model binding, handling query parameters, or processing incoming HTTP payloads in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Designing request DTOs for endpoints
- Handling file uploads or multipart requests
- Reviewing controller action signatures

---

## Critical Patterns

### Request DTO Naming

Always suffix `Request`. Never reuse domain entities as request bodies:

```
CreateUserRequest
UpdateUserRequest
SearchUsersRequest
```

### Pagination — Mandatory cap on `PageSize`

```csharp
public class PaginationRequest
{
    public int Page { get; set; } = 1;
    public int PageSize { get; set; } = 20;

    // Mandatory cap — prevents pageSize=9999
    public int ValidatedPageSize => Math.Min(PageSize, 100);
}

public class SearchUsersRequest : PaginationRequest
{
    public string? Name { get; set; }
    public bool? IsActive { get; set; }
}
```

### Route Constraints — GUID on public resources

```csharp
// Rejects invalid IDs before the service — prevents BOLA enumeration
[HttpGet("{id:guid}")]
```

Never `{id:int}` on public endpoints — see [`../security/SKILL.md`](../security/SKILL.md).

### Binding Sources — Always explicit

Do not rely on ASP.NET Core's implicit inference. Always annotate:

```csharp
[HttpGet("{id:guid}")]
public async Task<IActionResult> GetById(
    [FromRoute] Guid id,
    [FromHeader(Name = "X-Tenant-ID")] string? tenantId) { ... }

[HttpGet]
public async Task<IActionResult> Search([FromQuery] SearchUsersRequest request) { ... }

[HttpPost]
public async Task<IActionResult> Create([FromBody] CreateUserRequest request) { ... }
```

### Null Safety in Request DTOs

Required fields non-nullable with default. Optional fields nullable:

```csharp
public class CreateUserRequest
{
    public string FirstName { get; set; } = string.Empty;   // required
    public string Email { get; set; } = string.Empty;       // required
    public string? PhoneNumber { get; set; }                // optional
}
```

### File Uploads

The file upload validation pattern (D-30) lives in [`../validations/SKILL.md`](../validations/SKILL.md).

> **Critical note**: `[Consumes("multipart/form-data")]` at the action level **overrides** the `[Consumes("application/json")]` on the controller. Always annotate it on file upload actions — without it ASP.NET Core rejects the request with 415.

---

## Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| Domain entity as request body | Couples the API to the data model; security risk |
| Binding source without annotation | Implicit inference can fail in non-obvious ways |
| Pagination without `PageSize` cap | `pageSize=9999` is a DoS vector |
| `int` IDs on public endpoints | Trivial BOLA enumeration |
| `ResponseDTO<T>` as HTTP response body | Violates D-25/D-26 — return raw resource or `ProblemDetails` |
| Reading `IFormFile` stream in validator | Blocks the stream — see D-30 in `validations` |

---

## Resources

- **General structure**: See [../general/SKILL.md](../general/SKILL.md)
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md)
- **Validations (file upload D-30)**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Responses**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Security (BOLA, GUID IDs)**: See [../security/SKILL.md](../security/SKILL.md)

---

## Changelog

### v1.2 — 2026-03-28
- **Removed**: Full binding sources tutorial with table and examples (the agent already knows `[FromBody]`, `[FromQuery]`, etc.)
- **Removed**: Content negotiation `[Produces]`/`[Consumes]` boilerplate — framework standard
- **Removed**: Extended file upload section — consolidated in `validations/SKILL.md` (D-30), only referenced here
- **Kept**: Naming convention, pagination with ValidatedPageSize cap, GUID route constraint, explicit binding, null safety, anti-patterns

### v1.1 — 2026-03-24
- File uploads rewritten with D-30 pattern; anti-patterns expanded
