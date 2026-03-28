---
name: net8-apirest-requests
description: >
  Best practices for handling HTTP requests in ASP.NET Core 8 REST APIs covering model binding, DTO design, input binding sources, and request pipeline conventions.
  Trigger: When designing request DTOs, configuring model binding, handling query parameters, or processing incoming HTTP payloads in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Diseñando request DTOs para endpoints
- Manejando file uploads o requests multipart
- Revisando action signatures de controllers

---

## Critical Patterns

### Request DTO Naming

Siempre sufijo `Request`. Nunca reusar domain entities como request bodies:

```
CreateUserRequest
UpdateUserRequest
SearchUsersRequest
```

### Pagination — Cap obligatorio en `PageSize`

```csharp
public class PaginationRequest
{
    public int Page { get; set; } = 1;
    public int PageSize { get; set; } = 20;

    // Cap obligatorio — previene pageSize=9999
    public int ValidatedPageSize => Math.Min(PageSize, 100);
}

public class SearchUsersRequest : PaginationRequest
{
    public string? Name { get; set; }
    public bool? IsActive { get; set; }
}
```

### Route Constraints — GUID en recursos públicos

```csharp
// Rechaza IDs inválidos antes del service — previene BOLA enumeration
[HttpGet("{id:guid}")]
```

Nunca `{id:int}` en endpoints públicos — ver [`../security/SKILL.md`](../security/SKILL.md).

### Binding Sources — Siempre explícito

No depender del inference implícito de ASP.NET Core. Anotar siempre:

```csharp
[HttpGet("{id:guid}")]
public async Task<IActionResult> GetById(
    [FromRoute] Guid id,
    [FromHeader(Name = "X-Tenant-ID")] string? tenantId) { ... }

[HttpGet]
public async Task<IActionResult> Search([FromQuery] SearchUsersRequest request) { ... }

[HttpPost]
public async Task<IActionResult> Create([FromBody] CreateUserRequest request) { ... }
```

### Null Safety en Request DTOs

Required fields non-nullable con default. Optional fields nullable:

```csharp
public class CreateUserRequest
{
    public string FirstName { get; set; } = string.Empty;   // required
    public string Email { get; set; } = string.Empty;       // required
    public string? PhoneNumber { get; set; }                // optional
}
```

### File Uploads

El pattern de validación de file uploads (D-30) vive en [`../validations/SKILL.md`](../validations/SKILL.md).

> **Nota crítica**: `[Consumes("multipart/form-data")]` a nivel de action **overrides** el `[Consumes("application/json")]` del controller. Siempre anotarlo en actions de file upload — sin él ASP.NET Core rechaza el request con 415.

---

## Anti-Patterns

| Anti-pattern | Problema |
|---|---|
| Domain entity como request body | Acopla API al data model; riesgo de seguridad |
| Binding source sin anotar | El inference implícito puede fallar de formas no obvias |
| Pagination sin cap de `PageSize` | `pageSize=9999` es un vector de DoS |
| `int` IDs en endpoints públicos | BOLA enumeration trivial |
| `ResponseDTO<T>` como HTTP response body | Viola D-25/D-26 — retornar raw resource o `ProblemDetails` |
| Leer `IFormFile` stream en validator | Bloquea el stream — ver D-30 en `validations` |

---

## Resources

- **General structure**: See [../general/SKILL.md](../general/SKILL.md)
- **Mapping**: See [../mapping/SKILL.md](../mapping/SKILL.md)
- **Validations (file upload D-30)**: See [../validations/SKILL.md](../validations/SKILL.md)
- **Responses**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Security (BOLA, GUID IDs)**: See [../security/SKILL.md](../security/SKILL.md)

---

## Changelog

### v1.2 — 2026-03-28
- **Removed**: Tutorial completo de binding sources con tabla y ejemplos (el agente conoce `[FromBody]`, `[FromQuery]`, etc.)
- **Removed**: Content negotiation `[Produces]`/`[Consumes]` boilerplate — estándar del framework
- **Removed**: File upload section extensa — consolidada en `validations/SKILL.md` (D-30), solo referencia aquí
- **Kept**: Naming convention, pagination con ValidatedPageSize cap, GUID route constraint, binding explícito, null safety, anti-patterns

### v1.1 — 2026-03-24
- File uploads rewritten con D-30 pattern; anti-patterns expandidos
