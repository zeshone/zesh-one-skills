---
name: net8-apirest-security
description: >
  Security standards for ASP.NET Core 8 REST APIs based on OWASP API Security Top 10 2023, covering authorization, authentication, BOLA prevention, input validation, rate limiting, and secure credential handling.
  Trigger: When implementing authentication, authorization, JWT tokens, CORS, rate limiting, or any security control in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.1"
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

### OWASP API Security Top 10 — Coverage Map

| OWASP Risk | Implementation |
|---|---|
| API1 — Broken Object Level Authorization (BOLA) | Ownership check in service layer + GUID IDs |
| API2 — Broken Authentication | JWT with short expiry + refresh tokens + brute-force protection |
| API3 — Broken Object Property Level Authorization | Role-specific DTOs, never expose full entities |
| API4 — Unrestricted Resource Consumption | Rate limiting + max page size + request body size limit |
| API5 — Broken Function Level Authorization | Policy-based authorization per endpoint |
| API7 — Server-Side Request Forgery (SSRF) | Validate and allowlist outbound URLs |
| API8 — Security Misconfiguration | Disable dev-only features in production |

> **Out of scope**: API6 (Unrestricted Access to Sensitive Business Flows), API9 (Improper Inventory Management), and API10 (Unsafe Consumption of APIs) are not covered — they require infrastructure-level controls beyond the scope of this API-level skill.

### API1 — BOLA Prevention (Mandatory for all user-owned resources)

Every service method that accesses a resource by ID **must** verify ownership:

```csharp
public async Task<OrderDto> GetByIdAsync(Guid orderId, Guid requestingUserId)
{
    var order = await _repository.GetByIdAsync(orderId);
    if (order is null)
        throw new NotFoundException(nameof(Order), orderId);
    // Ownership check — critical
    if (order.UserId != requestingUserId)
        throw new ForbiddenException("You do not have access to this resource.");
    return order.ToDto();
}
```

Get the `requestingUserId` from the JWT claims in the controller, then pass it to the service:

```csharp
[HttpGet("{id:guid}")]
[Authorize]
public async Task<IActionResult> GetById([FromRoute] Guid id)
{
    var userId = Guid.Parse(User.FindFirstValue(ClaimTypes.NameIdentifier)!);
    var data = await _orderService.GetByIdAsync(id, userId);
    return Ok(new ResponseDTO<OrderDto> { Success = true, Data = data });
}
```

> **Use GUIDs** for all public resource identifiers. Sequential IDs enable enumeration attacks.

**`ForbiddenException` — define it in `Shared/Exceptions/`** (analogous to `NotFoundException` in `general`):

```csharp
// Shared/Exceptions/ForbiddenException.cs
public class ForbiddenException : Exception
{
    public ForbiddenException(string message)
        : base(message) { }
}
```

> Cross-ref: The `ExceptionHandlingMiddleware` (→ [`../responses/SKILL.md`](../responses/SKILL.md)) must map `ForbiddenException` → HTTP 403 Forbidden. `NotFoundException` → HTTP 401 Unauthorized is **not** the same — 401 means unauthenticated, 403 means authenticated but not authorized.

### API2 — Authentication with JWT Bearer

```csharp
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ValidIssuer = builder.Configuration["Jwt:Issuer"],
            ValidAudience = builder.Configuration["Jwt:Audience"],
            IssuerSigningKey = new SymmetricSecurityKey(
                Encoding.UTF8.GetBytes(builder.Configuration["Jwt:Secret"]!)),
            ClockSkew = TimeSpan.Zero           // no grace period on expiry
        };
    });
```

JWT token rules:
- **Expiry**: short-lived (`15min` to `1h`)
- **Refresh tokens**: long-lived, stored server-side (or `HttpOnly` cookie)
- **Secret**: minimum 256 bits, stored in environment variable or secrets manager — never in `appsettings.json`

### API2 — Secure Password Storage

Never store plain-text passwords. Use PBKDF2 or BCrypt:

