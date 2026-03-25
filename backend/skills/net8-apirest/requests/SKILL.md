---
name: net8-apirest-requests
description: >
  Best practices for handling HTTP requests in ASP.NET Core 8 REST APIs covering model binding, DTO design, input binding sources, and request pipeline conventions.
  Trigger: When designing request DTOs, configuring model binding, handling query parameters, or processing incoming HTTP payloads in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.1"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Designing request DTOs for API endpoints
- Configuring how parameters are bound (body, query, route, header)
- Handling file uploads or multipart requests
- Applying consistent conventions for all incoming request shapes
- Reviewing controller action signatures

---

## Critical Patterns

### Request DTO Naming Convention

Use suffix `Request` for input DTOs:

```
CreateUserRequest
UpdateUserRequest
SearchUsersRequest
```

Never reuse domain entities as request bodies. Always use dedicated DTOs.

### Binding Sources — Explicit Annotation

Always annotate binding sources explicitly. Do not rely on implicit binding inference:

| Source | Attribute | Use case |
|---|---|---|
| JSON body | `[FromBody]` | POST / PUT payloads |
| Route segment | `[FromRoute]` | `/users/{id}` |
| Query string | `[FromQuery]` | Filters, pagination |
| Header | `[FromHeader]` | Custom headers (e.g., `X-Tenant-ID`) |
| Form | `[FromForm]` | File uploads, multipart |

```csharp
[HttpGet("{id:guid}")]
public async Task<IActionResult> GetById(
    [FromRoute] Guid id,
    [FromHeader(Name = "X-Tenant-ID")] string? tenantId)
{ ... }

[HttpGet]
public async Task<IActionResult> Search([FromQuery] SearchUsersRequest request) { ... }

[HttpPost]
public async Task<IActionResult> Create([FromBody] CreateUserRequest request) { ... }
```

### Route Constraints — Enforce Types at Routing Level

Use route constraints to reject invalid IDs before they reach the service:

```csharp
[HttpGet("{id:guid}")]     // only valid GUIDs
[HttpGet("{id:int:min(1)")] // only positive integers (avoid for public APIs — use GUIDs)
```

> **Rule**: Public resource endpoints must use `{id:guid}` to prevent BOLA enumeration attacks.

### Pagination — Standard Query Parameters

All list/search endpoints that may return large datasets must support pagination:

```csharp
public class PaginationRequest
{
    public int Page { get; set; } = 1;
    public int PageSize { get; set; } = 20;

    // Prevent abuse: cap max page size
    public int ValidatedPageSize => Math.Min(PageSize, 100);
}

public class SearchUsersRequest : PaginationRequest
{
    public string? Name { get; set; }
    public bool? IsActive { get; set; }
    public DateTime? RegisteredAfter { get; set; }
}
```

### Request Validation Flow

Requests must be validated **before** reaching the service layer. The controller must not call the service with an invalid request:

```
HTTP Request
    → Model Binding         (ASP.NET Core)
    → FluentValidation      (auto via filter or manual)
    → Controller Action     (only if valid)
    → Service Layer
```

See [../validations/SKILL.md](../validations/SKILL.md) for validation implementation details.

### Null Safety in Request DTOs

Use nullable types only when the field is truly optional. Required fields must be non-nullable:

```csharp
public class CreateUserRequest
{
    public string FirstName { get; set; } = string.Empty;   // required
    public string LastName { get; set; } = string.Empty;    // required
    public string Email { get; set; } = string.Empty;       // required
    public string? PhoneNumber { get; set; }                // optional
    public DateTime? BirthDate { get; set; }                // optional
}
```

### Content Negotiation

APIs must accept and return `application/json` by default. Enforce this via `[Produces]` and `[Consumes]`:

```csharp
[ApiController]
[Route("api/[controller]")]
[Produces("application/json")]
[Consumes("application/json")]
public class UsersController : ControllerBase { ... }
```

### File Uploads — Hybrid Validation Pattern (D-30)

File upload endpoints use a **hybrid validation approach**: FluentValidation handles metadata checks
(presence, size, MIME type), while the controller applies a minimal defensive guard for stream integrity.
This preserves SRP — the validator never reads the stream.

#### Step 1 — Request DTO

`IFormFile` must be initialized with `default!` (required for model binding in .NET 8).
The `.NotNull()` rule in FluentValidation provides the runtime guarantee:

```csharp
public class UploadAvatarRequest
{
    // Use default! + validator .NotNull() — do NOT use `required` keyword.
    // `required` does not integrate cleanly with [FromForm] model binding in .NET 8.
    public IFormFile File { get; set; } = default!;
}
```

#### Step 2 — FluentValidation: metadata only (no stream reading)

```csharp
public class UploadAvatarRequestValidator : AbstractValidator<UploadAvatarRequest>
{
    private static readonly string[] AllowedMimeTypes = ["image/jpeg", "image/png"];
    private const long MaxFileSize = 5 * 1024 * 1024; // 5 MB

    public UploadAvatarRequestValidator()
    {
        RuleFor(x => x.File)
            .NotNull().WithMessage("File is required.");

        // .Length and .ContentType are safe to access in validators — they are metadata,
        // NOT stream reads. Never call .OpenReadStream() inside a validator.
        RuleFor(x => x.File.Length)
            .LessThanOrEqualTo(MaxFileSize).WithMessage("File exceeds 5 MB limit.")
            .When(x => x.File is not null);

        RuleFor(x => x.File.ContentType)
            .Must(ct => AllowedMimeTypes.Contains(ct))
            .WithMessage("Invalid file type. Allowed: JPEG, PNG.")
            .When(x => x.File is not null);
    }
}
```

