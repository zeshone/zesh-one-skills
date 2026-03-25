---
name: net8-apirest-responses
description: >
  Unified response standards for ASP.NET Core 8 REST APIs covering the dual response contract:
  `ResponseDTO<T>` for internal inter-layer communication (Repository → Service) and standard
  HTTP responses (raw resource JSON or `ProblemDetails`) for external controller → client responses.
  Trigger: When building API responses, defining return types for services and repositories,
  or standardizing how data and errors flow from the data layer to the client in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.1"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Building controller action return values
- Defining service method return types
- Standardizing success and error responses
- Implementing repository return patterns
- Reviewing how data and errors flow through the layers toward the client

---

## Dual Response Contract (D-25)

This skill enforces a **two-boundary model**. The response shape depends on where in the system you are:

| Boundary | Contract | Who uses it |
|---|---|---|
| **Internal** (Repository → Service) | `ResponseDTO<T>` | Repositories, Services |
| **External** (Controller → Client) | Raw HTTP: resource JSON (2xx) or `ProblemDetails` (4xx/5xx) | Controllers |

> **Rule**: `ResponseDTO<T>` is an internal transport contract. It must **never** be serialized as the HTTP response body sent to external clients.

### Why this split?

- Internal layers need a uniform carrier for success/failure state across service calls.
- External clients expect standard HTTP semantics: status codes, raw resource bodies, and RFC 9457 `ProblemDetails` for errors — not custom envelopes.

---

## Critical Patterns

### ResponseDTO\<T\> — Internal Contract Only

Used exclusively between Repository and Service layers:

```csharp
// Shared/Models/ResponseDTO.cs
public class ResponseDTO<T>
{
    public bool Success { get; set; }
    public string Message { get; set; } = string.Empty;
    public T? Data { get; set; }
}
```

### Three Rules for ResponseDTO (Internal Use)

| Scenario | Success | Message | Data |
|---|---|---|---|
| Operation succeeded with data | `true` | `""` or informative | The result object/list |
| Operation succeeded, no results | `true` | Informative (e.g., "No records found.") | Empty value, **never null** |
| Operation failed / error | `false` | Clear error description | `null` |

**Empty value by type** (when success but no data):

| Type `T` | Data value |
|---|---|
| `List<T>` / `IEnumerable<T>` | `new List<T>()` / `[]` |
| `string` | `string.Empty` |
| `int`, `double`, etc. | `0` |
| Complex object | `new T()` (if parameterless constructor exists) |

```json
// Internal success with data (ResponseDTO — stays inside the API, never sent to client)
{ "success": true, "message": "", "data": { "id": "...", "name": "Alice" } }

// Internal success, no data
{ "success": true, "message": "No users found.", "data": [] }

// Internal error indicator
{ "success": false, "message": "User not found.", "data": null }
```

### Layer Communication Flow

```
Repository  →      Service      →   Controller  →   Client
  Entity     ResponseDTO<DTO>     IActionResult     HTTP 2xx (raw JSON)
                                                    HTTP 4xx/5xx (ProblemDetails)
```

| Layer | Returns | Rule |
|---|---|---|
| Repository | Domain entities or `null` | Never return DTOs or `ResponseDTO` from repositories |
| Service | `ResponseDTO<DTO>` (or throws domain exceptions) | Never return raw entities to the controller |
| Controller | `IActionResult` — raw resource JSON (2xx) or `ProblemDetails` (4xx/5xx) | Extract `.Data` from `ResponseDTO<T>`; do NOT forward the envelope to the client |

---

### Repository — Return Patterns

```csharp
// Returns entity or null — never throw for "not found" at this layer
public async Task<User?> GetByIdAsync(Guid id) =>
    await _context.Users.AsNoTracking().FirstOrDefaultAsync(u => u.Id == id);

// Returns list — never return null, return empty list
public async Task<List<User>> GetAllActiveAsync() =>
    await _context.Users.Where(u => u.IsActive).AsNoTracking().ToListAsync();
```

### Service — Return Patterns

