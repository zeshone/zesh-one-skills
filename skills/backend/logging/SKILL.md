---
name: net8-apirest-logging
description: >
  Structured logging, end-to-end traceability, and canonical exception handling for ASP.NET Core 8 REST APIs using Serilog with ProblemDetails, CorrelationId, and mandatory observability fields.
  Trigger: When implementing logging, configuring Serilog, writing exception handlers, adding Correlation ID tracking, or diagnosing incidents in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.5"
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

### Structured Logging — Named placeholders, never string interpolation

```csharp
// CORRECT — structured, queryable
_logger.LogInformation("User {UserId} retrieved {Count} records", userId, count);

// WRONG — not structured, not queryable
_logger.LogInformation($"User {userId} retrieved {count} records");
```

### Log Levels — HTTP Status Mapping

| HTTP Status | Log Level |
|---|---|
| `400` validation errors | `LogInformation` — expected flow |
| `403`, `404`, `409` domain errors | `LogWarning` — recoverable |
| `500` unhandled exceptions | `LogError` |
| Unrecoverable / startup failures | `LogCritical` |

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

> `UserId` is NOT enriched here — `Authentication` has not run yet. `UserId` enrichment is done in `ExceptionHandlingMiddleware`.

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
            catch (ValidationException ex)
            {
                _logger.LogInformation(ex, "Validation error for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status400BadRequest, "Validation Error", ex.Message);
            }
            catch (ForbiddenException ex)
            {
                _logger.LogWarning(ex, "Forbidden access for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status403Forbidden, "Forbidden", ex.Message);
            }
            catch (NotFoundException ex)
            {
                _logger.LogWarning(ex, "Resource not found for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status404NotFound, "Not Found", ex.Message);
            }
            catch (ConflictException ex)
            {
                _logger.LogWarning(ex, "Conflict for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status409Conflict, "Conflict", ex.Message);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Unhandled exception for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status500InternalServerError,
                    "Internal Server Error", "An unexpected error occurred.");
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
    .Enrich.WithMachineName()
    .Enrich.WithEnvironmentName()
    .WriteTo.Console(outputTemplate:
        "[{Timestamp:HH:mm:ss} {Level:u3}] [{CorrelationId}] [{Application}/{Environment}] {Message:lj}{NewLine}{Exception}")
    .WriteTo.File("logs/api-.log", rollingInterval: RollingInterval.Day)
    .CreateLogger();

builder.Host.UseSerilog();
```

> **Non-obvious**: The `outputTemplate` above injects `CorrelationId`, `Application`, and `Environment` — these fields are pushed via `LogContext` by `CorrelationIdMiddleware`. The `rollingInterval: RollingInterval.Day` keeps log files manageable without configuration overhead.

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
