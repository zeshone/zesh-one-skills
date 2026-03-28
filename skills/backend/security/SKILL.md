---
name: net8-apirest-security
description: >
  Security standards for ASP.NET Core 8 REST APIs based on OWASP API Security Top 10 2023, covering authorization, authentication, BOLA prevention, input validation, rate limiting, and secure credential handling.
  Trigger: When implementing authentication, authorization, JWT tokens, CORS, rate limiting, or any security control in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.6"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Implementing JWT authentication or OAuth 2.0
- Configuring authorization policies
- Preventing unauthorized access to resources (BOLA/IDOR)
- Setting up CORS, rate limiting, or request size limits
- Handling passwords or secrets securely
- Reviewing any endpoint that accesses user-owned data

---

## Critical Patterns

### GUID IDs — Obligatorio en todos los recursos públicos

Nunca usar `int` secuenciales en endpoints públicos — habilitan enumeración y BOLA.

```csharp
public Guid Id { get; set; } = Guid.NewGuid();
```

### BOLA Prevention — Ownership check en service layer

Todo método que accede a un recurso por ID **debe** verificar ownership. Nunca en el controller, siempre en el service:

```csharp
public async Task<OrderDto> GetByIdAsync(Guid orderId, Guid requestingUserId)
{
    var order = await _repository.GetByIdAsync(orderId);
    if (order is null) throw new NotFoundException(nameof(Order), orderId);
    if (order.UserId != requestingUserId)
        throw new ForbiddenException("You do not have access to this resource.");
    return order.ToDto();
}
```

El `requestingUserId` viene del JWT claim en el controller:

```csharp
// Guid.Parse with ! throws NullReferenceException if the claim is absent — use TryParse instead.
if (!Guid.TryParse(User.FindFirstValue(ClaimTypes.NameIdentifier), out var userId))
    return Unauthorized(); // missing or malformed identity claim — the request is unauthenticated
var data = await _orderService.GetByIdAsync(id, userId);
return Ok(data);
```

### `ForbiddenException` — Domain exception para 403

```csharp
// Shared/Exceptions/ForbiddenException.cs
public class ForbiddenException : Exception
{
    public ForbiddenException(string message) : base(message) { }
}
```

> Cross-ref: `ExceptionHandlingMiddleware` mapea `ForbiddenException` → HTTP 403. Ver [`../responses/SKILL.md`](../responses/SKILL.md).
> **`401` NO es una domain exception** — es propiedad del pipeline de auth. No crear `UnauthorizedException`.

### Authentication Error Messages — Nunca revelar existencia de usuario

```csharp
// CORRECT — plain message, no envelope
return Unauthorized(new { message = "Invalid credentials." });

// WRONG — reveals user existence
return Unauthorized(new { message = "User not found." });
```

> **Note**: Never wrap controller responses in `ResponseDTO<T>`. See [`../responses/SKILL.md`](../responses/SKILL.md) for the dual-contract pattern.

### Rate Limiting — Algoritmo por escenario

> **Middleware order (W-16)**: `app.UseRateLimiter()` MUST come BEFORE `app.UseCors()` in `Program.cs`. Inverting the order means CORS preflight requests (`OPTIONS`) consume rate limit slots and may be blocked, causing legitimate cross-origin requests to fail with 429 before CORS headers are evaluated.
>
> ```csharp
> // ✅ Correct order
> app.UseRateLimiter();
> app.UseCors();
> app.UseAuthentication();
> app.UseAuthorization();
>
> // ❌ Wrong — CORS OPTIONS requests burn rate limit slots
> app.UseCors();
> app.UseRateLimiter();
> ```

| Escenario | Algoritmo | Razón |
|---|---|---|
| Auth endpoints (login, register, forgot-password) | `SlidingWindowLimiter` | Previene burst en el límite de ventana — más preciso contra brute-force |
| API general (CRUD, queries) | `FixedWindowLimiter` | Simple, predecible, suficiente para tráfico regular |
| APIs con tráfico irregular / spikes | `TokenBucketLimiter` | Permite bursts controlados sin rechazar tráfico legítimo |

```csharp
builder.Services.AddRateLimiter(options =>
{
    options.AddSlidingWindowLimiter("auth", o =>
    {
        o.PermitLimit = 5;
        o.Window = TimeSpan.FromMinutes(1);
        o.SegmentsPerWindow = 4;
        o.QueueLimit = 0;
    });
    options.AddFixedWindowLimiter("api", o =>
    {
        o.PermitLimit = 100;
        o.Window = TimeSpan.FromMinutes(1);
    });
    options.RejectionStatusCode = StatusCodes.Status429TooManyRequests;
});

[HttpPost("login")]
[EnableRateLimiting("auth")]
public async Task<IActionResult> Login([FromBody] LoginRequest request) { ... }
```

> Para configuración `OnRejected` con `ProblemDetails` → ver [`../responses/SKILL.md`](../responses/SKILL.md)

### CORS — Nunca `AllowAnyOrigin` en producción

