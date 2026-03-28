---
name: net8-apirest-responses
description: >
  Unified response standards for ASP.NET Core 8 REST APIs covering the dual response contract:
  `ResponseDTO<T>` for internal inter-layer communication (Repository → Service) and standard
  HTTP responses (raw resource JSON or `ProblemDetails`) for external controller → client responses.
  Trigger: When building API responses, defining return types for services and repositories,
  or standardizing how data and errors flow from the data layer to the client in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Definiendo return types de services y repositories
- Estandarizando responses de éxito y error
- Implementando exception handling middleware
- Revisando el flow de datos desde la capa de datos hasta el cliente

---

## Critical Patterns

### Dual Response Contract (D-25)

| Boundary | Contract | Quién lo usa |
|---|---|---|
| **Internal** (Repository → Service) | `ResponseDTO<T>` | Repositories, Services |
| **External** (Controller → Client) | Raw HTTP: resource JSON (2xx) o `ProblemDetails` (4xx/5xx) | Controllers |

> **Regla**: `ResponseDTO<T>` es un contrato interno de transporte. **Nunca** se serializa como HTTP response body para clientes externos.

### `ResponseDTO<T>` — Definición

```csharp
// Shared/Models/ResponseDTO.cs
public class ResponseDTO<T>
{
    public bool Success { get; set; }
    public string Message { get; set; } = string.Empty;
    public T? Data { get; set; }
}
```

### Layer Contract — Qué retorna cada capa

| Layer | Retorna | Nunca |
|---|---|---|
| Repository | Domain entities o `null` | DTOs, `ResponseDTO<T>` |
| Service | DTOs (después de mapear) o throws domain exception | Raw entities |
| Controller | `IActionResult` — raw resource JSON (2xx) o `ProblemDetails` (4xx/5xx) | El envelope `ResponseDTO<T>` al cliente |

```
Repository  →      Service      →   Controller  →   Client
  Entity     throws / DTO          IActionResult    HTTP 2xx (raw JSON)
                                                    HTTP 4xx/5xx (ProblemDetails)
```

### `ResponseDTO<T>` — Tres reglas de uso (interno)

| Escenario | Success | Message | Data |
|---|---|---|---|
| Operación exitosa con data | `true` | `""` o informativo | El objeto/lista resultado |
| Exitosa, sin resultados | `true` | Informativo (ej. "No records found.") | Empty value, **nunca null** |
| Error / fallo | `false` | Descripción clara del error | `null` |

Empty value por tipo cuando success pero sin data: `[]` para listas, `string.Empty` para strings, `0` para numéricos, `new T()` si tiene constructor sin parámetros.

### Paged Response

```csharp
// Shared/Models/PagedResult.cs
public class PagedResult<T>
{
    public List<T> Items { get; set; } = new();
    public int TotalCount { get; set; }
    public int Page { get; set; }
    public int PageSize { get; set; }
    public int TotalPages => (int)Math.Ceiling((double)TotalCount / PageSize);
}
```

---

## HTTP Status Code Reference

| Code | Cuando usar |
|---|---|
| `200 OK` | GET, PUT retornando resource |
| `201 Created` | POST — incluir `Location` header |
| `204 No Content` | DELETE, PUT sin body |
| `400 Bad Request` | Errores de validación |
| `401 Unauthorized` | **Pipeline/auth — NO es domain exception** |
| `403 Forbidden` | Autenticado pero sin permiso |
| `404 Not Found` | ID no encontrado |
| `409 Conflict` | Duplicado, modificación concurrente |
| `429 Too Many Requests` | Rate limiter activado |
| `500 Internal Server Error` | Exception no manejada |

---

## Exception Handling — HTTP Mapping (Single Source of Truth)

> **Ownership**: El contrato exception → HTTP vive aquí (`responses`). La implementación del middleware, logging y correlationId viven en [`logging`](../logging/SKILL.md).

