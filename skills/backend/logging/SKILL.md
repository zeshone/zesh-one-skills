---
name: net8-apirest-logging
description: >
  Structured logging, end-to-end traceability, and canonical exception handling for ASP.NET Core 8 REST APIs using Serilog with ProblemDetails, CorrelationId, and mandatory observability fields.
  Trigger: When implementing logging, configuring Serilog, writing exception handlers, adding Correlation ID tracking, or diagnosing incidents in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.0"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Configuring Serilog in a new or existing API
- Adding logging to services, middlewares, or controllers
- Implementing global exception handling middleware (ProblemDetails contract)
- Adding Correlation ID tracking across requests
- Deciding log levels for different scenarios
- Diagnosing incidents reported by users (using correlationId from error response)

> **Breaking change v2.0**: Error responses now use `ProblemDetails` (RFC 9457) instead of `ResponseDTO<T>`. See [Changelog](#changelog).

---

## Critical Patterns

### Logging Library — Serilog

Always use **Serilog** for structured logging. Never use `Console.WriteLine` for diagnostic output in production code.

Required NuGet packages:
- `Serilog.AspNetCore`
- `Serilog.Sinks.Console`
- `Serilog.Sinks.File`
- `Serilog.Enrichers.Environment`

### Log Levels — When to Use Each

| Level | Use case |
|---|---|
| `LogDebug` | Detailed internal state, only during development |
| `LogInformation` | Normal operation events (start, completion of flows). **Also**: validation errors (400) — expected business flow |
| `LogWarning` | Unexpected but recoverable situations. **Also**: 403, 404, 409 — domain-level issues |
| `LogError` | Unhandled exceptions that affect operation (500) |
| `LogCritical` | Unrecoverable failures that require immediate attention |

> **Rule**: Never log sensitive data (passwords, tokens, PII). Sanitize BEFORE logging, never after. See [Datos Sensibles — Campos Prohibidos](#datos-sensibles--campos-prohibidos).

### Structured Logging — Message Templates

Use **named placeholders**, never string interpolation in log messages:

```csharp
// CORRECT — structured, queryable
_logger.LogInformation("User {UserId} retrieved {Count} records", userId, count);

// WRONG — not structured, not queryable
_logger.LogInformation($"User {userId} retrieved {count} records");
```

### Correlation ID Middleware — v2.0

Every request MUST carry a `Correlation-ID` that is:
1. Generated if not present in the incoming request
2. Pushed to Serilog context with `Application` and `Environment` enrichers
3. Stored in `HttpContext.Items` for downstream access
4. Returned in the response header

**Posición en el pipeline**: PRIMERO — antes de `Authentication` y cualquier otro middleware. Esto garantiza que TODO log del request, incluyendo errores de auth, tenga correlationId.

```csharp
public class CorrelationIdMiddleware
{
    private const string CorrelationIdHeader = "X-Correlation-ID";
    private readonly RequestDelegate _next;
    private readonly string _applicationName;
    private readonly string _environment;

    public CorrelationIdMiddleware(
        RequestDelegate next,
        IConfiguration configuration,
        IWebHostEnvironment env)
    {
        _next = next;
        _applicationName = configuration["ApplicationName"] ?? "API";
        _environment = env.EnvironmentName;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        if (!context.Request.Headers.TryGetValue(CorrelationIdHeader, out var correlationId))
            correlationId = Guid.NewGuid().ToString();

        context.Items["CorrelationId"] = correlationId.ToString();
        context.Response.Headers[CorrelationIdHeader] = correlationId.ToString();

        using (LogContext.PushProperty("CorrelationId", correlationId.ToString()))
        using (LogContext.PushProperty("Application", _applicationName))
        using (LogContext.PushProperty("Environment", _environment))
        {
            await _next(context);
        }
    }
}
```

> **Nota**: `UserId` no se enriquece aquí porque `Authentication` aún no ejecutó — `HttpContext.User` no está poblado. El enriquecimiento de `UserId` DEBE hacerse en `ExceptionHandlingMiddleware` (ver sección siguiente).

### Exception Handling Middleware — CANÓNICO

El middleware custom `ExceptionHandlingMiddleware` es el **patrón canónico** para manejo global de excepciones. Ejecuta DESPUÉS de `Authentication` (para tener `HttpContext.User` disponible) y ANTES de `Authorization`.

Responsabilidades:
- Capturar todas las excepciones no manejadas
- Mapear excepciones al código HTTP y nivel de log correspondientes
- Enriquecer el log con `UserId` (disponible aquí porque auth ya ejecutó)
- Retornar respuesta `ProblemDetails` con `Extensions["correlationId"]`
- Nunca exponer stack traces al cliente

#### Tabla Canónica de Mapeo de Excepciones

| Excepción | HTTP Status | Log Level | Semántica |
|---|---|---|---|
| `ValidationException` | 400 Bad Request | `Information` | Error de validación esperado — flujo normal de negocio |
| `ForbiddenException` | 403 Forbidden | `Warning` | Acceso denegado — usuario autenticado pero sin permisos |
| `NotFoundException` | 404 Not Found | `Warning` | Recurso no encontrado — error de dominio recuperable |
| `ConflictException` | 409 Conflict | `Warning` | Conflicto de estado — duplicado o precondición fallida |
| `Exception` (no manejada) | 500 Internal Server Error | `Error` | Error inesperado — requiere investigación |

> **Regla**: Este mapeo es el contrato canónico. NO duplicarlo en `responses/SKILL.md` ni en `general/SKILL.md`. Toda referencia cross-skill apunta aquí.

```csharp
public class ExceptionHandlingMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<ExceptionHandlingMiddleware> _logger;

    public ExceptionHandlingMiddleware(
        RequestDelegate next,
        ILogger<ExceptionHandlingMiddleware> logger)
    {
        _next = next;
        _logger = logger;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        // Enriquecemos UserId aquí: Authentication ya ejecutó, HttpContext.User está poblado
        var userId = context.User?.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "anonymous";
        using (LogContext.PushProperty("UserId", userId))
        {
            try
            {
                await _next(context);
            }
            catch (ValidationException ex)
            {
                _logger.LogInformation(ex,
                    "Validation error for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status400BadRequest,
                    "Validation Error", ex.Message);
            }
            catch (ForbiddenException ex)
            {
                _logger.LogWarning(ex,
                    "Forbidden access for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status403Forbidden,
                    "Forbidden", ex.Message);
            }
            catch (NotFoundException ex)
            {
                _logger.LogWarning(ex,
                    "Resource not found for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status404NotFound,
                    "Not Found", ex.Message);
            }
            catch (ConflictException ex)
            {
                _logger.LogWarning(ex,
                    "Conflict for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status409Conflict,
                    "Conflict", ex.Message);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex,
                    "Unhandled exception for {UserId} at {RequestPath}. CorrelationId: {CorrelationId}",
                    userId, context.Request.Path, context.Items["CorrelationId"]);
                await WriteErrorResponse(context, StatusCodes.Status500InternalServerError,
                    "Internal Server Error", "An unexpected error occurred.");
            }
        }
    }

    private static async Task WriteErrorResponse(
        HttpContext context, int statusCode, string title, string detail)
    {
        context.Response.StatusCode = statusCode;
        context.Response.ContentType = "application/problem+json";
        var problem = new ProblemDetails
        {
            Status = statusCode,
            Title = title,
            Detail = detail
        };
        problem.Extensions["correlationId"] = context.Items["CorrelationId"]?.ToString();
        await context.Response.WriteAsJsonAsync(problem);
    }
}
```

**Ejemplo de respuesta de error (ProblemDetails)**:
```json
{
  "status": 404,
  "title": "Not Found",
  "detail": "User with id '3fa85f64-5717-4562-b3fc-2c963f66afa6' was not found.",
  "correlationId": "7f3e2a1b-9c4d-4e5f-8a6b-1d2e3f4a5b6c"
}
```

> **Content-Type canónico**: siempre `application/problem+json` (no `application/json`) para respuestas de error.

### IExceptionHandler — Alternativa No Canónica

ASP.NET Core 8 introduce `IExceptionHandler` como alternativa nativa. **No es el patrón canónico** en esta skill, pero puede usarse con el mismo contrato de `ProblemDetails` + `correlationId` si el equipo lo prefiere por consistencia con el framework.

```csharp
// Alternativa — mismo contrato de ProblemDetails, diferente mecanismo
public class GlobalExceptionHandler : IExceptionHandler
{
    public async ValueTask<bool> TryHandleAsync(
        HttpContext httpContext,
        Exception exception,
        CancellationToken cancellationToken)
    {
        var statusCode = exception is NotFoundException ? 404 : 500;
        httpContext.Response.StatusCode = statusCode;
        httpContext.Response.ContentType = "application/problem+json";
        var problem = new ProblemDetails { Status = statusCode, Title = "Error", Detail = exception.Message };
        problem.Extensions["correlationId"] = httpContext.Items["CorrelationId"]?.ToString();
        await httpContext.Response.WriteAsJsonAsync(problem, cancellationToken);
        return true;
    }
}
```

Registro en `Program.cs`:
```csharp
builder.Services.AddExceptionHandler<GlobalExceptionHandler>();
builder.Services.AddProblemDetails();
```

> **Nota**: Si usás `IExceptionHandler`, debés igualmente garantizar el mapeo canónico de excepciones y el enriquecimiento de `UserId`. El middleware custom es preferido porque centraliza ambas responsabilidades en un solo lugar.

### Campos Mínimos Obligatorios de Trazabilidad

Todo log de nivel `Warning` o superior DEBE contener estos 9 campos. Son la cadena de evidencia mínima para reconstruir un incidente.

| Campo | Fuente | Momento de enriquecimiento |
|---|---|---|
| `CorrelationId` | `HttpContext.Items["CorrelationId"]` | `CorrelationIdMiddleware` (inicio del pipeline) |
| `Application` | `IConfiguration["ApplicationName"]` | `CorrelationIdMiddleware` (inicio del pipeline) |
| `Environment` | `IWebHostEnvironment.EnvironmentName` | `CorrelationIdMiddleware` (inicio del pipeline) |
| `UserId` | `HttpContext.User` claims | `ExceptionHandlingMiddleware` (después de auth) |
| `RequestPath` | `HttpContext.Request.Path` | `UseSerilogRequestLogging` o middleware |
| `RequestMethod` | `HttpContext.Request.Method` | `UseSerilogRequestLogging` o middleware |
| `StatusCode` | `HttpContext.Response.StatusCode` | `UseSerilogRequestLogging` (al finalizar request) |
| `ExceptionType` | `exception.GetType().Name` | `ExceptionHandlingMiddleware` (en catch) |
| `Duration` | Calculado por Serilog request logging | `UseSerilogRequestLogging` (al finalizar request) |

> **Regla de oro**: `UserId` NUNCA se enriquece en `CorrelationIdMiddleware` porque `Authentication` no ejecutó aún. El orden del pipeline lo garantiza — ver [Pipeline Order en general/SKILL.md](../general/SKILL.md).

### Request Logging con Serilog

Log HTTP request/response al nivel del pipeline. Enriquecer con los campos de trazabilidad que Serilog no captura automáticamente:

```csharp
// Program.cs — después de UseMiddleware<CorrelationIdMiddleware>
app.UseSerilogRequestLogging(options =>
{
    options.EnrichDiagnosticContext = (diagnosticContext, httpContext) =>
    {
        diagnosticContext.Set("RequestHost", httpContext.Request.Host.Value);
        diagnosticContext.Set("UserAgent", httpContext.Request.Headers["User-Agent"].ToString());
        diagnosticContext.Set("CorrelationId", httpContext.Items["CorrelationId"]?.ToString());
        diagnosticContext.Set("UserId",
            httpContext.User?.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "anonymous");
    };
});
```

> `StatusCode` y `Duration` los captura Serilog automáticamente en request logging. No duplicarlos manualmente.

### Serilog Bootstrap — Configuración Canónica

```csharp
// Program.cs
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

> **Enrichers incluidos en v2.0**: `WithEnvironmentName()` del paquete `Serilog.Enrichers.Environment` para que `Environment` esté disponible como fallback estático (además del push dinámico del middleware).

### Datos Sensibles — Campos Prohibidos

Los siguientes datos NUNCA deben aparecer en logs. Sanitizar ANTES de loguear, nunca después:

| Categoría | Ejemplos de campos prohibidos |
|---|---|
| Credenciales | `password`, `passwordHash`, `pin`, `secret` |
| Tokens y claves | `jwtToken`, `accessToken`, `refreshToken`, `apiKey`, `bearerToken` |
| Datos financieros | `cardNumber`, `cvv`, `bankAccount`, `iban` |
| PII directo | `email` (completo), `phone`, `nationalId`, `ssn`, `dateOfBirth` |
| Infraestructura | `connectionString`, `databasePassword`, `smtpPassword` |

> **Alineación con security**: Esta lista complementa las reglas de `security/SKILL.md`. Si hay conflicto, la regla más restrictiva aplica.

---

## Code Examples

### Service-Level Logging

```csharp
public class UserService : IUserService
{
    private readonly IUserRepository _repository;
    private readonly ILogger<UserService> _logger;

    public UserService(IUserRepository repository, ILogger<UserService> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<UserDto> GetByIdAsync(Guid id)
    {
        _logger.LogInformation("Fetching user {UserId}", id);
        var user = await _repository.GetByIdAsync(id);
        if (user is null)
        {
            _logger.LogWarning("User {UserId} not found", id);
            throw new NotFoundException(nameof(User), id);
        }
        _logger.LogInformation("User {UserId} retrieved successfully", id);
        return user.ToDto();
    }
}
```

### Troubleshooting — Cadena de Evidencia

Flujo completo cuando un usuario reporta un error:

```
1. Usuario recibe error → obtiene correlationId del cuerpo ProblemDetails.Extensions["correlationId"]
2. Soporte busca en logs: WHERE CorrelationId = '{valor reportado}'
3. Reconstrucción:
   - Log de request (Serilog): Method, Path, StatusCode, Duration, UserId
   - Log de excepción (ExceptionHandlingMiddleware): ExceptionType, Message, StackTrace (interno)
   - Logs de servicio: Operación, parámetros (sin PII)
4. Con estos 3 tipos de log se puede reproducir el escenario completo
```

**Ejemplo de query en Seq / Kibana**:
```
CorrelationId = "7f3e2a1b-9c4d-4e5f-8a6b-1d2e3f4a5b6c"
```

---

## Configuración

### appsettings.json — Serilog

```json
{
  "ApplicationName": "MyApi",
  "Serilog": {
    "MinimumLevel": {
      "Default": "Information",
      "Override": {
        "Microsoft": "Warning",
        "System": "Warning"
      }
    },
    "WriteTo": [
      { "Name": "Console" },
      {
        "Name": "File",
        "Args": {
          "path": "logs/api-.log",
          "rollingInterval": "Day"
        }
      }
    ],
    "Enrich": ["FromLogContext", "WithMachineName", "WithEnvironmentName"]
  }
}
```

### Paquetes NuGet requeridos

```bash
dotnet add package Serilog.AspNetCore
dotnet add package Serilog.Sinks.Console
dotnet add package Serilog.Sinks.File
dotnet add package Serilog.Enrichers.Environment
```

---

## Frontera de Scope — Fase 2 (OUT OF SCOPE)

Las siguientes capacidades están **explícitamente fuera del alcance** de esta skill v2.0. Se abordarán en una iteración futura (Fase 2):

- **Multi-sink avanzado**: Bugsnag, Seq, Datadog, Application Insights como sinks primarios
- **OpenTelemetry**: Distributed tracing, métricas, exporters (OTLP, Jaeger, Zipkin)
- **Async sinks encolados**: Sinks con buffer y escritura asíncrona (Channel/Queue-based)
- **Distributed tracing**: Propagación de trace context (W3C TraceContext, B3) entre servicios
- **Health checks avanzados**: `/health` endpoints con checks de dependencias y formato detallado
- **Log sampling**: Reducción de volumen en producción con sampling configurable

> Si necesitás alguna de estas capacidades hoy, implementala de forma independiente sin modificar esta skill. La Fase 2 definirá el contrato canónico cuando corresponda.

---

## Cross-References

| Skill | Relación |
|---|---|
| [`responses`](../responses/SKILL.md) | Define el contrato general de `ProblemDetails`. La tabla de mapeo de excepciones **vive aquí** en `logging`. No duplicarla en `responses`. |
| [`general`](../general/SKILL.md) | Define el pipeline order. `CorrelationIdMiddleware` va en posición 0 — ver nota en `general`. |
| [`security`](../security/SKILL.md) | Define políticas de auth/JWT. `ForbiddenException` y prohibición de loguear tokens se alinean con esta skill. |

---

## Commands

```bash
# Install Serilog packages
dotnet add package Serilog.AspNetCore
dotnet add package Serilog.Sinks.Console
dotnet add package Serilog.Sinks.File
dotnet add package Serilog.Enrichers.Environment
```

---

## Resources

- **Standards**: See [../../../../rules-to-skills/Standardized_NET_Rules.md](../../../../rules-to-skills/Standardized_NET_Rules.md)
- **General conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **Responses**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Security**: See [../security/SKILL.md](../security/SKILL.md)
- **RFC 9457 — Problem Details**: https://www.rfc-editor.org/rfc/rfc9457

---

## Changelog

### v2.0 — 2026-03-25 ⚠️ Breaking Change

- **BREAKING**: Error responses migradas de `ResponseDTO<T>` a `ProblemDetails` (RFC 9457). `Content-Type` cambia a `application/problem+json`.
- **BREAKING**: `correlationId` ahora se expone en `ProblemDetails.Extensions["correlationId"]` (antes no se exponía al cliente).
- **Added**: Tabla canónica de mapeo de excepciones: `ValidationException`→400/Info, `ForbiddenException`→403/Warning, `NotFoundException`→404/Warning, `ConflictException`→409/Warning, `Exception`→500/Error.
- **Added**: `CorrelationIdMiddleware` v2.0 con enrichment de `Application` y `Environment`.
- **Added**: Enriquecimiento de `UserId` en `ExceptionHandlingMiddleware` (después de auth).
- **Added**: Tabla de 9 campos mínimos obligatorios de trazabilidad con fuente y momento de enriquecimiento.
- **Added**: Sección de campos prohibidos (datos sensibles / PII).
- **Added**: Sección de troubleshooting con cadena de evidencia completa.
- **Added**: Frontera explícita de Fase 2 (OUT OF SCOPE).
- **Changed**: `IExceptionHandler` degradado a alternativa no canónica; middleware custom es el patrón principal.
- **Changed**: Serilog bootstrap agrega `WithEnvironmentName()` enricher.
- **Meta**: Cambio bajo SDD formal (engram, OpenCode only). Decisiones D-33/D-34/D-35.
