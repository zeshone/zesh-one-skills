---
name: net8-apirest-logging
description: >
  Structured logging, end-to-end traceability, and canonical exception handling for ASP.NET Core 8 REST APIs using Serilog with ProblemDetails, CorrelationId, and mandatory observability fields.
  Trigger: When implementing logging, configuring Serilog, writing exception handlers, adding Correlation ID tracking, or diagnosing incidents in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.4"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Configurando Serilog en una API nueva o existente
- Implementando el global exception handling middleware
- Agregando Correlation ID tracking
- Decidiendo log levels para distintos escenarios
- Diagnosticando incidentes con correlationId

---

## Critical Patterns

### Serilog — Librería canónica

Siempre usar **Serilog** para structured logging. Nunca `Console.WriteLine` en código de producción.

Packages requeridos:
```bash
dotnet add package Serilog.AspNetCore
dotnet add package Serilog.Sinks.Console
dotnet add package Serilog.Sinks.File
dotnet add package Serilog.Enrichers.Environment
```

### Structured Logging — Named placeholders, nunca string interpolation

```csharp
// CORRECT — estructurado, queryable
_logger.LogInformation("User {UserId} retrieved {Count} records", userId, count);

// WRONG — no estructurado, no queryable
_logger.LogInformation($"User {userId} retrieved {count} records");
```

### Log Levels

| Level | Cuándo usar |
|---|---|
| `LogDebug` | Estado interno detallado, solo en desarrollo |
| `LogInformation` | Eventos de operación normal + errores de validación (400) — flujo esperado |
| `LogWarning` | Situaciones recuperables + 403, 404, 409 — errores de dominio |
| `LogError` | Excepciones no manejadas que afectan la operación (500) |
| `LogCritical` | Fallos irrecuperables que requieren atención inmediata |

### Correlation ID Middleware — PRIMERO en el pipeline

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

> `UserId` NO se enriquece aquí — `Authentication` aún no ejecutó. El enriquecimiento de `UserId` se hace en `ExceptionHandlingMiddleware`.

### Exception Handling Middleware — Canónico

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

> **Content-Type canónico**: siempre `application/problem+json` (no `application/json`) para respuestas de error.

Para el contrato de mapeo exception → HTTP status → ver [`responses/SKILL.md`](../responses/SKILL.md).

### 9 Campos Mínimos Obligatorios de Trazabilidad

Todo log `Warning` o superior DEBE contener estos campos:

| Campo | Fuente | Cuándo se enriquece |
|---|---|---|
| `CorrelationId` | `HttpContext.Items["CorrelationId"]` | `CorrelationIdMiddleware` (inicio del pipeline) |
| `Application` | `IConfiguration["ApplicationName"]` | `CorrelationIdMiddleware` |
| `Environment` | `IWebHostEnvironment.EnvironmentName` | `CorrelationIdMiddleware` |
| `UserId` | `HttpContext.User` claims | `ExceptionHandlingMiddleware` (después de auth) |
| `RequestPath` | `HttpContext.Request.Path` | `UseSerilogRequestLogging` |
| `RequestMethod` | `HttpContext.Request.Method` | `UseSerilogRequestLogging` |
| `StatusCode` | `HttpContext.Response.StatusCode` | `UseSerilogRequestLogging` (al finalizar request) |
| `ExceptionType` | `exception.GetType().Name` | `ExceptionHandlingMiddleware` (en catch) |
| `Duration` | Calculado por Serilog | `UseSerilogRequestLogging` (al finalizar request) |

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

### Request Logging con Serilog

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

> `StatusCode` y `Duration` los captura Serilog automáticamente. No duplicarlos.

### Datos Sensibles — Campos Prohibidos en Logs

Sanitizar ANTES de loguear, nunca después:

| Categoría | Campos prohibidos |
|---|---|
| Credenciales | `password`, `passwordHash`, `pin`, `secret` |
| Tokens | `jwtToken`, `accessToken`, `refreshToken`, `apiKey`, `bearerToken` |
| Datos financieros | `cardNumber`, `cvv`, `bankAccount`, `iban` |
| PII | `email` (completo), `phone`, `nationalId`, `ssn`, `dateOfBirth` |
| Infraestructura | `connectionString`, `databasePassword`, `smtpPassword` |

### Troubleshooting — Cadena de Evidencia

```
1. Usuario recibe error → obtiene correlationId de ProblemDetails.Extensions["correlationId"]
2. Soporte busca: WHERE CorrelationId = '{valor}'
3. Reconstrucción:
   - Log de request (Serilog): Method, Path, StatusCode, Duration, UserId
   - Log de excepción (ExceptionHandlingMiddleware): ExceptionType, Message, StackTrace (interno)
   - Logs de service: Operación, parámetros (sin PII)
```

---

## Cross-References

| Skill | Relación |
|---|---|
| [`responses`](../responses/SKILL.md) | Define el contrato exception → HTTP. La implementación del middleware vive aquí. |
| [`general`](../general/SKILL.md) | Define el pipeline order. `CorrelationIdMiddleware` va en posición 0. |
| [`security`](../security/SKILL.md) | `ForbiddenException` y prohibición de loguear tokens. |

---

## Changelog

### v2.4 — 2026-03-28
- **Fixed (FIX-C)**: Changed `context.Items["CorrelationId"]?.ToString()` to `context.Items["CorrelationId"] as string` in two locations: `WriteErrorResponse` (`problem.Extensions["correlationId"]`) and `UseSerilogRequestLogging` enricher (`diagnosticContext.Set("CorrelationId", ...)`). `Items["CorrelationId"]` is typed `object?` but is always stored as `string` by `CorrelationIdMiddleware`; `as string` is explicit about the type and consistent with FIX-6's philosophy of removing redundant `.ToString()` calls.

### v2.3 — 2026-03-28
- **Fixed (FIX-6)**: Removed redundant `.ToString()` on `correlationId` in `LogContext.PushProperty(...)`. The variable is already typed as `string`; calling `.ToString()` implies type uncertainty and adds noise.

### v2.2 — 2026-03-28
- **Fixed (C-03)**: `CorrelationIdMiddleware` now sanitizes the caller-supplied `X-Correlation-ID` header before storing it. Non-printable/control characters are stripped and the value is capped at 64 characters to prevent log-injection attacks. A fallback to a fresh GUID is applied if the sanitized value is empty.

### v2.1 — 2026-03-28
- **Removed**: Sección `IExceptionHandler` como alternativa no canónica (ruido — si alguien lo elige, sabe cómo usarlo)
- **Removed**: `appsettings.json` Serilog config example (config estándar, no decisión)
- **Removed**: NuGet commands duplicados al final
- **Removed**: Frontera de scope Fase 2 (out of scope section)
- **Kept**: Todo el middleware canónico, 9 campos, sensitive data list, Serilog bootstrap, request logging, troubleshooting chain

### v2.0 — 2026-03-25
- BREAKING: Error responses migradas a ProblemDetails; correlationId en Extensions; 9 campos de trazabilidad
