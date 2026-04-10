---
name: net8-apirest-responses
description: >
  Unified response standards for ASP.NET Core 8 REST APIs. Covers two internal contracts:
  `Result<T>` (preferred — Railway-Oriented) and `ResponseDTO<T>` (legacy), plus the external
  HTTP contract (raw resource JSON or `ProblemDetails`) for controller → client responses.
  Trigger: When building API responses, defining return types for services, or standardizing
  how data and errors flow from the data layer to the client in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.7"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Defining return types for services and repositories
- Standardizing success and error responses
- Implementing exception handling middleware
- Reviewing the data flow from the data layer to the client

---

## Critical Patterns

### Internal Contract — Choose One

| Contract | Status | When to use |
|---|---|---|
| `Result<T>` | ✅ **Preferred** | New projects and new services |
| `ResponseDTO<T>` | ⚠️ **Legacy** | Existing codebases already using it — do not mix both in the same service |

---

### `Result<T>` — Railway-Oriented Pattern (Preferred)

`Result<T>` makes success and failure explicit return types — no exceptions for business failures, no null checks, no internal envelope leaking to the client.

```csharp
// Shared/Models/Result.cs
public class Result<T>
{
    public bool IsSuccess { get; private set; }
    public T? Value { get; private set; }
    public Uri? Location { get; private set; }
    public string? ErrorMessage { get; private set; }
    public int StatusCode { get; private set; }
    public string? Title { get; private set; }
    public Dictionary<string, string[]>? ValidationErrors { get; private set; }

    private Result() { }

    public static Result<T> Success(T value) =>
        new() { IsSuccess = true, Value = value, StatusCode = 200 };

    public static Result<T> Created(T value, Uri? location = null) =>
        new() { IsSuccess = true, Value = value, StatusCode = 201, Location = location };

    public static Result<T> NoContent() =>
        new() { IsSuccess = true, StatusCode = 204 };

    // WARNING: Never pass raw exception messages from unhandled exceptions here
    // (SQL text, table names, internal details). Provide a sanitized message explicitly.
    // Unhandled exceptions should usually flow to ExceptionHandlingMiddleware instead.
    public static Result<T> Fail(string message, int statusCode = 500) =>
        new() { IsSuccess = false, ErrorMessage = message, StatusCode = statusCode };

    public static Result<T> BusinessFail(string message, int statusCode = 400, string title = "Bad Request") =>
        new() { IsSuccess = false, ErrorMessage = message, StatusCode = statusCode, Title = title };

    public static Result<T> ValidationFail(Dictionary<string, string[]> errors) =>
        new() { IsSuccess = false, StatusCode = 400, Title = "Validation Failed", ValidationErrors = errors };
```

> **`ToHttpResponse()` contract**: `201 → CreatedResult(Location, Value)`; `204 → NoContentResult`; `ValidationErrors → ValidationProblemDetails(400)`; all other failures → `ProblemDetails` with the factory's `StatusCode` and `Title`. `Location` on 201 is required per REST — see anti-pattern table.

```csharp
    public IActionResult ToHttpResponse() => IsSuccess
        ? StatusCode == 204
            ? new NoContentResult()
            : StatusCode == 201
            ? new CreatedResult(Location?.ToString() ?? string.Empty, Value)
            : new OkObjectResult(Value)
        : ValidationErrors is not null
            ? new BadRequestObjectResult(new ValidationProblemDetails(ValidationErrors) { Status = 400 })
            : new ObjectResult(new ProblemDetails
            {
                Status = StatusCode,
                Title = Title ?? "Error",
                Detail = ErrorMessage
            }) { StatusCode = StatusCode };
}
```

**Controller pattern — always:**
```csharp
public async Task<IActionResult> CreateUser([FromBody] CreateUserRequest request)
{
    var result = await _userService.CreateAsync(request);
    return result.ToHttpResponse();
}
```

**Service pattern:**
```csharp
public async Task<Result<UserDto>> CreateAsync(CreateUserRequest request)
{
    if (await _repository.ExistsByEmailAsync(request.Email))
        return Result<UserDto>.BusinessFail("Email already registered.", 409, "Conflict");

    var entity = request.ToEntity();
    await _repository.AddAsync(entity);
    return Result<UserDto>.Created(entity.ToDto(), new Uri($"/users/{entity.Id}", UriKind.Relative));
}
```

**DELETE convention (`204 No Content`):**
```csharp
public async Task<Result<object>> DeleteAsync(Guid id)
{
    var user = await _repository.GetByIdAsync(id);
    if (user is null) throw new NotFoundException(nameof(User), id);

    await _repository.DeleteAsync(id);
    return Result<object>.NoContent();
}
```

> **Convention**: use `Result<object>.NoContent()` for DELETE operations when you need a generic placeholder type. A `204` response has no body, so the generic type is never serialized to the client.

