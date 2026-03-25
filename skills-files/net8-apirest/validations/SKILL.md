---
name: net8-apirest-validations
description: >
  Input validation standards for ASP.NET Core 8 REST APIs using FluentValidation, covering validator structure, ASP.NET Core integration, auto-registration, and decision rules between FluentValidation and manual validation.
  Trigger: When adding validation to request DTOs, creating FluentValidation validators, or deciding how to validate incoming data in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.1"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Adding validation rules to any request DTO
- Creating a new FluentValidation validator class
- Deciding between FluentValidation and manual validation
- Handling validation errors in controllers
- Configuring automatic validation via filters

---

## Critical Patterns

### Decision Tree — FluentValidation vs Manual

```
Is the rule purely format/structural? (length, required, email, regex)
  YES → FluentValidation

Does the rule require a database lookup? (unique email, existing user)
  YES → FluentValidation with async custom rule (.MustAsync)

Is the rule complex business logic internal to the service?
  YES → throw domain exception inside the service, not a validator

Is validation part of an external system call?
  YES → validate in service, throw specific exception
```

> **Default rule**: All boundary validation (HTTP input) uses FluentValidation. All business-rule violations use domain exceptions in the service layer.

### Validator File Location

```
Features/{Feature}/
  Validators/
    CreateUserRequestValidator.cs
    UpdateUserRequestValidator.cs
```

One validator class per DTO. Never share validators across DTOs.

### FluentValidation — Validator Structure

```csharp
// Features/Users/Validators/CreateUserRequestValidator.cs
public class CreateUserRequestValidator : AbstractValidator<CreateUserRequest>
{
    public CreateUserRequestValidator()
    {
        RuleFor(x => x.FirstName)
            .NotEmpty().WithMessage("First name is required.")
            .MaximumLength(100).WithMessage("First name must not exceed 100 characters.");

        RuleFor(x => x.LastName)
            .NotEmpty().WithMessage("Last name is required.")
            .MaximumLength(100).WithMessage("Last name must not exceed 100 characters.");

        RuleFor(x => x.Email)
            .NotEmpty().WithMessage("Email is required.")
            .EmailAddress().WithMessage("Email format is invalid.")
            .MaximumLength(200).WithMessage("Email must not exceed 200 characters.");

        RuleFor(x => x.PhoneNumber)
            .Matches(@"^\+?[1-9]\d{7,14}$").WithMessage("Phone number format is invalid.")
            .When(x => !string.IsNullOrEmpty(x.PhoneNumber));

        RuleFor(x => x.Password)
            .NotEmpty().WithMessage("Password is required.")
            .MinimumLength(8).WithMessage("Password must be at least 8 characters.")
            .Matches("[A-Z]").WithMessage("Password must contain at least one uppercase letter.")
            .Matches("[0-9]").WithMessage("Password must contain at least one number.");
    }
}
```

### Auto-Registration — All Validators from Assembly

Register all validators at once in `Program.cs`. Never register them individually:

```csharp
builder.Services.AddValidatorsFromAssemblyContaining<Program>();
```

### Integration Mode — Automatic Validation Filter (Recommended)

Use the FluentValidation ASP.NET Core integration package to validate automatically via the model binding pipeline:

```csharp
// Program.cs
builder.Services.AddFluentValidationAutoValidation();
builder.Services.AddValidatorsFromAssemblyContaining<Program>();
```

With auto-validation enabled, the action is only reached if the model is valid. The framework returns `400 Bad Request` automatically with a **`ProblemDetails`** body (configured via `InvalidModelStateResponseFactory` — see [`responses/SKILL.md`](../responses/SKILL.md) for the canonical setup).

> **Rule**: Prefer **manual validation** (inject `IValidator<T>`) when you need programmatic access to the validation result: logging errors before returning, multi-step flows that conditionally branch on specific failures, or service-layer logic that evaluates errors before deciding the response.

### Manual Validation — Programmatic Access to Result

Use manual validation when the controller or service needs to inspect, log, or conditionally react to specific validation errors — not merely to control the response shape:

```csharp
[HttpPost]
public async Task<IActionResult> Create([FromBody] CreateUserRequest request)
{
    var validation = await _validator.ValidateAsync(request);
    if (!validation.IsValid)
    {
        // Log specific errors, evaluate conditions, or build a structured detail
        return BadRequest(new ProblemDetails
        {
            Status = StatusCodes.Status400BadRequest,
            Title = "Bad Request",
            Detail = validation.Errors.First().ErrorMessage
        });
    }
    var data = await _userService.CreateAsync(request);
    return CreatedAtAction(nameof(GetById), new { id = data.Id },
        new ResponseDTO<UserDto> { Success = true, Data = data });
}
```

Inject the validator via constructor:
```csharp
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;
    private readonly IValidator<CreateUserRequest> _validator;

    public UsersController(IUserService userService, IValidator<CreateUserRequest> validator)
    {
        _userService = userService;
        _validator = validator;
    }
}
```

> **Note**: MediatR Pipeline Behaviors are an **optional extension** for cross-cutting validation in command/query pipelines. They are not the canonical approach for this skill — FluentValidation with auto-validation filter or direct `IValidator<T>` injection covers the boundary contract without coupling to a mediator.

### Async Validation — Database Rules

Use `.MustAsync` for rules requiring async calls (uniqueness checks, existence checks):

