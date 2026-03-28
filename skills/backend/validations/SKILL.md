---
name: net8-apirest-validations
description: >
  Input validation standards for ASP.NET Core 8 REST APIs using FluentValidation, covering validator structure, ASP.NET Core integration, auto-registration, and decision rules between FluentValidation and manual validation.
  Trigger: When adding validation to request DTOs, creating FluentValidation validators, or deciding how to validate incoming data in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.3"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Adding validation rules to any request DTO
- Deciding between FluentValidation vs manual validation
- Configuring automatic validation via filters

---

## Critical Patterns

### Where each type of validation belongs

```
Is the rule format/structural? (length, required, email, regex)
  → FluentValidation

Does the rule require a DB lookup? (unique email, existing user)
  → FluentValidation with .MustAsync

Is it business logic internal to the service?
  → throw domain exception in the service, not in the validator

External system validation?
  → validate in service, throw specific exception
```

> **Base rule**: All boundary validation (HTTP input) uses FluentValidation. All business rule violations use domain exceptions in the service layer.

### Setup — Auto-registration + Auto-validation

> **FluentValidation v11 breaking change**: `AddFluentValidationAutoValidation()` and `AddFluentValidationClientsideAdapters()` were removed from the core `FluentValidation` package. They now live exclusively in `FluentValidation.AspNetCore`. Install both packages:

```bash
dotnet add package FluentValidation.AspNetCore   # contains AddFluentValidationAutoValidation
dotnet add package FluentValidation.DependencyInjectionExtensions  # contains AddValidatorsFromAssemblyContaining
```

```csharp
// Program.cs — a single line registers all validators in the assembly
builder.Services.AddFluentValidationAutoValidation();
builder.Services.AddValidatorsFromAssemblyContaining<Program>();
```

With auto-validation active, the action only executes if the model is valid. The framework returns `400 Bad Request` with `ProblemDetails` automatically (configured via `InvalidModelStateResponseFactory` — see [`responses/SKILL.md`](../responses/SKILL.md)).

### Manual Validation — Only when you need the result programmatically

Use injected `IValidator<T>` when the controller or service needs to inspect, log, or conditionally react to specific errors — not just to control the shape of the response:

```csharp
var validation = await _validator.ValidateAsync(request);
if (!validation.IsValid)
    return BadRequest(new ProblemDetails
    {
        Status = StatusCodes.Status400BadRequest,
        Title = "Bad Request",
        Detail = validation.Errors.First().ErrorMessage
    });
```

### File Upload — Hybrid Validation Pattern (D-30)

FluentValidation handles metadata (presence, size, MIME type). The controller applies only a stream integrity guard. **The validator never reads the stream.**

```csharp
// Request DTO — default! + .NotNull() in validator (do not use `required` with [FromForm])
public class UploadAvatarRequest
{
    public IFormFile File { get; set; } = default!;
}

// Validator — metadata only, never .OpenReadStream()
public class UploadAvatarRequestValidator : AbstractValidator<UploadAvatarRequest>
{
    private static readonly string[] AllowedMimeTypes = ["image/jpeg", "image/png"];
    private const long MaxFileSize = 5 * 1024 * 1024;

    public UploadAvatarRequestValidator()
    {
        RuleFor(x => x.File).NotNull().WithMessage("File is required.");
        RuleFor(x => x.File.Length)
            .LessThanOrEqualTo(MaxFileSize).WithMessage("File exceeds 5 MB limit.")
            .When(x => x.File is not null);
        RuleFor(x => x.File.ContentType)
            .Must(ct => AllowedMimeTypes.Contains(ct))
            .WithMessage("Invalid file type. Allowed: JPEG, PNG.")
            .When(x => x.File is not null);
    }
}

// Controller — stream integrity guard only
[HttpPost("avatar")]
[Consumes("multipart/form-data")]
public async Task<IActionResult> UploadAvatar([FromForm] UploadAvatarRequest request)
{
    if (request.File.Length == 0)
        return BadRequest(new ProblemDetails
        {
            Status = StatusCodes.Status400BadRequest,
            Title = "Bad Request",
            Detail = "File stream is empty or corrupted."
        });

    await _avatarService.UploadAsync(request.File);
    return NoContent();
}
```

---

## Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| `DataAnnotations` on DTOs | Less expressive, no async support — always use FluentValidation |
| Business logic in validators | Mixes concerns; use domain exceptions in the service |
| Format validation in the service | Too late in the pipeline |
| A single global validator for multiple DTOs | Violates SRP — one validator per DTO |
| Reading `IFormFile` stream in FluentValidation | Blocks the stream before the controller reads it |
| Exposing `ResponseDTO<T>` as the body of HTTP 400 | Violates D-25/D-26 — always use `ProblemDetails` for HTTP errors |

---

## Commands

```bash
# FluentValidation v11+ — three separate packages required
dotnet add package FluentValidation                                  # core validators
dotnet add package FluentValidation.AspNetCore                       # AddFluentValidationAutoValidation (v11: split from core)
dotnet add package FluentValidation.DependencyInjectionExtensions    # AddValidatorsFromAssemblyContaining
```

---

## Resources

- **FluentValidation docs**: https://docs.fluentvalidation.net/en/latest/aspnet.html
- **Responses (HTTP contract + InvalidModelStateResponseFactory)**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Requests (file upload context)**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Testing validators**: See [../testing-unit/SKILL.md](../testing-unit/SKILL.md)

---

## Changelog

### v1.3 — 2026-03-28
- **Fixed (W-03)**: Added FluentValidation v11 breaking change note — `AddFluentValidationAutoValidation()` was removed from the core `FluentValidation` package and moved to `FluentValidation.AspNetCore`. Updated Commands section to list all three packages with explicit comments.

### v1.2 — 2026-03-28
- **Removed**: FluentValidation tutorial (rules, syntax, validator examples) — the agent already knows it
- **Removed**: Extensive validator structure section with full CreateUserRequestValidator example
- **Removed**: Common validation rules reference (`NotEmpty`, `EmailAddress`, `Matches`, etc.)
- **Removed**: Full validator example with multiple rule sets
- **Kept**: Decision tree (where each validation belongs), D-30 hybrid file upload, auto-registration setup, manual validation with criteria, anti-patterns

### v1.1 — 2026-03-24
- External 400 contract migrated to ProblemDetails; D-30 file upload pattern; expanded anti-patterns