### Tabla Canónica

| Exception | HTTP Status | Log Level | Client Message |
|---|---|---|---|
| `NotFoundException` | `404 Not Found` | Warning | `ex.Message` (safe) |
| `ForbiddenException` | `403 Forbidden` | Warning | `ex.Message` (safe) |
| `ConflictException` | `409 Conflict` | Warning | `ex.Message` (safe) |
| `ValidationException` (manual) | `400 Bad Request` | Information | First error message |
| `Exception` (unhandled) | `500 Internal Server Error` | Error | Generic — **nunca exponer stack trace** |

> **`401 Unauthorized` — Pipeline/Auth ownership (D-27)**  
> `401` NO es una domain exception. Lo resuelve el pipeline de auth (JWT bearer, `[Authorize]`).  
> **No crear `UnauthorizedException`.**

### Exception Location Reference

| Exception | Definida en |
|---|---|
| `NotFoundException` | `general` skill — `Shared/Exceptions/NotFoundException.cs` |
| `ForbiddenException` | `security` skill — `Shared/Exceptions/ForbiddenException.cs` |
| `ConflictException` | `Shared/Exceptions/ConflictException.cs` |

### Uniform 400 — `InvalidModelStateResponseFactory`

```csharp
builder.Services.AddControllers()
    .ConfigureApiBehaviorOptions(options =>
    {
        options.InvalidModelStateResponseFactory = context =>
        {
            var firstError = context.ModelState.Values
                .SelectMany(v => v.Errors)
                .Select(e => e.ErrorMessage)
                .FirstOrDefault() ?? "Validation failed.";

            return new BadRequestObjectResult(new ProblemDetails
            {
                Status = StatusCodes.Status400BadRequest,
                Title = "Bad Request",
                Detail = firstError
            });
        };
    });
```

### Uniform 429 — Rate Limiter `OnRejected`

```csharp
options.OnRejected = async (context, ct) =>
{
    context.HttpContext.Response.StatusCode = StatusCodes.Status429TooManyRequests;
    await context.HttpContext.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status429TooManyRequests,
        Title = "Too Many Requests",
        Detail = "Rate limit exceeded. Please try again later."
    }, ct);
};
```

---

## Anti-Patterns

| Anti-pattern | Problema |
|---|---|
| `ResponseDTO<T>` en el HTTP response body | Viola D-25 — el cliente recibe un envelope interno |
| Retornar `null` desde service en "not found" | El controller debe null-check en todos lados; mejor throw |
| HTTP 200 para errores | Engaña a los clientes — siempre usar status codes correctos |
| `ResponseDTO<T>` en repositories | Repositories son data layer, no API layer |
| Crear `UnauthorizedException` como domain exception | `401` es del pipeline — no del dominio |
| No incluir `Location` header en 201 Created | Viola REST |

---

## Resources

- **General / NotFoundException**: See [../general/SKILL.md](../general/SKILL.md)
- **Requests**: See [../requests/SKILL.md](../requests/SKILL.md)
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md)
- **Logging / ExceptionHandlingMiddleware**: See [../logging/SKILL.md](../logging/SKILL.md)
- **Validations**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Security / ForbiddenException**: See [../security/SKILL.md](../security/SKILL.md)

---

## Changelog

### v1.2 — 2026-03-28
- **Removed**: Full CRUD controller example (ya está en `general/SKILL.md`)
- **Removed**: Explicación verbosa del "why" del dual contract
- **Removed**: Snippets completos del middleware (viven en `logging/SKILL.md`)
- **Removed**: Repository y Service return pattern examples (consolidados en la tabla de layer contract)
- **Kept**: Dual contract table, ResponseDTO<T> definition, layer contract table, paged result, HTTP status table, exception mapping como single source of truth, uniform 400/429

### v1.1 — 2026-03-24
- Dual Response Contract D-25/D-26/D-27; ProblemDetails; exception ownership clarificado
