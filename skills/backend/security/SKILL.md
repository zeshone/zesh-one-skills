---
name: net8-apirest-security
description: >
  Security standards for ASP.NET Core 8 REST APIs based on OWASP API Security Top 10 2023, covering authorization, authentication, BOLA prevention, input validation, rate limiting, and secure credential handling.
  Trigger: When implementing authentication, authorization, JWT tokens, CORS, rate limiting, or any security control in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.8"
allowed-tools:
  - Read
  - Edit
  - Write
  - Bash
  - Glob
  - Grep
---

## When to Use

- Implementing JWT authentication or OAuth 2.0
- JWT bootstrap (`AddAuthentication`, `TokenValidationParameters`) is in `general/SKILL.md`. This skill covers authorization, BOLA, OWASP, rate limiting, and CORS.
- Configuring authorization policies
- Preventing unauthorized access to resources (BOLA/IDOR)
- Setting up CORS, rate limiting, or request size limits
- Handling passwords or secrets securely
- Reviewing any endpoint that accesses user-owned data

---

## Critical Patterns

### GUID IDs — Required on all public resources

Never use sequential `int` IDs on public endpoints — they enable enumeration and BOLA.

```csharp
public Guid Id { get; set; } = Guid.NewGuid();
```

### BOLA Prevention — Ownership check in service layer

Every method that accesses a resource by ID **must** verify ownership. Never in the controller, always in the service:

```csharp
// Legacy pattern — prefer Result<OrderDto> in new services (see ../responses/SKILL.md)
public async Task<OrderDto> GetByIdAsync(Guid orderId, Guid requestingUserId)
{
    var order = await _repository.GetByIdAsync(orderId);
    if (order is null) throw new NotFoundException(nameof(Order), orderId);
    if (order.UserId != requestingUserId)
        throw new ForbiddenException("You do not have access to this resource.");
    return order.ToDto();
}
```

The `requestingUserId` comes from the JWT claim in the controller:

```csharp
// Guid.Parse with ! throws NullReferenceException if the claim is absent — use TryParse instead.
if (!Guid.TryParse(User.FindFirstValue(ClaimTypes.NameIdentifier), out var userId))
    return Unauthorized(new ProblemDetails
    {
        Status = StatusCodes.Status401Unauthorized,
        Title = "Unauthorized",
        Detail = "Missing or invalid identity claim."
    }); // missing or malformed identity claim — the request is unauthenticated
var data = await _orderService.GetByIdAsync(id, userId);
return Ok(data);
```

### `ForbiddenException` — Domain exception for 403

```csharp
// Shared/Exceptions/ForbiddenException.cs
public class ForbiddenException : Exception
{
    public ForbiddenException(string message) : base(message) { }
}
```

> Cross-ref: `ExceptionHandlingMiddleware` maps `ForbiddenException` → HTTP 403. See [`../responses/SKILL.md`](../responses/SKILL.md).
> **`401` is NOT a domain exception** — it belongs to the auth pipeline. Do not create `UnauthorizedException`.

### Authentication Error Messages — Never reveal user existence

```csharp
// CORRECT — ProblemDetails matches the external response contract
return Unauthorized(new ProblemDetails
{
    Status = StatusCodes.Status401Unauthorized,
    Title = "Unauthorized",
    Detail = "Invalid credentials."
});

// WRONG — reveals user existence
return Unauthorized(new { message = "User not found." });
```

> **Note**: Never wrap controller responses in `ResponseDTO<T>`. See [`../responses/SKILL.md`](../responses/SKILL.md) for the response contract.

### Rate Limiting — Algorithm per scenario

> **Middleware order**: `UseRateLimiter()` must be placed in the correct pipeline position. See [`../general/SKILL.md`](../general/SKILL.md) for the canonical order.

| Scenario | Algorithm | Reason |
|---|---|---|
| Auth endpoints (login, register, forgot-password) | `SlidingWindowLimiter` | Prevents burst at the window boundary — more precise against brute-force |
| General API (CRUD, queries) | `FixedWindowLimiter` | Simple, predictable, sufficient for regular traffic |
| APIs with irregular traffic / spikes | `TokenBucketLimiter` | Allows controlled bursts without rejecting legitimate traffic |

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

> For `OnRejected` configuration with `ProblemDetails` → see [`../responses/SKILL.md`](../responses/SKILL.md)

### CORS — Never `AllowAnyOrigin` in production

`AllowAnyOrigin()` + `AllowCredentials()` → `InvalidOperationException` at runtime. ASP.NET Core explicitly forbids this.

```csharp
policy.WithOrigins(builder.Configuration.GetSection("Cors:AllowedOrigins").Get<string[]>()!)
      .AllowAnyHeader()
      .AllowAnyMethod()
      .AllowCredentials();
```