```csharp
public class CreateUserRequestValidator : AbstractValidator<CreateUserRequest>
{
    private readonly IUserRepository _repository;

    public CreateUserRequestValidator(IUserRepository repository)
    {
        _repository = repository;

        RuleFor(x => x.Email)
            .NotEmpty()
            .EmailAddress()
            .MustAsync(BeUniqueEmailAsync)
            .WithMessage("Email is already registered.");
    }

    private async Task<bool> BeUniqueEmailAsync(string email, CancellationToken ct) =>
        !await _repository.ExistsByEmailAsync(email, ct);
}
```

### Validation Error Response Format

All HTTP 400 responses triggered by validation failures **must use `ProblemDetails`** (RFC 7807). There is a single canonical format for external consumers:

```json
{
  "type": "https://tools.ietf.org/html/rfc7807",
  "title": "Bad Request",
  "status": 400,
  "detail": "Email format is invalid."
}
```

This applies to both **automatic** (via `InvalidModelStateResponseFactory`) and **manual** validation paths. Refer to [`responses/SKILL.md`](../responses/SKILL.md) for the `InvalidModelStateResponseFactory` configuration that produces this format automatically from the model state.

> `ResponseDTO<T>` is an **internal contract** (service/application layer → controller). Controllers **always** map a failed `ResponseDTO<T>` to a `ProblemDetails` HTTP 400 before returning to the client — it is never exposed directly as an HTTP response body.

### Common Validation Rules Reference

```csharp
// Required
RuleFor(x => x.Name).NotEmpty();

// Length
RuleFor(x => x.Name).Length(2, 100);

// Email
RuleFor(x => x.Email).EmailAddress();

// Regex
RuleFor(x => x.Code).Matches(@"^[A-Z]{3}\d{3}$");

// Numeric range
RuleFor(x => x.Age).InclusiveBetween(18, 120);

// Conditional
RuleFor(x => x.CompanyName)
    .NotEmpty()
    .When(x => x.AccountType == AccountType.Business);

// Collection
RuleForEach(x => x.Tags)
    .NotEmpty()
    .MaximumLength(50);

// Custom
RuleFor(x => x.StartDate)
    .Must(date => date > DateTime.UtcNow)
    .WithMessage("Start date must be in the future.");
```

### Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| `DataAnnotations` on DTOs | Less expressive, harder to test, no async support |
| Business logic in validators | Mixes concerns; use service exceptions instead |
| Validation inside service for format rules | Too late in the pipeline; validate at boundary |
| One global validator for multiple DTOs | Breaks SRP; one validator per DTO |
| Swallowing validation errors | Clients cannot fix requests without error details |
| `if (model == null) return BadRequest()` alone | Missing field-level error messages |
| Exposing `ResponseDTO<T>` as HTTP 400 body | Violates dual-contract rule (D-25); external errors always use `ProblemDetails` |

---

## Code Examples

### Full Validator with Multiple Rule Sets

```csharp
public class UpdateUserRequestValidator : AbstractValidator<UpdateUserRequest>
{
    public UpdateUserRequestValidator()
    {
        RuleFor(x => x.FirstName)
            .NotEmpty()
            .MaximumLength(100)
            .When(x => x.FirstName is not null);

        RuleFor(x => x.LastName)
            .NotEmpty()
            .MaximumLength(100)
            .When(x => x.LastName is not null);

        RuleFor(x => x.PhoneNumber)
            .Matches(@"^\+?[1-9]\d{7,14}$")
            .When(x => !string.IsNullOrEmpty(x.PhoneNumber));
    }
}
```

---

## Commands

```bash
# Add FluentValidation packages
dotnet add package FluentValidation
dotnet add package FluentValidation.AspNetCore
dotnet add package FluentValidation.DependencyInjectionExtensions
```

---

## Changelog

### v1.1.0
- **Contrato externo 400**: Todos los ejemplos de respuesta HTTP 400 por validación ahora usan `ProblemDetails` (RFC 7807) en lugar de `ResponseDTO<T>` — alineación con decisiones D-25/D-26 de `responses` v1.1.
- **Estrategia canónica agnóstica a MediatR**: FluentValidation con auto-validation filter o `IValidator<T>` directo es el enfoque canónico (D-28). MediatR Pipeline Behaviors son extensión opcional documentada explícitamente.
- **Motivación de validación manual**: Actualizada para reflejar casos de uso reales — acceso programático al resultado, logging, flujos multi-paso y lógica condicional sobre errores específicos (D-30).
- **Anti-pattern nuevo**: Agregado "Exponer `ResponseDTO<T>` como HTTP 400 body" a la tabla de anti-patterns (D-32).
- **Cross-references**: Agregadas referencias explícitas a `responses/SKILL.md` (contrato HTTP y `InvalidModelStateResponseFactory`), `testing-unit/SKILL.md` (testing de validators) y `general/SKILL.md` (estructura de carpetas `Exceptions/`).

---

## Resources

- **FluentValidation docs**: https://docs.fluentvalidation.net/en/latest/aspnet.html
- **Standards**: See [../../rules-to-skills/Standardized_NET_Rules.md](../../rules-to-skills/Standardized_NET_Rules.md)
- **Requests**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Responses (HTTP contract + InvalidModelStateResponseFactory)**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Testing validators**: See [../testing-unit/SKILL.md](../testing-unit/SKILL.md)
- **Architecture (Exceptions/ folder structure)**: See [../general/SKILL.md](../general/SKILL.md)