```csharp
// Using ASP.NET Core built-in hasher
public class PasswordHasher : IPasswordHasher
{
    private readonly Microsoft.AspNetCore.Identity.PasswordHasher<string> _hasher = new();

    public string Hash(string password) => _hasher.HashPassword(string.Empty, password);

    public bool Verify(string hashedPassword, string plainPassword) =>
        _hasher.VerifyHashedPassword(string.Empty, hashedPassword, plainPassword)
            == PasswordVerificationResult.Success;
}
```

Authentication error messages must be **generic** — never hint at whether the user exists:
```csharp
// CORRECT
return Unauthorized(new ResponseDTO<object> { Success = false, Message = "Invalid credentials." });

// WRONG — reveals user existence
return Unauthorized(new ResponseDTO<object> { Success = false, Message = "User not found." });
```

### API3 — Role-Specific DTOs (Property Level Authorization)

Use different DTOs per role to prevent property-level data leakage:

```csharp
// Full DTO for admins
public class UserAdminDto
{
    public Guid Id { get; set; }
    public string Email { get; set; } = string.Empty;
    public string Role { get; set; } = string.Empty;
    public bool IsActive { get; set; }
    public DateTime LastLogin { get; set; }
    public string InternalNotes { get; set; } = string.Empty;
}

// Limited DTO for regular users
public class UserPublicDto
{
    public Guid Id { get; set; }
    public string DisplayName { get; set; } = string.Empty;
}
```

### API4 — Rate Limiting + Request Size Limits

Choose the right algorithm for each scenario:

| Scenario | Recommended Algorithm | Reason |
|---|---|---|
| Auth endpoints (login, register, forgot-password) | `SlidingWindowLimiter` | Prevents burst at window boundary; more precise against brute-force |
| General API (CRUD, queries) | `FixedWindowLimiter` | Simple, predictable, sufficient for regular traffic |
| APIs with irregular traffic / spikes | `TokenBucketLimiter` | Allows controlled bursts without rejecting legitimate traffic |

> For advanced configuration and tuning of rate limiters → [`../performance/SKILL.md`](../performance/SKILL.md)

```csharp
// Rate limiting
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

// Apply tighter limit on auth endpoints
[HttpPost("login")]
[EnableRateLimiting("auth")]
public async Task<IActionResult> Login([FromBody] LoginRequest request) { ... }
```

Request body size limit:
```csharp
builder.Services.Configure<KestrelServerOptions>(options =>
{
    options.Limits.MaxRequestBodySize = 10 * 1024 * 1024; // 10 MB max
});
```

### API5 — Policy-Based Authorization (Function Level)

Define granular policies and apply them explicitly:

```csharp
builder.Services.AddAuthorization(options =>
{
    options.AddPolicy("AdminOnly", policy => policy.RequireRole("Admin"));
    options.AddPolicy("OwnerOrAdmin", policy =>
        policy.RequireAssertion(ctx =>
            ctx.User.IsInRole("Admin") ||
            ctx.User.HasClaim(c => c.Type == ClaimTypes.NameIdentifier)));
});
```

Apply per endpoint — never rely on global `[Authorize]` alone for sensitive operations:
```csharp
[HttpDelete("{id:guid}")]
[Authorize(Policy = "AdminOnly")]
public async Task<IActionResult> Delete([FromRoute] Guid id) { ... }
```

### CORS — Restrictive Configuration

Never use `AllowAnyOrigin` in production. Combining `AllowAnyOrigin()` with `AllowCredentials()` causes an `InvalidOperationException` at runtime — ASP.NET Core forbids it explicitly:

```csharp
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowedOrigins", policy =>
    {
        policy.WithOrigins(builder.Configuration.GetSection("Cors:AllowedOrigins").Get<string[]>()!)
              .AllowAnyHeader()
              .AllowAnyMethod()
              .AllowCredentials();
    });
});

app.UseCors("AllowedOrigins"); // Named policy — see CORS config above
```

### Secrets Management

| Item | Storage |
|---|---|
| JWT Secret | Environment variable / Azure Key Vault |
| DB Connection String | Environment variable / User Secrets (dev) |
| API Keys | Environment variable / Key Vault |
| Any sensitive config | **Never** in `appsettings.json` committed to source control |

```csharp
// Access secret safely
var jwtSecret = builder.Configuration["Jwt:Secret"]
    ?? throw new InvalidOperationException("JWT secret not configured.");
```

---

