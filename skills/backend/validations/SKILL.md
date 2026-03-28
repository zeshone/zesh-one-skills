---
name: net8-apirest-validations
description: >
  Input validation standards for ASP.NET Core 8 REST APIs using FluentValidation, covering validator structure, ASP.NET Core integration, auto-registration, and decision rules between FluentValidation and manual validation.
  Trigger: When adding validation to request DTOs, creating FluentValidation validators, or deciding how to validate incoming data in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Agregando reglas de validación a cualquier request DTO
- Decidiendo entre FluentValidation vs validación manual
- Configurando validación automática vía filters

---

## Critical Patterns

### Dónde va cada tipo de validación

```
¿La regla es format/estructural? (length, required, email, regex)
  → FluentValidation

¿La regla requiere un DB lookup? (email único, usuario existente)
  → FluentValidation con .MustAsync

¿Es lógica de negocio interna al service?
  → throw domain exception en el service, no en el validator

¿Validación de sistema externo?
  → validar en service, throw specific exception
```

> **Regla base**: Toda validación de boundary (input HTTP) usa FluentValidation. Toda violación de regla de negocio usa domain exceptions en el service layer.

### Setup — Auto-registration + Auto-validation

```csharp
// Program.cs — una sola línea registra todos los validators del assembly
builder.Services.AddFluentValidationAutoValidation();
builder.Services.AddValidatorsFromAssemblyContaining<Program>();
```

Con auto-validation activo, la action solo se ejecuta si el modelo es válido. El framework retorna `400 Bad Request` con `ProblemDetails` automáticamente (configurado via `InvalidModelStateResponseFactory` — ver [`responses/SKILL.md`](../responses/SKILL.md)).

### Manual Validation — Solo cuando necesitás el resultado programáticamente

Usar `IValidator<T>` inyectado cuando el controller o service necesita inspeccionar, loguear o reaccionar condicionalmente a errores específicos — no solo para controlar el shape de la respuesta:

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

FluentValidation maneja metadata (presencia, tamaño, MIME type). El controller aplica solo un guard de integridad de stream. **El validator nunca lee el stream.**

```csharp
// Request DTO — default! + .NotNull() en validator (no usar `required` con [FromForm])
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

// Controller — solo stream integrity guard
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

| Anti-pattern | Problema |
|---|---|
| `DataAnnotations` en DTOs | Menos expresivo, sin async support — siempre FluentValidation |
| Business logic en validators | Mezcla concerns; usar domain exceptions en service |
| Validación de formato en el service | Demasiado tarde en el pipeline |
| Un validator global para múltiples DTOs | Viola SRP — un validator por DTO |
| Leer `IFormFile` stream en FluentValidation | Bloquea el stream antes de que el controller lo lea |
| Exponer `ResponseDTO<T>` como body de HTTP 400 | Viola D-25/D-26 — siempre `ProblemDetails` en errores HTTP |

---

## Commands

```bash
dotnet add package FluentValidation
dotnet add package FluentValidation.AspNetCore
dotnet add package FluentValidation.DependencyInjectionExtensions
```

---

## Resources

- **FluentValidation docs**: https://docs.fluentvalidation.net/en/latest/aspnet.html
- **Responses (HTTP contract + InvalidModelStateResponseFactory)**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Requests (file upload context)**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Testing validators**: See [../testing-unit/SKILL.md](../testing-unit/SKILL.md)

---

## Changelog

### v1.2 — 2026-03-28
- **Removed**: Tutorial de FluentValidation (reglas, sintaxis, ejemplos de validators) — el agente ya lo conoce
- **Removed**: Sección extensa de validator structure con ejemplo completo de CreateUserRequestValidator
- **Removed**: Common validation rules reference (`NotEmpty`, `EmailAddress`, `Matches`, etc.)
- **Removed**: Full validator example con múltiples rule sets
- **Kept**: Decision tree (dónde va cada validación), D-30 hybrid file upload, setup de auto-registration, manual validation con criterio, anti-patterns

### v1.1 — 2026-03-24
- Contrato externo 400 migrado a ProblemDetails; D-30 file upload pattern; anti-patterns expandidos
