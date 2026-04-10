---
name: net8-apirest-logging
description: >
  Structured logging, end-to-end traceability, and canonical exception handling for ASP.NET Core 8 REST APIs using Serilog with ProblemDetails, CorrelationId, and mandatory observability fields.
  Trigger: When implementing logging, configuring Serilog, writing exception handlers, adding Correlation ID tracking, or diagnosing incidents in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.8"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Configuring Serilog in a new or existing API
- Implementing the global exception handling middleware
- Adding Correlation ID tracking
- Deciding log levels for different scenarios
- Diagnosing incidents with correlationId

---

## Critical Patterns

### Serilog — Canonical Library

Always use **Serilog** for structured logging. Never `Console.WriteLine` in production code.

Required packages: `Serilog.AspNetCore`, `Serilog.Sinks.Console`, `Serilog.Sinks.File`, `Serilog.Enrichers.Environment`.

### Log Levels — HTTP Status Mapping

| HTTP Status | Log Level |
|---|---|
| `400` validation errors | `LogInformation` — expected flow, not a warning |

### Correlation ID Middleware — FIRST in the pipeline

```csharp
public class CorrelationIdMiddleware
{
    private const string CorrelationIdHeader = "X-Correlation-ID";
    private readonly RequestDelegate _next;
    private readonly string _applicationName;
    private readonly string _environment;

    public CorrelationIdMiddleware(RequestDelegate next, IConfiguration configuration, IWebHostEnvironment env)
    {
        _next = next;
        _applicationName = configuration["ApplicationName"] ?? "API";
        _environment = env.EnvironmentName;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        // Security: sanitize the caller-supplied header value to prevent log injection.
        // Unsanitized values can contain newlines or structured payloads that forge log entries.
        string rawId = context.Request.Headers.TryGetValue(CorrelationIdHeader, out var headerVal)
            ? headerVal.ToString()
            : Guid.NewGuid().ToString();

        // Strip non-printable/control characters and enforce a max length of 64 chars.
        var sanitized = new string(rawId.Where(c => !char.IsControl(c)).ToArray());
        if (sanitized.Length > 64) sanitized = sanitized[..64];
        var correlationId = string.IsNullOrWhiteSpace(sanitized) ? Guid.NewGuid().ToString() : sanitized;

        context.Items["CorrelationId"] = correlationId;
        context.Response.Headers[CorrelationIdHeader] = correlationId;

        using (LogContext.PushProperty("CorrelationId", correlationId))
        using (LogContext.PushProperty("Application", _applicationName))
        using (LogContext.PushProperty("Environment", _environment))
        {
            await _next(context);
        }
    }
}
```

> `UserId` is NOT enriched here — `Authentication` has not run yet. `UserId` is pushed in `ExceptionHandlingMiddleware` (error context) and optionally by `UserLogContextMiddleware` (when per-user log partitioning is active).

### Exception Handling Middleware — Canonical