### Secrets — Never in `appsettings.json`

| Item | Storage |
|---|---|
| JWT Secret | External config file (path provided by env var) / Azure Key Vault |
| DB Connection String | External config file (path provided by env var) / User Secrets (dev) |
| API Keys | External config file (path provided by env var) / Key Vault |

```csharp
var jwtSecret = builder.Configuration["Jwt:Secret"]
    ?? throw new InvalidOperationException("JWT secret not configured.");
```

> See [`../general/SKILL.md`](../general/SKILL.md) — Configuration section for the external file pattern. The environment variable points to the config file path, not to the secret value itself.
>
> To generate a secure secret in PowerShell: `[System.Convert]::ToBase64String([System.Security.Cryptography.RandomNumberGenerator]::GetBytes(32))`

---

## Anti-Patterns

| Anti-pattern | Why it is dangerous |
|---|---|
| `int` ID on public endpoints | Resource enumeration — trivial BOLA |
| Ownership check in controller instead of service | Can be bypassed; the service is the real boundary |
| `AllowAnyOrigin()` in production | CSRF and data exfiltration from any domain |
| `AllowAnyOrigin()` + `AllowCredentials()` combined | `InvalidOperationException` at runtime |
| Error messages that reveal user existence | Enables user enumeration |
| Secret in committed `appsettings.json` | Credentials exposed in git history |
| Creating `UnauthorizedException` as a domain exception | `401` belongs to the pipeline/auth, not the domain |

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

### v1.8 — 2026-04-09
- **Fixed (Round 3)**: Replaced the bare `return Unauthorized();` in the BOLA ownership guard example with an explicit `ProblemDetails` payload for consistency with the canonical auth error contract.

### v1.7 — 2026-04-09
- **Fixed (CRITICAL — pipeline contradiction)**: Removed inline pipeline order block from Rate Limiting section. It contradicted `general/SKILL.md` on middleware positions. Replaced with a cross-ref to the canonical pipeline order in `general`.
- **Fixed**: Removed `dotnet add package` and `dotnet user-secrets` commands from Commands section — generic CLI knowledge, not project decisions. The PowerShell JWT secret generator is now inlined as a note in the Secrets section instead of keeping a dedicated Commands heading.

### v1.6 — 2026-03-28
- **Fixed (W-16)**: Added middleware order warning — `UseRateLimiter()` MUST come before `UseCors()`. Inverted order causes CORS `OPTIONS` preflight requests to consume rate limit slots, resulting in legitimate 429 rejections before CORS headers are evaluated.

### v1.5 — 2026-03-28
- **Fixed (FIX-B)**: Annotated the v1.3 changelog entry for C-04 to record that `Forbid()` / 403 was subsequently identified as semantically wrong for a missing claim, and was corrected to `Unauthorized()` / 401 in v1.4. Prevents an agent reading v1.3 as baseline from reinstating the incorrect pattern.

### v1.4 — 2026-03-28
- **Fixed (FIX-2)**: Replaced `return Forbid()` (HTTP 403) with `return Unauthorized()` (HTTP 401) in the ownership guard when the identity claim is absent or malformed. A missing/malformed claim means the request is unauthenticated — 401 is semantically correct; 403 ("known identity, insufficient permission") is wrong here. Also removed the alternative `throw new ForbiddenException(...)` from the inline comment, which compounded the semantic error.

### v1.3 — 2026-03-28
- **Fixed (C-04)**: Replaced `Guid.Parse(User.FindFirstValue(ClaimTypes.NameIdentifier)!)` with `Guid.TryParse(...)` + guard returning `Forbid()`. The null-forgiving operator on a potentially-null claim caused `NullReferenceException` (raw 500) for requests with absent or malformed identity claims. (Note: `Forbid()` / 403 was subsequently identified as semantically wrong for a missing claim — corrected to `Unauthorized()` / 401 in v1.4)

### v1.2 — 2026-03-28
- **Removed**: OWASP coverage map table and descriptions (official docs — the agent already knows them)
- **Removed**: Full JWT setup boilerplate (`AddAuthentication`, `TokenValidationParameters`) — framework standard
- **Removed**: Password hashing with `PasswordHasher<T>` — ASP.NET Core Identity standard
- **Removed**: Policy-based authorization examples — framework standard
- **Kept**: Required GUIDs, ownership check in service, ForbiddenException, rate limiter algorithms, auth error messages, CORS gotcha, secrets rule

### v1.1 — 2026-03-28
- Fixed: controller example used `ResponseDTO<object>` in auth error → corrected to `new { message = "..." }`
- Added: ForbiddenException cross-ref to responses/SKILL.md