## Anti-Patterns

| Anti-pattern | Why it's dangerous |
|---|---|
| `AllowAnyOrigin()` in production | Enables CSRF and data exfiltration from any domain |
| `AllowAnyOrigin()` + `AllowCredentials()` combined | ASP.NET Core throws `InvalidOperationException` at runtime — technically impossible |
| JWT without `ValidateLifetime = true` | Expired tokens remain valid indefinitely |
| Secret in `appsettings.json` committed to repo | Credentials exposed in git history — mandatory rotation required |
| Error messages that reveal user existence | "User not found" / "Wrong password" → enables user enumeration |
| BOLA without ownership check in service layer | Any authenticated user can access other users' resources — OWASP API1 |
| Global `[Authorize]` without granular policies | False sense of security — does not prevent horizontal/vertical privilege escalation |
| Sequential integer IDs (`int`) as public identifiers | Enables enumeration and resource scraping — use GUIDs |

---

## Code Examples

### Secure Middleware Pipeline Order

```csharp
app.UseAuthentication();      // 1. Who are you?
app.UseRateLimiter();         // 2. Too many requests?
app.UseCors("AllowedOrigins"); // 3. Allowed origin? — Named policy (see CORS config above)
app.UseMiddleware<ExceptionHandlingMiddleware>(); // 4. Catch all errors
app.UseAuthorization();       // 5. Are you allowed?
app.MapControllers();
```

### SSRF Prevention — Validate Outbound URLs

```csharp
public async Task<string> FetchExternalDataAsync(string url)
{
    var allowedHosts = new[] { "api.trusted.com", "data.myorg.com" };
    var uri = new Uri(url);
    if (!allowedHosts.Contains(uri.Host, StringComparer.OrdinalIgnoreCase))
        throw new InvalidOperationException($"Host '{uri.Host}' is not allowed.");
    return await _httpClient.GetStringAsync(uri);
}
```

---

## Commands

```bash
# Generate a secure JWT secret (256-bit)
[System.Convert]::ToBase64String([System.Security.Cryptography.RandomNumberGenerator]::GetBytes(32))

# Add JWT Bearer package
dotnet add package Microsoft.AspNetCore.Authentication.JwtBearer

# Add User Secrets (dev only)
dotnet user-secrets init
dotnet user-secrets set "Jwt:Secret" "your-secret-here"
```

---

## Changelog

### v1.1 — 2026-03-24
- **Added**: Section `## Anti-Patterns` with 8 critical security items (2-column table)
- **Added**: `ForbiddenException` definition in `Shared/Exceptions/` with cross-ref to `ExceptionHandlingMiddleware`
- **Added**: Rate limiting algorithm decision table (`SlidingWindow` vs `FixedWindow` vs `TokenBucket`)
- **Added**: Explicit note on OWASP API6/API9/API10 out-of-scope after coverage map
- **Added**: Cross-references to `responses`, `validations`, `testing-unit` in `## Resources`
- **Fixed**: `UseCors("AllowedOrigins")` — added inline comment clarifying named policy
- **Fixed**: CORS section — added explicit warning on `AllowAnyOrigin()` + `AllowCredentials()` runtime error
- **Updated**: Rate limiting auth example changed from `FixedWindowLimiter` → `SlidingWindowLimiter` (aligned with recommendation)
- **Updated**: Frontmatter version `"1.0"` → `"1.1"`

---

## Resources

- **OWASP API Security 2023**: https://owasp.org/API-Security/editions/2023/en/0x00-introduction/
- **Standards**: See [../../../../rules-to-skills/Standardized_NET_Rules.md](../../../../rules-to-skills/Standardized_NET_Rules.md)
- **HTTP responses & exception mapping**: See [`../responses/SKILL.md`](../responses/SKILL.md)
- **Input validation (API3/API4)**: See [`../validations/SKILL.md`](../validations/SKILL.md)
- **Testing ownership checks & policies**: See [`../testing-unit/SKILL.md`](../testing-unit/SKILL.md)
- **Performance / Rate Limiting (advanced config)**: See [`../performance/SKILL.md`](../performance/SKILL.md)
- **General conventions**: See [`../general/SKILL.md`](../general/SKILL.md)
