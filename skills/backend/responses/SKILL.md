---
name: net8-apirest-responses
description: >
  Unified response standards for ASP.NET Core 8 REST APIs covering the dual response contract:
  `ResponseDTO<T>` for internal inter-layer communication (Service → Controller) and standard
  HTTP responses (raw resource JSON or `ProblemDetails`) for external controller → client responses.
  Trigger: When building API responses, defining return types for services and repositories,
  or standardizing how data and errors flow from the data layer to the client in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.3"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Defining return types for services and repositories
- Standardizing success and error responses
- Implementing exception handling middleware
- Reviewing the data flow from the data layer to the client

---

## Critical Patterns

### Dual Response Contract (D-25)

| Boundary | Contract | Who uses it |
|---|---|---|
| **Internal** (Service → Controller) | `ResponseDTO<T>` | Services |
| **External** (Controller → Client) | Raw HTTP: resource JSON (2xx) or `ProblemDetails` (4xx/5xx) | Controllers |

> **Rule**: `ResponseDTO<T>` is an internal transport contract. It is **never** serialized as an HTTP response body for external clients.

### `ResponseDTO<T>` — Definition

```csharp
// Shared/Models/ResponseDTO.cs
public class ResponseDTO<T>
{
    public bool Success { get; set; }
    public string Message { get; set; } = string.Empty;
    public T? Data { get; set; }
}
```

### Layer Contract — What each layer returns

| Layer | Returns | Never |
|---|---|---|
| Repository | Domain entities or `null` | DTOs, `ResponseDTO<T>` |
| Service | DTOs (after mapping) or throws domain exception | Raw entities |
| Controller | `IActionResult` — raw resource JSON (2xx) or `ProblemDetails` (4xx/5xx) | The `ResponseDTO<T>` envelope to the client |

```
Repository  →      Service      →   Controller  →   Client
  Entity     throws / DTO          IActionResult    HTTP 2xx (raw JSON)
                                                    HTTP 4xx/5xx (ProblemDetails)
```

### `ResponseDTO<T>` — Usage Rules (internal)

- `Success: true` → `Data` is the result; **never `null`** when data exists. Use empty values for empty results: `[]` for lists, `string.Empty` for strings, `0` for numerics, `new T()` if parameterless constructor.
- `Success: false` → `Message` is a clear error description; `Data` is `null`.

### Paged Response

```csharp
// Shared/Models/PagedResult.cs
public class PagedResult<T>
{
    public List<T> Items { get; set; } = new();
    public int TotalCount { get; set; }
    public int Page { get; set; }
    public int PageSize { get; set; }
    public int TotalPages => (int)Math.Ceiling((double)TotalCount / PageSize);
}
```

---

## Exception Handling — HTTP Mapping (Single Source of Truth)

> **Ownership**: The exception → HTTP contract lives here (`responses`). The middleware implementation, logging and correlationId live in [`logging`](../logging/SKILL.md).

### Canonical Table

| Exception | HTTP Status | Log Level | Client Message |
|---|---|---|---|
| `NotFoundException` | `404 Not Found` | Warning | `ex.Message` (safe) |
| `ForbiddenException` | `403 Forbidden` | Warning | `ex.Message` (safe) |
| `ConflictException` | `409 Conflict` | Warning | `ex.Message` (safe) |
| `ValidationException` (manual) | `400 Bad Request` | Information | First error message |
| `Exception` (unhandled) | `500 Internal Server Error` | Error | Generic — **never expose stack trace** |

> **`401 Unauthorized` — Pipeline/Auth ownership (D-27)**  
> `401` is NOT a domain exception. It is resolved by the auth pipeline (JWT bearer, `[Authorize]`).  
> **Do not create `UnauthorizedException`.**

### Exception Location Reference

| Exception | Defined in |
|---|---|
| `NotFoundException` | `general` skill — `Shared/Exceptions/NotFoundException.cs` |
| `ForbiddenException` | `security` skill — `Shared/Exceptions/ForbiddenException.cs` |
| `ConflictException` | `Shared/Exceptions/ConflictException.cs` |

### Uniform 400 — `InvalidModelStateResponseFactory`

```csharp
builder.Services.AddControllers()
    .ConfigureApiBehaviorOptions(options =>
    {
        options.InvalidModelStateResponseFactory = context =>
        {
            var firstError = context.ModelState.Values
                .SelectMany(v => v.Errors)
                .Select(e => e.ErrorMessage)
                .FirstOrDefault() ?? "Validation failed.";

            return new BadRequestObjectResult(new ProblemDetails
            {
                Status = StatusCodes.Status400BadRequest,
                Title = "Bad Request",
                Detail = firstError
            });
        };
    });
```

### Uniform 429 — Rate Limiter `OnRejected`

```csharp
options.OnRejected = async (context, ct) =>
{
    context.HttpContext.Response.StatusCode = StatusCodes.Status429TooManyRequests;
    await context.HttpContext.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status429TooManyRequests,
        Title = "Too Many Requests",
        Detail = "Rate limit exceeded. Please try again later."
    }, ct);
};
```

---

## Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| `ResponseDTO<T>` in the HTTP response body | Violates D-25 — the client receives an internal envelope |
| Returning `null` from service on "not found" | The controller must null-check everywhere; better to throw |
| HTTP 200 for errors | Misleads clients — always use correct status codes |
| `ResponseDTO<T>` in repositories | Repositories are the data layer, not the API layer |
| Creating `UnauthorizedException` as a domain exception | `401` belongs to the pipeline — not the domain |
| Not including `Location` header on 201 Created | Violates REST |

---

## Resources

- **General / NotFoundException**: See [../general/SKILL.md](../general/SKILL.md)
- **Requests**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md)
- **Logging / ExceptionHandlingMiddleware**: See [../logging/SKILL.md](../logging/SKILL.md)
- **Validations**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Security / ForbiddenException**: See [../security/SKILL.md](../security/SKILL.md)

---

## Changelog

### v1.3 — 2026-04-09
- **Fixed (CRITICAL)**: Dual Contract table boundary label corrected from "Repository → Service" to "Service → Controller". Repositories return pure entities — they never produce `ResponseDTO<T>`. Updated frontmatter description to match.
- **Fixed (W-12)**: Collapsed three-row ResponseDTO usage table into two bullet rules — cleaner, same information, removes redundant scenario split.

### v1.2 — 2026-03-28
- **Removed**: Full CRUD controller example (already in `general/SKILL.md`)
- **Removed**: Verbose explanation of the "why" behind the dual contract
- **Removed**: Full middleware snippets (they live in `logging/SKILL.md`)
- **Removed**: Repository and Service return pattern examples (consolidated in the layer contract table)
- **Kept**: Dual contract table, ResponseDTO<T> definition, layer contract table, paged result, HTTP status table, exception mapping as single source of truth, uniform 400/429

### v1.1 — 2026-03-24
- Dual Response Contract D-25/D-26/D-27; ProblemDetails; exception ownership clarified