#### Step 3 — Controller: defensive stream guard only

The controller action receives the request **only if FluentValidation passed**. Its sole responsibility
here is to guard against a corrupted/empty stream — it does NOT re-validate metadata:

```csharp
[HttpPost("avatar")]
// [Consumes("multipart/form-data")] at action level overrides the controller-level
// [Consumes("application/json")]. This is required — do not remove it.
[Consumes("multipart/form-data")]
public async Task<IActionResult> UploadAvatar([FromForm] UploadAvatarRequest request)
{
    // Defensive guard: stream integrity only (NOT metadata — that is the validator's job).
    if (request.File.Length == 0)
        return BadRequest(new ProblemDetails
        {
            Status = StatusCodes.Status400BadRequest,
            Title = "Bad Request",
            Detail = "File stream is empty or corrupted."
        });

    await _avatarService.UploadAsync(request.File);
    return NoContent(); // or CreatedAtAction(...) depending on resource semantics
}
```

> **Content negotiation note**: `[Consumes("multipart/form-data")]` on an action **overrides** the
> controller-level `[Consumes("application/json")]`. Always annotate file upload actions explicitly —
> omitting it causes ASP.NET Core to reject `multipart/form-data` requests with 415 Unsupported Media Type.

#### Data flow summary

```
HTTP multipart/form-data
  → Model Binding ([FromForm] UploadAvatarRequest)
  → FluentValidation: UploadAvatarRequestValidator
      ├── .NotNull()                              — File present
      ├── .Length <= MaxFileSize                  — size metadata
      └── .ContentType in AllowedMimeTypes        — MIME metadata
  → Controller Action  (only if all rules pass)
      └── if (file.Length == 0) → ProblemDetails 400   — stream guard
  → Service Layer (storage, processing)
```

### Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| Domain entity as request body | Couples API to data model; security risk |
| No binding source annotation | Implicit binding can cause unexpected behavior |
| Unbounded pagination | Allows `pageSize=9999`, causes performance issues |
| int IDs on public endpoints | Enables BOLA / resource enumeration |
| Validation inside the service | Mixes concerns; validation belongs at the boundary |
| Generic `object` as request body | No type safety, no validation |
| `ResponseDTO<T>` as HTTP response body | Exposes internal application wrapper to HTTP clients; violates D-25/D-26/D-27. Return raw resources or `ProblemDetails` — never envelope at HTTP boundary |
| Reading `IFormFile` stream in FluentValidation | Blocks the stream before the controller reads it; breaks SRP. Validators must only access metadata (`.Length`, `.ContentType`) — never call `.OpenReadStream()` |

---

## Code Examples

### Paginated Search Endpoint

```csharp
[HttpGet]
public async Task<IActionResult> Search([FromQuery] SearchUsersRequest request)
{
    var result = await _userService.SearchAsync(request);
    return Ok(result); // D-25: raw PagedResult<UserDto>, no envelope
}
```

### Accepting Optional Tenant Context via Header

```csharp
[HttpGet]
public async Task<IActionResult> GetAll(
    [FromQuery] PaginationRequest pagination,
    [FromHeader(Name = "X-Tenant-ID")] Guid? tenantId)
{
    var result = await _userService.GetAllAsync(pagination, tenantId);
    return Ok(result);
}
```

---

## Resources

- **Standards**: See [../../../../rules-to-skills/Standardized_NET_Rules.md](../../../../rules-to-skills/Standardized_NET_Rules.md)
- **General structure**: See [../general/SKILL.md](../general/SKILL.md) — `Features/{Feature}/DTOs/` is where Request DTOs live
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md) — request-to-entity mapping patterns
- **Validations**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Responses**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Security**: See [../security/SKILL.md](../security/SKILL.md) — OWASP API4 file upload considerations, request body size limits

---

## Changelog

### v1.1 (2026-03-24)

- **File Uploads rewritten**: Replaced manual metadata validation in controller with the hybrid D-30 pattern.
  FluentValidation handles all metadata (presence, size, MIME); controller applies minimal stream-integrity guard only.
- **Paginated Search corrected**: `Search` example now returns `Ok(result)` directly — removed `ResponseDTO<PagedResult<UserDto>>` wrapper (D-25).
- **Anti-patterns expanded**: Added `ResponseDTO<T>` as HTTP response body and stream reading in FluentValidation.
- **IFormFile nullable safety documented**: Clarified `default!` + `.NotNull()` pattern and why `required` keyword is not used.
- **Content negotiation note added**: `[Consumes("multipart/form-data")]` at action level overrides controller-level attribute.
- **Cross-references added**: `general` (DTO folder structure) and `mapping` (request-to-entity) added to Resources.

### v1.0 (initial)

- Initial release: binding sources, pagination, null safety, content negotiation, basic file upload, anti-patterns table.