> **Rule**: Domain exceptions (`NotFoundException`, `ForbiddenException`) are still thrown for exceptional flows and caught by `ExceptionHandlingMiddleware`. `Result<T>` handles **expected business outcomes** (not found by design, conflict, validation) — replacing the throw/catch cycle for predictable conditions.

### Layer Contract with `Result<T>`

| Layer | Returns | Never |
|---|---|---|
| Repository | Domain entities or `null` | DTOs, `Result<T>`, `ResponseDTO<T>` |
| Service | `Result<T>` wrapping DTO | Raw entities, null returns for not-found, silently swallowed errors |
| Controller | `result.ToHttpResponse()` | The `Result<T>` object directly to the client |

```
Repository  →        Service         →   Controller  →   Client
  Entity       Result<T>.Success(dto)    ToHttpResponse()   HTTP 2xx (raw JSON)
  null         Result<T>.BusinessFail()                     HTTP 4xx (ProblemDetails)
```

---

### `ResponseDTO<T>` — Legacy Contract

> ⚠️ **Legacy**: Use only in codebases already built on this pattern. Do NOT introduce in new services. Do NOT mix `Result<T>` and `ResponseDTO<T>` in the same service.
>
> ⛔ **Never add an `Exception?` property to `ResponseDTO<T>`** — exceptions serialized as JSON leak stack traces and internal details to the client.

```csharp
// Shared/Models/ResponseDTO.cs
public class ResponseDTO<T>
{
    public bool Success { get; set; }
    public string Message { get; set; } = string.Empty;
    public T? Data { get; set; }
    // ⛔ NEVER: public Exception? Exception { get; set; }
}
```

**Legacy layer contract:**

| Boundary | Contract | Who uses it |
|---|---|---|
| **Internal** (Service → Controller) | `ResponseDTO<T>` | Services |
| **External** (Controller → Client) | Raw HTTP: resource JSON (2xx) or `ProblemDetails` (4xx/5xx) | Controllers |

> **Rule**: `ResponseDTO<T>` is an internal transport contract. It is **never** serialized as an HTTP response body.

**Usage rules:**
- `Success: true` → `Data` is the result; **never `null`**. Use empty values: `[]`, `string.Empty`, `0`.
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
| `ConflictException` | `general` skill — `Shared/Exceptions/ConflictException.cs` |

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
| `ResponseDTO<T>` in the HTTP response body | Internal envelope exposed to client — violates the external contract |
| `Result<T>` and `ResponseDTO<T>` in the same service | Two conflicting contracts create unpredictable controller code |
| `Exception?` property on any response DTO | Serializes stack trace and internals to the client — security risk |
| Returning `null` from service on "not found" | Use `Result<T>.BusinessFail("Not found.", 404, "Not Found")` or throw `NotFoundException` — never null |
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

### v1.6 — 2026-04-09
- **Fixed (Round 4)**: Changed `Result<T>.Fail(Exception ex)` to `Result<T>.Fail(string message, int statusCode = 500)` so callers must provide a sanitized message explicitly instead of leaking `ex.Message` by default.
- **Fixed (Round 4)**: Corrected the `Layer Contract with Result<T>` table — the Service row now documents what to avoid (`Raw entities, null returns for not-found, silently swallowed errors`) instead of incorrectly saying services should throw for business failures.

### v1.5 — 2026-04-09
- **Fixed (Round 3)**: Added a warning on `Result<T>.Fail(...)` — never forward raw unhandled exception messages to the client; use a generic message and let `ExceptionHandlingMiddleware` sanitize unexpected failures.
- **Fixed (Round 3)**: Added a minimal DELETE service example using `Result<object>.NoContent()` plus a note explaining the placeholder-type convention for `204 No Content` responses.
- **Fixed (Round 3)**: Corrected the anti-pattern example from `.BusinessFail(404)` to `.BusinessFail("Not found.", 404, "Not Found")` so the status code is no longer passed into the message parameter.

### v1.4 — 2026-04-09
- **Added**: `Result<T>` Railway-Oriented pattern as the preferred internal contract for new projects — factory methods `Success`, `Created`, `NoContent`, `Fail`, `BusinessFail`, `ValidationFail`; `ToHttpResponse()` for controller integration; service and controller usage patterns.
- **Updated**: `ResponseDTO<T>` marked as legacy — do not introduce in new services, do not mix with `Result<T>`. Added explicit warning against `Exception?` property (stack trace serialization risk).
- **Updated**: Layer contract table updated to reflect `Result<T>` as the primary path.
- **Updated**: Anti-patterns table expanded with `Result<T>`/`ResponseDTO<T>` mixing and `Exception?` DTO field.
- **Updated**: Frontmatter description updated to reflect dual-contract evolution.

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