The service is responsible for:
1. Calling the repository
2. Mapping entities to DTOs
3. Throwing domain exceptions for business errors

```csharp
public async Task<UserDto> GetByIdAsync(Guid id)
{
    var user = await _repository.GetByIdAsync(id);
    if (user is null)
        throw new NotFoundException(nameof(User), id);
    return user.ToDto();
}

public async Task<List<UserDto>> GetAllActiveAsync()
{
    var users = await _repository.GetAllActiveAsync();
    return users.Count == 0
        ? new List<UserDto>()
        : users.Select(u => u.ToDto()).ToList();
}

public async Task<UserDto> CreateAsync(CreateUserRequest request)
{
    var entity = request.ToEntity();
    await _repository.AddAsync(entity);
    return entity.ToDto();
}
```

### Controller — HTTP Status Code Mapping (External Contract)

Controllers extract `.Data` from the service result and return it directly as the HTTP response body. **No `ResponseDTO<T>` envelope is sent to the client.**

```csharp
[HttpGet("{id:guid}")]
public async Task<IActionResult> GetById([FromRoute] Guid id)
{
    var data = await _userService.GetByIdAsync(id);
    return Ok(data); // HTTP 200 — raw UserDto, no envelope
}

[HttpPost]
public async Task<IActionResult> Create([FromBody] CreateUserRequest request)
{
    var data = await _userService.CreateAsync(request);
    return CreatedAtAction(nameof(GetById), new { id = data.Id }, data);
    // HTTP 201 — raw UserDto, Location header included
}

[HttpPut("{id:guid}")]
public async Task<IActionResult> Update([FromRoute] Guid id, [FromBody] UpdateUserRequest request)
{
    var data = await _userService.UpdateAsync(id, request);
    return Ok(data); // HTTP 200 — raw UserDto
}

[HttpDelete("{id:guid}")]
public async Task<IActionResult> Delete([FromRoute] Guid id)
{
    await _userService.DeleteAsync(id);
    return NoContent(); // HTTP 204 — no body
}
```