`AllowAnyOrigin()` + `AllowCredentials()` → `InvalidOperationException` en runtime. ASP.NET Core lo prohíbe explícitamente.

```csharp
policy.WithOrigins(builder.Configuration.GetSection("Cors:AllowedOrigins").Get<string[]>()!)
      .AllowAnyHeader()
      .AllowAnyMethod()
      .AllowCredentials();
```

### Secrets — Nunca en `appsettings.json`

| Item | Storage |
|---|---|
| JWT Secret | Environment variable / Azure Key Vault |
| DB Connection String | Environment variable / User Secrets (dev) |
| API Keys | Environment variable / Key Vault |

```csharp
var jwtSecret = builder.Configuration["Jwt:Secret"]
    ?? throw new InvalidOperationException("JWT secret not configured.");
```

---

## Anti-Patterns

| Anti-pattern | Por qué es peligroso |
|---|---|
| `int` ID en endpoints públicos | Enumeración de recursos — BOLA trivial |
| Ownership check en controller en vez de service | Se puede bypassear; el service es el boundary real |
| `AllowAnyOrigin()` en producción | CSRF y exfiltración desde cualquier dominio |
| `AllowAnyOrigin()` + `AllowCredentials()` combinados | `InvalidOperationException` en runtime |
| Error messages que revelan existencia de usuario | Habilita user enumeration |
| Secret en `appsettings.json` committeado | Credenciales expuestas en git history |
| Crear `UnauthorizedException` como domain exception | `401` es del pipeline/auth, no del dominio |

---

## Commands

```bash
# Generate a secure JWT secret (256-bit)
[System.Convert]::ToBase64String([System.Security.Cryptography.RandomNumberGenerator]::GetBytes(32))

dotnet add package Microsoft.AspNetCore.Authentication.JwtBearer
dotnet user-secrets init
dotnet user-secrets set "Jwt:Secret" "your-secret-here"
```

---

## Resources

- **HTTP responses & exception mapping**: See [`../responses/SKILL.md`](../responses/SKILL.md)
- **Input validation**: See [`../validations/SKILL.md`](../validations/SKILL.md)
- **Testing ownership checks & policies**: See [`../testing-unit/SKILL.md`](../testing-unit/SKILL.md)
- **Performance / Rate Limiting (advanced config)**: See [`../performance/SKILL.md`](../performance/SKILL.md)
- **General conventions**: See [`../general/SKILL.md`](../general/SKILL.md)
- **OWASP API Security 2023**: https://owasp.org/API-Security/editions/2023/en/0x00-introduction/

---

## Changelog

### v1.6 — 2026-03-28
- **Fixed (W-16)**: Added middleware order warning — `UseRateLimiter()` MUST come before `UseCors()`. Inverted order causes CORS `OPTIONS` preflight requests to consume rate limit slots, resulting in legitimate 429 rejections before CORS headers are evaluated.

### v1.5 — 2026-03-28
- **Fixed (FIX-B)**: Annotated the v1.3 changelog entry for C-04 to record that `Forbid()` / 403 was subsequently identified as semantically wrong for a missing claim, and was corrected to `Unauthorized()` / 401 in v1.4. Prevents an agent reading v1.3 as baseline from reinstating the incorrect pattern.

### v1.4 — 2026-03-28
- **Fixed (FIX-2)**: Replaced `return Forbid()` (HTTP 403) with `return Unauthorized()` (HTTP 401) in the ownership guard when the identity claim is absent or malformed. A missing/malformed claim means the request is unauthenticated — 401 is semantically correct; 403 ("known identity, insufficient permission") is wrong here. Also removed the alternative `throw new ForbiddenException(...)` from the inline comment, which compounded the semantic error.

### v1.3 — 2026-03-28
- **Fixed (C-04)**: Replaced `Guid.Parse(User.FindFirstValue(ClaimTypes.NameIdentifier)!)` with `Guid.TryParse(...)` + guard returning `Forbid()`. The null-forgiving operator on a potentially-null claim caused `NullReferenceException` (raw 500) for requests with absent or malformed identity claims. (Note: `Forbid()` / 403 was subsequently identified as semantically wrong for a missing claim — corrected to `Unauthorized()` / 401 in v1.4)

### v1.2 — 2026-03-28
- **Removed**: OWASP coverage map table y descripciones (doc oficial — el agente la conoce)
- **Removed**: JWT setup boilerplate completo (`AddAuthentication`, `TokenValidationParameters`) — estándar del framework
- **Removed**: Password hashing con `PasswordHasher<T>` — estándar de ASP.NET Core Identity
- **Removed**: Policy-based authorization examples — estándar del framework
- **Kept**: GUIDs obligatorios, ownership check en service, ForbiddenException, algoritmos rate limiter, auth error messages, CORS gotcha, secrets rule

### v1.1 — 2026-03-28
- Fixed: controller example usaba `ResponseDTO<object>` en auth error → corregido a `new { message = "..." }`
- Added: ForbiddenException cross-ref a responses/SKILL.md