```csharp
public class ExceptionHandlingMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<ExceptionHandlingMiddleware> _logger;

    public ExceptionHandlingMiddleware(RequestDelegate next, ILogger<ExceptionHandlingMiddleware> logger)
    {
        _next = next;
        _logger = logger;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        var userId = context.User?.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "anonymous";
        using (LogContext.PushProperty("UserId", userId))
        {
            try { await _next(context); }
            catch (ValidationException ex) // FluentValidation.ValidationException
            {
                using (LogContext.PushProperty("ExceptionType", ex.GetType().Name))
                {
                    _logger.LogInformation(ex, "Validation error for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                        userId, context.Request.Path, context.Items["CorrelationId"]);
                    await WriteErrorResponse(context, StatusCodes.Status400BadRequest, "Validation Error", ex.Message);
                }
            }
            catch (ForbiddenException ex)
            {
                using (LogContext.PushProperty("ExceptionType", ex.GetType().Name))
                {
                    _logger.LogWarning(ex, "Forbidden access for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                        userId, context.Request.Path, context.Items["CorrelationId"]);
                    await WriteErrorResponse(context, StatusCodes.Status403Forbidden, "Forbidden", ex.Message);
                }
            }
            catch (NotFoundException ex)
            {
                using (LogContext.PushProperty("ExceptionType", ex.GetType().Name))
                {
                    _logger.LogWarning(ex, "Resource not found for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                        userId, context.Request.Path, context.Items["CorrelationId"]);
                    await WriteErrorResponse(context, StatusCodes.Status404NotFound, "Not Found", ex.Message);
                }
            }
            catch (ConflictException ex)
            {
                using (LogContext.PushProperty("ExceptionType", ex.GetType().Name))
                {
                    _logger.LogWarning(ex, "Conflict for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                        userId, context.Request.Path, context.Items["CorrelationId"]);
                    await WriteErrorResponse(context, StatusCodes.Status409Conflict, "Conflict", ex.Message);
                }
            }
            catch (Exception ex)
            {
                using (LogContext.PushProperty("ExceptionType", ex.GetType().Name))
                {
                    _logger.LogError(ex, "Unhandled exception for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                        userId, context.Request.Path, context.Items["CorrelationId"]);
                    await WriteErrorResponse(context, StatusCodes.Status500InternalServerError,
                        "Internal Server Error", "An unexpected error occurred.");
                }
            }
        }
    }

    private static async Task WriteErrorResponse(HttpContext context, int statusCode, string title, string detail)
    {
        context.Response.StatusCode = statusCode;
        context.Response.ContentType = "application/problem+json";
        var problem = new ProblemDetails { Status = statusCode, Title = title, Detail = detail };
        problem.Extensions["correlationId"] = context.Items["CorrelationId"] as string;
        await context.Response.WriteAsJsonAsync(problem);
    }
}
```

> **Canonical Content-Type**: always `application/problem+json` (not `application/json`) for error responses.

For the exception → HTTP status mapping contract, see [`responses/SKILL.md`](../responses/SKILL.md).

### 9 Mandatory Minimum Traceability Fields

Every `Warning` log or above MUST contain these fields:

| Field | Source | When enriched |
|---|---|---|
| `CorrelationId` | `HttpContext.Items["CorrelationId"]` | `CorrelationIdMiddleware` (start of pipeline) |
| `Application` | `IConfiguration["ApplicationName"]` | `CorrelationIdMiddleware` |
| `Environment` | `IWebHostEnvironment.EnvironmentName` | `CorrelationIdMiddleware` |
| `UserId` | `HttpContext.User` claims | `ExceptionHandlingMiddleware` (after auth) |
| `RequestPath` | `HttpContext.Request.Path` | `UseSerilogRequestLogging` |
| `RequestMethod` | `HttpContext.Request.Method` | `UseSerilogRequestLogging` |
| `StatusCode` | `HttpContext.Response.StatusCode` | `UseSerilogRequestLogging` (at request completion) |
| `ExceptionType` | `exception.GetType().Name` | `ExceptionHandlingMiddleware` (in catch) |
| `Duration` | Calculated by Serilog | `UseSerilogRequestLogging` (at request completion) |

### Serilog Bootstrap

```csharp
Log.Logger = new LoggerConfiguration()
    .MinimumLevel.Information()
    .MinimumLevel.Override("Microsoft", LogEventLevel.Warning)
    .MinimumLevel.Override("System", LogEventLevel.Warning)
    .Enrich.FromLogContext()
    .WriteTo.Console(outputTemplate:
        "[{Timestamp:HH:mm:ss} {Level:u3}] [{CorrelationId}] [{Application}/{Environment}] {Message:lj}{NewLine}{Exception}")
    .WriteTo.File("logs/api-.log", rollingInterval: RollingInterval.Day)
    .CreateLogger();

builder.Host.UseSerilog();
```

> **Non-obvious**: The `outputTemplate` above injects `CorrelationId`, `Application`, and `Environment` — these fields are pushed via `LogContext` by `CorrelationIdMiddleware`. The `rollingInterval: RollingInterval.Day` keeps log files manageable without configuration overhead.
>
> **Production path note**: `"logs/api-.log"` is a relative path and resolves from the process working directory. In production, prefer an absolute path from configuration, for example `Path.Combine(builder.Configuration["Logging:BasePath"]!, "api-.log")`.