> **Error responses** from controllers are handled by the global exception middleware — see [Exception Handling — HTTP Mapping](#exception-handling--http-mapping) below.

---

### HTTP Status Code Reference

| Code | Meaning | When to use |
|---|---|---|
| `200 OK` | Success | GET, PUT (when returning resource) |
| `201 Created` | Resource created | POST — include `Location` header |
| `204 No Content` | Success, no body | DELETE, PUT (no resource returned) |
| `400 Bad Request` | Invalid input | Validation errors, malformed request |
| `401 Unauthorized` | Not authenticated | **Resolved by pipeline/auth — NOT a domain exception** |
| `403 Forbidden` | Not authorized | Authenticated but no permission |
| `404 Not Found` | Resource absent | ID not found |
| `409 Conflict` | State conflict | Duplicate resource, concurrent modification |
| `429 Too Many Requests` | Rate limited | Rate limiter triggered |
| `500 Internal Server Error` | Unhandled server error | Global exception middleware |

---

## Exception Handling — HTTP Mapping

This section defines the **canonical mapping** of domain exceptions to HTTP status codes and `ProblemDetails` response shape. This is the **single source of truth** for the entire skill ecosystem.

> **Ownership**: The exception → HTTP mapping contract lives here (`responses`).  
> The actual middleware implementation, logging, and correlation ID tracking live in [`logging`](../logging/SKILL.md).

### Canonical Exception → Status Code Table

| Exception | HTTP Status | Log Level | Client Message |
|---|---|---|---|
| `NotFoundException` | `404 Not Found` | Warning | `ex.Message` (safe — contains resource + id) |
| `ForbiddenException` | `403 Forbidden` | Warning | `ex.Message` (safe — "You do not have access…") |
| `ConflictException` | `409 Conflict` | Warning | `ex.Message` (safe — "Email already registered") |
| `ValidationException` (manual) | `400 Bad Request` | Information | First error message |
| `Exception` (unhandled) | `500 Internal Server Error` | Error | Generic — **never expose stack trace** |

> **`401 Unauthorized` — Pipeline/Auth ownership (D-27)**  
> `401` is **not a domain exception**. It is resolved entirely by the authentication pipeline (JWT bearer middleware, `[Authorize]` attributes, or a dedicated auth middleware).  
> **Do NOT create an `UnauthorizedException` domain class.** If you need to signal authentication failure inside a service, propagate it as an HTTP concern at the pipeline level, not as a business exception.

### Exception Location Reference

| Exception | Defined in | Constructor |
|---|---|---|
| `NotFoundException` | `general` skill — `Shared/Exceptions/NotFoundException.cs` | `NotFoundException(string entity, object id)` |
| `ForbiddenException` | `security` skill — `Shared/Exceptions/ForbiddenException.cs` | `ForbiddenException(string message)` |
| `ConflictException` | `Shared/Exceptions/ConflictException.cs` | `ConflictException(string message)` |

> See [`general`](../general/SKILL.md) and [`security`](../security/SKILL.md) for exception definitions. Never duplicate them here.

### Middleware Error Response — ProblemDetails Shape (D-26)

All error responses sent to clients follow [RFC 9457 `ProblemDetails`](https://datatracker.ietf.org/doc/html/rfc9457). The exception handler middleware formats responses as:

```csharp
// Catch block example inside ExceptionHandlingMiddleware.cs
// Full implementation and logging logic: see logging/SKILL.md
catch (NotFoundException ex)
{
    context.Response.StatusCode = StatusCodes.Status404NotFound;
    await context.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status404NotFound,
        Title = "Not Found",
        Detail = ex.Message
    });
}
catch (ForbiddenException ex)
{
    context.Response.StatusCode = StatusCodes.Status403Forbidden;
    await context.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status403Forbidden,
        Title = "Forbidden",
        Detail = ex.Message
    });
}
catch (ConflictException ex)
{
    context.Response.StatusCode = StatusCodes.Status409Conflict;
    await context.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status409Conflict,
        Title = "Conflict",
        Detail = ex.Message
    });
}
catch (Exception)
{
    context.Response.StatusCode = StatusCodes.Status500InternalServerError;
    await context.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status500InternalServerError,
        Title = "Internal Server Error",
        Detail = "An unexpected error occurred."
        // NEVER expose stack trace or internal details
    });
}
```

> For logging, correlation IDs, and the full middleware implementation, see [`logging`](../logging/SKILL.md).

---

## Uniform 400 — Auto-Validation with ProblemDetails (D-26)

When FluentValidation or model binding rejects a request automatically (before the controller action runs), `InvalidModelStateResponseFactory` intercepts and formats the response as `ProblemDetails`:

```csharp
// Program.cs — inside AddControllers()
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

> For validation rules and FluentValidation setup, see [`validations`](../validations/SKILL.md).

---

## Uniform 429 — Rate Limiter with ProblemDetails (D-26)

When the rate limiter rejects a request, the `OnRejected` callback formats the response as `ProblemDetails`:

```csharp
// Program.cs — inside AddRateLimiter()
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

> For rate limiter configuration and policies, see [`security`](../security/SKILL.md).

---

### Paged Response — Standard Structure

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

Controller response (raw, no envelope):
```csharp
var pagedResult = await _userService.GetAllAsync(pagination);
return Ok(pagedResult); // HTTP 200 — raw PagedResult<UserDto>
```

---

### Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| Returning raw entity from controller | Exposes data model; security risk |
| Returning `null` from service on "not found" | Controller must null-check everywhere; throw instead |
| Mixing business errors with HTTP concerns in service | Couples service to HTTP layer |
| `ResponseDTO<T>` inside repositories | Repositories are data-layer, not API-layer |
| Wrapping `ResponseDTO<T>` in the HTTP response body sent to the client | Violates the dual contract; clients should receive raw JSON or `ProblemDetails` |
| HTTP 200 for errors | Misleads clients; always use correct status codes |
| No `Location` header on 201 Created | Violates REST convention |
| Creating `UnauthorizedException` as a domain exception | `401` is owned by the pipeline/auth — not the business domain |

---

## Code Examples

### Full CRUD — Controller Pattern

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

    [HttpGet]
    public async Task<IActionResult> GetAll([FromQuery] PaginationRequest pagination)
    {
        var data = await _userService.GetAllAsync(pagination);
        return Ok(data); // HTTP 200 — raw PagedResult<UserDto>
    }

    [HttpGet("{id:guid}")]
    public async Task<IActionResult> GetById([FromRoute] Guid id)
    {
        var data = await _userService.GetByIdAsync(id);
        return Ok(data); // HTTP 200 — raw UserDto
    }

    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreateUserRequest request)
    {
        var data = await _userService.CreateAsync(request);
        return CreatedAtAction(nameof(GetById), new { id = data.Id }, data);
        // HTTP 201 — raw UserDto, Location header included
    }

    [HttpPut("{id:guid}")]
    public async Task<IActionResult> Update([FromRoute] Guid id, [FromBody] UpdateUserRequest request)
    {
        var data = await _userService.UpdateAsync(id, request);
        return Ok(data); // HTTP 200 — raw UserDto
    }

    [HttpDelete("{id:guid}")]
    public async Task<IActionResult> Delete([FromRoute] Guid id)
    {
        await _userService.DeleteAsync(id);
        return NoContent(); // HTTP 204 — no body
    }
}
```

> Errors (404, 403, 409, 500) are handled automatically by `ExceptionHandlingMiddleware` — no try/catch needed in controllers. See [Exception Handling — HTTP Mapping](#exception-handling--http-mapping).

---

## Resources

- **Standards**: See [../../rules-to-skills/Standardized_NET_Rules.md](../../rules-to-skills/Standardized_NET_Rules.md)
- **General / NotFoundException**: See [../general/SKILL.md](../general/SKILL.md)
- **Requests**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md)
- **Logging / ExceptionHandlingMiddleware**: See [../logging/SKILL.md](../logging/SKILL.md)
- **Validations / FluentValidation**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Security / ForbiddenException / Rate Limiting**: See [../security/SKILL.md](../security/SKILL.md)

---

## Changelog

### v1.1 — 2026-03-24

**Dual Response Contract (D-25)**
- Clarified that `ResponseDTO<T>` is an **internal** transport contract only (Repository → Service).
- Controllers now return **raw resource JSON** (2xx) or **`ProblemDetails`** (4xx/5xx) to external clients.
- Updated all controller snippets and the Full CRUD example to remove `ResponseDTO<T>` from HTTP response bodies.
- Added explicit anti-pattern: "Wrapping `ResponseDTO<T>` in the HTTP response body sent to the client."
- Updated `Layer Communication Flow` table to reflect the correct external boundary.

**ProblemDetails for External Errors (D-26)**
- Added new section **"Exception Handling — HTTP Mapping"** with the canonical exception → HTTP status code table.
- Middleware catch-block examples now produce `ProblemDetails` (RFC 9457), not `ResponseDTO<T>`.
- Uniform `400 Bad Request`: `InvalidModelStateResponseFactory` now returns `ProblemDetails`.
- Uniform `429 Too Many Requests`: `OnRejected` callback now returns `ProblemDetails`.

**`401` belongs to the pipeline/auth — not the domain (D-27)**
- Added explicit note in the HTTP Status Code table: `401 Unauthorized` is resolved by the auth pipeline.
- Added canonical table note: **Do NOT create `UnauthorizedException` as a domain exception.**
- Added anti-pattern: "Creating `UnauthorizedException` as a domain exception."

**Exception Ownership**
- Added `ConflictException` (409) to the canonical table and exception location reference.
- Cross-referenced `NotFoundException` to `general` skill and `ForbiddenException` to `security` skill.

**Cross-References**
- Added `general` and `security` to the Resources section.
- Reinforced bidirectional cross-ref: `logging` owns middleware implementation; `responses` owns the mapping contract.

### v1.0 — Initial release

- `ResponseDTO<T>` unified contract across all layers.
- HTTP Status Code Reference table.
- Repository, Service, and Controller return patterns.
- Paged response structure.
- Anti-patterns table.