### Request Logging with Serilog

```csharp
app.UseSerilogRequestLogging(options =>
{
    options.EnrichDiagnosticContext = (diagnosticContext, httpContext) =>
    {
        diagnosticContext.Set("RequestHost", httpContext.Request.Host.Value);
        diagnosticContext.Set("UserAgent", httpContext.Request.Headers["User-Agent"].ToString());
        diagnosticContext.Set("CorrelationId", httpContext.Items["CorrelationId"] as string);
        diagnosticContext.Set("UserId", httpContext.User?.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "anonymous");
    };
});
```

> `StatusCode` and `Duration` are captured automatically by Serilog. Do not duplicate them.

### Per-User Log Partitioning — `WriteTo.Map`

For applications with distinct users (agents, customers, admins, tenants), partition log files by user identity using `WriteTo.Map`. This makes per-user debugging trivial — no log queries needed.

```csharp
// Install: Serilog.Sinks.Map
// The property name ("UserId" in this example) must be pushed to LogContext
// BEFORE WriteTo.Map is evaluated — use a middleware that runs after Authentication.

Log.Logger = new LoggerConfiguration()
    .MinimumLevel.Information()
    .Enrich.FromLogContext()
    .WriteTo.Map(
        keyPropertyName: "UserId",       // LogContext property to partition by
        defaultKey: "anonymous",         // file when UserId is not yet available
        configure: (userId, writeTo) =>
            writeTo.File(
                path: $"logs/{userId}/log-.txt",
                rollingInterval: RollingInterval.Day,
                retainedFileCountLimit: 31,
                outputTemplate: "[{Timestamp:HH:mm:ss} {Level:u3}] [{CorrelationId}] {Message:lj}{NewLine}{Exception}"))
    .CreateLogger();
```

**Middleware to push user identity into `LogContext`** (runs after `UseAuthentication`):

```csharp
public class UserLogContextMiddleware
{
    private readonly RequestDelegate _next;

    public UserLogContextMiddleware(RequestDelegate next) => _next = next;

    public async Task InvokeAsync(HttpContext context)
    {
        // Resolve the identity key from the claim appropriate to your domain:
        // ClaimTypes.NameIdentifier for user ID, a custom claim for tenant/agent/etc.
        var userId = context.User?.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "anonymous";

        using (LogContext.PushProperty("UserId", userId))
        {
            await _next(context);
        }
    }
}
```

**Register after `UseAuthentication()` and before `UseRateLimiter()` in `Program.cs` (see pipeline order in [`../general/SKILL.md`](../general/SKILL.md)):**
```csharp
app.UseAuthentication();
app.UseMiddleware<UserLogContextMiddleware>(); // must be AFTER auth — identity not available before
app.UseRateLimiter();
```

> **When to use**: multi-tenant SaaS, agent/broker portals, admin panels with multiple operator roles — any scenario where isolating one user's activity in a dedicated file saves significant debugging time.
>
> **When NOT to use**: public APIs with anonymous traffic or very high user counts (thousands of distinct files per day creates filesystem pressure). For high-cardinality user populations, use structured logging with a log aggregator (Seq, Grafana Loki) and query by `UserId` field instead.
>
> **The partition key is your domain**: use whatever claim your domain uses — user ID, tenant code, agent code, organization ID. The pattern is identical; only the claim name and path template change.
>
> **Production path note**: In production, use an absolute path from configuration (for example `Path.Combine(builder.Configuration["Logging:BasePath"]!, userId, "log-.txt")`) — relative paths resolve to the process working directory.

### Sensitive Data — Fields Forbidden in Logs

Sanitize BEFORE logging, never after:

| Category | Forbidden fields |
|---|---|
| Credentials | `password`, `passwordHash`, `pin`, `secret` |
| Tokens | `jwtToken`, `accessToken`, `refreshToken`, `apiKey`, `bearerToken` |
| Financial data | `cardNumber`, `cvv`, `bankAccount`, `iban` |
| PII | `email` (full), `phone`, `nationalId`, `ssn`, `dateOfBirth` |
| Infrastructure | `connectionString`, `databasePassword`, `smtpPassword` |

### Troubleshooting — Evidence Chain

```
1. User receives error → obtains correlationId from ProblemDetails.Extensions["correlationId"]
2. Support searches: WHERE CorrelationId = '{value}'
3. Reconstruction:
   - Request log (Serilog): Method, Path, StatusCode, Duration, UserId
   - Exception log (ExceptionHandlingMiddleware): ExceptionType, Message, StackTrace (internal)
   - Service logs: Operation, parameters (without PII)
```

---

## Cross-References

| Skill | Relationship |
|---|---|
| [`responses`](../responses/SKILL.md) | Defines the exception → HTTP contract. The middleware implementation lives here. |
| [`general`](../general/SKILL.md) | Defines the pipeline order. `CorrelationIdMiddleware` goes at position 0. |
| [`security`](../security/SKILL.md) | `ForbiddenException` and the prohibition on logging tokens. |

---

## Changelog

### v2.7 — 2026-04-09
- **Fixed (Round 4)**: Added explicit production-path notes for both the bootstrap file sink and the per-user `WriteTo.Map` sink. Relative `logs/...` paths resolve from the process working directory, so production deployments should use absolute paths from configuration.

### v2.6 — 2026-04-09
- **Added**: Per-User Log Partitioning with `WriteTo.Map` — partition log files by any user identity claim (user ID, tenant, agent, organization). Includes `UserLogContextMiddleware`, bootstrap configuration, and guidance on when to use vs when a log aggregator is the better choice. Pattern derived from production multi-tenant microservices audit.

### v2.5 — 2026-04-09
- **Fixed (W-07)**: Added explanatory note to the Serilog bootstrap block. The `outputTemplate` and `rollingInterval` are non-obvious ZeshOne decisions (custom fields from `CorrelationIdMiddleware`, daily rotation as default). The rest of the bootstrap is standard Serilog configuration the agent already knows.

### v2.4 — 2026-03-28
- **Fixed (FIX-C)**: Changed `context.Items["CorrelationId"]?.ToString()` to `context.Items["CorrelationId"] as string` in two locations: `WriteErrorResponse` (`problem.Extensions["correlationId"]`) and `UseSerilogRequestLogging` enricher (`diagnosticContext.Set("CorrelationId", ...)`). `Items["CorrelationId"]` is typed `object?` but is always stored as `string` by `CorrelationIdMiddleware`; `as string` is explicit about the type and consistent with FIX-6's philosophy of removing redundant `.ToString()` calls.

### v2.3 — 2026-03-28
- **Fixed (FIX-6)**: Removed redundant `.ToString()` on `correlationId` in `LogContext.PushProperty(...)`. The variable is already typed as `string`; calling `.ToString()` implies type uncertainty and adds noise.

### v2.2 — 2026-03-28
- **Fixed (C-03)**: `CorrelationIdMiddleware` now sanitizes the caller-supplied `X-Correlation-ID` header before storing it. Non-printable/control characters are stripped and the value is capped at 64 characters to prevent log-injection attacks. A fallback to a fresh GUID is applied if the sanitized value is empty.

### v2.1 — 2026-03-28
- **Removed**: `IExceptionHandler` section as a non-canonical alternative (noise — if someone chooses it, they know how to use it)
- **Removed**: `appsettings.json` Serilog config example (standard config, not a decision)
- **Removed**: Duplicate NuGet commands at the end
- **Removed**: Phase 2 scope boundary (out of scope section)
- **Kept**: All canonical middleware, 9 fields, sensitive data list, Serilog bootstrap, request logging, troubleshooting chain

### v2.0 — 2026-03-25
- BREAKING: Error responses migrated to ProblemDetails; correlationId in Extensions; 9 traceability fields
