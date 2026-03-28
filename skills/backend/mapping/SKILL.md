---
name: net8-apirest-mapping
description: >
  Object mapping standards for ASP.NET Core 8 REST APIs. Establishes extension methods
  as the default mapping strategy, defines layer ownership (Service owns all transformations),
  provides the canonical Request→Entity→DTO→Response flow, and documents AutoMapper usage
  limits. Trigger: When implementing DTO mappings, creating AutoMapper profiles, or building
  entity-to-DTO transformation logic in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Creating mappings between domain entities and DTOs
- Deciding whether to use AutoMapper or extension methods
- Building reusable manual mapping libraries
- Placing mapping logic correctly within the layered architecture
- Reviewing or refactoring existing mapping profiles

---

## Mapping Strategy Precedence (D-31, D-32)

> **Default: Extension methods. AutoMapper is the exception, not the rule.**

Apply the following precedence rules **in order**:

1. **Is this `Request → Entity` (write path)?** → Extension method — always. No exceptions.
2. **Is this `Entity → DTO` with any computed property, conditional, or transformation?** → Extension method.
3. **Is this `Entity → DTO` with 100% matching property names, same types, zero `ForMember`?** → AutoMapper permitted.
4. **In doubt?** → Extension method.

| Scenario | Approach |
|---|---|
| `Request → Entity` (CREATE / UPDATE) | **Extension method — mandatory** |
| `Entity → DTO` with computed props (`FullName`, `IsActive`) | **Extension method** |
| `Entity → DTO` 1:1, same names, same types, zero `ForMember` | AutoMapper permitted |
| `Entity → DTO` — high-frequency read-only (hot path) | Extension method or compiled mapping |
| Flattening nested objects | Extension method (explicit) |
| Collection mapping | `.ToDtoList()` extension method |

> **Why this precedence?** Extension methods are explicit, debuggable, testable as pure functions,
> and produce zero reflection overhead. AutoMapper's `ForMember` in `Request → Entity` hides
> critical assignments (`Id = Guid.NewGuid()`, `CreatedAt = DateTime.UtcNow`) that must be
> auditable. Manual is safer and more maintainable.

---

## Layer Ownership — Who Maps What

> **The Service is the sole owner of all mapping. No other layer transforms entities to DTOs or requests to entities.**

```
[Controller]  — Binds HTTP input to Request DTO. Returns IActionResult with DTO or ProblemDetails.
                Does NOT map. Does NOT transform.

[Service]     — OWNS ALL MAPPING:
                  request.ToEntity()  →  creates Entity from Request DTO
                  entity.ToDto()      →  creates DTO from Entity
                  entities.ToDtoList()→  creates DTO collection

[Repository]  — Returns Entity or null. Returns List<Entity> (never empty = null).
                Does NOT return DTOs. Does NOT map. Does NOT see Request DTOs.
```

| Layer | May Call | Must NOT Call |
|---|---|---|
| Controller | Service methods | Any `.ToDto()` / `.ToEntity()` / `IMapper` |
| Service | `.ToEntity()`, `.ToDto()`, `.ToDtoList()`, `IMapper.Map<>` (1:1 only) | Repository with DTOs |
| Repository | EF Core queries | Any mapping method, `IMapper`, DTOs |

---

## Canonical Data Flow

```
HTTP Request
    │
    ▼
[Controller]
  Model Binding → Request DTO
  FluentValidation → validated Request DTO
    │
    ▼  (passes validated Request DTO to Service)
[Service]
  request.ToEntity()        ← CREATE: builds new Entity
  manual field merge        ← UPDATE: merges Request fields into existing Entity
  repository.AddAsync(entity) / repository.UpdateAsync(entity)
  entity.ToDto()            ← always maps Entity → DTO before returning
    │
    ▼  (returns DTO to Controller)
[Controller]
  Ok(dto)                   ← 200
  CreatedAtAction(dto)      ← 201
  NoContent()               ← 204
    │
    ▼
HTTP Response
  2xx  → raw resource JSON (DTO body)
  4xx/5xx → ProblemDetails (see responses/SKILL.md)
```

---

## Critical Patterns

### Extension Method Pattern (Default — D-31)

Use C# extension methods in `Features/{Feature}/Mappings/`. This keeps transformations
explicit, testable as pure functions, and easy to debug.

```csharp
// Features/Users/Mappings/UserMappingExtensions.cs
public static class UserMappingExtensions
{
    // Entity → DTO  (READ — returns computed props, never null)
    public static UserDto ToDto(this User user) => new()
    {
        Id = user.Id,
        FullName = $"{user.FirstName} {user.LastName}",
        Email = user.Email,
        IsActive = user.Status == UserStatus.Active,
        RegisteredAt = user.CreatedAt.ToString("yyyy-MM-dd")
    };

    // Request → Entity  (CREATE — sets server-side fields here, auditable)
    public static User ToEntity(this CreateUserRequest request) => new()
    {
        Id = Guid.NewGuid(),
        FirstName = request.FirstName,
        LastName = request.LastName,
        Email = request.Email.ToLowerInvariant(),
        CreatedAt = DateTime.UtcNow,
        Status = UserStatus.Active
    };

    // Collection mapping  (pure LINQ, no allocation overhead)
    public static IEnumerable<UserDto> ToDtoList(this IEnumerable<User> users) =>
        users.Select(u => u.ToDto());
}
```

Usage in Service (canonical):
```csharp
// CREATE
public async Task<UserDto> CreateAsync(CreateUserRequest request)
{
    var entity = request.ToEntity();        // Service owns this
    await _repository.AddAsync(entity);
    return entity.ToDto();                  // Service owns this
}

// READ
public async Task<UserDto> GetByIdAsync(Guid id)
{
    var user = await _repository.GetByIdAsync(id);
    if (user is null) throw new NotFoundException(nameof(User), id);
    return user.ToDto();
}

// READ LIST
public async Task<List<UserDto>> GetAllActiveAsync()
{
    var users = await _repository.GetAllActiveAsync();
    return users.Count == 0 ? [] : users.ToDtoList().ToList();
}
```

### File Location

```
Features/Users/
  Mappings/
    UserProfile.cs               ← AutoMapper profile (only if 1:1 trivial — see precedence)
    UserMappingExtensions.cs     ← Extension methods (default — always present)
```

### AutoMapper — Permitted Use Only (1:1 Trivial)

AutoMapper is **only** acceptable when ALL of the following are true:
- `Entity → DTO` direction only (never `Request → Entity`)
- All property names match exactly
- All property types are identical or auto-convertible
- Zero `ForMember` calls needed

```csharp
// ✅ PERMITTED — pure 1:1, no ForMember
public class UserProfile : Profile
{
    public UserProfile()
    {
        CreateMap<User, UserSimpleDto>();   // No ForMember = OK
    }
}

// ❌ NOT PERMITTED — ForMember reveals hidden logic (use extension method instead)
public class UserProfile : Profile
{
    public UserProfile()
    {
        CreateMap<CreateUserRequest, User>()
            .ForMember(dest => dest.Id, opt => opt.MapFrom(_ => Guid.NewGuid()))
            .ForMember(dest => dest.CreatedAt, opt => opt.MapFrom(_ => DateTime.UtcNow));
    }
}
```

Register in `Program.cs` (if AutoMapper is used):
```csharp
builder.Services.AddAutoMapper(typeof(Program).Assembly);
```

> **Rule**: Never call `Mapper.Map` statically. Always inject `IMapper` via constructor.

```csharp
public class UserService : IUserService
{
    private readonly IMapper _mapper;                  // Only if AutoMapper is used
    private readonly IUserRepository _repository;

    public UserService(IMapper mapper, IUserRepository repository)
    {
        _mapper = mapper;
        _repository = repository;
    }

    public async Task<UserSimpleDto> GetSummaryAsync(Guid id)
    {
        var user = await _repository.GetByIdAsync(id);
        if (user is null) throw new NotFoundException(nameof(User), id);
        return _mapper.Map<UserSimpleDto>(user);  // 1:1 trivial — no ForMember
    }
}
```

### Generic Mapping Helper — `IEntityMapper<T>` (Optional Pattern)

For projects with many similar mappings, a shared generic abstraction is acceptable.

```csharp
// Shared/Mappings/IEntityMapper.cs
public interface IEntityMapper<TEntity, TDto>
{
    TDto ToDto(TEntity entity);
    TEntity ToEntity(TDto dto);
    IEnumerable<TDto> ToDtoList(IEnumerable<TEntity> entities);
}
```

```csharp
// Features/Users/Mappings/UserMapper.cs
public class UserMapper : IEntityMapper<User, UserDto>
{
    public UserDto ToDto(User entity) => new()
    {
        Id = entity.Id,
        FullName = $"{entity.FirstName} {entity.LastName}",
        Email = entity.Email
    };

    public User ToEntity(UserDto dto) => new()
    {
        Id = dto.Id,
        Email = dto.Email
    };

    public IEnumerable<UserDto> ToDtoList(IEnumerable<User> entities) =>
        entities.Select(ToDto);
}
```

> ⚠️ **`IEntityMapper<T>` — Singleton vs Scoped warning (D-34)**
>
> If the mapper implementation is **stateless** (no injected dependencies), `Singleton` is safe:
> ```csharp
> builder.Services.AddSingleton<IEntityMapper<User, UserDto>, UserMapper>();
> ```
>
> If the mapper **injects stateful/scoped dependencies** (`HttpContext`, `IConfiguration`,
> `DbContext`, `ICurrentUserService`, etc.), **you MUST register as `Scoped`**:
> ```csharp
> builder.Services.AddScoped<IEntityMapper<Order, OrderDto>, OrderMapper>();
> ```
>
> Registering a Scoped-dependent mapper as Singleton causes **captured stale references**
> — a subtle runtime bug that won't appear in tests but will in production under concurrent load.

---

## Anti-Patterns to Avoid

| Anti-pattern | Problem |
|---|---|
| Mapping `Request → Entity` with AutoMapper `ForMember` | Hides critical server-side assignments (`Id`, `CreatedAt`). Use extension method — explicit and auditable |
| Mapping inside the Repository | Repository returns pure entities. Mapping to DTOs here violates layer separation (see [`dataaccess`](../dataaccess/SKILL.md)) |
| Mapping inside the Controller | Controller is thin — no transformation logic. Service owns all mapping (see [`responses`](../responses/SKILL.md) Layer Communication Flow) |
| Calling `Mapper.Map` statically | Untestable; always inject `IMapper` via constructor |
| AutoMapper with `ForMember` for complex/conditional logic | Hard to debug; brittle; use extension method |
| Mapping entities to DTOs inside the EF Core query (non-`ProjectTo`) | Leaks DTO concerns into the data layer |
| `IEntityMapper<T>` registered as `Singleton` with stateful dependencies | Captures stale scoped references — runtime bug under concurrency (see D-34 warning above) |
| Null-returning extension methods | `.ToDto()` must never return `null`; throw `NotFoundException` in the Service before calling `.ToDto()` |

---

## Code Examples

### Collection Mapping with LINQ

```csharp
// Extension method — preferred for any DTO with transformations
var dtos = users
    .Where(u => u.IsActive)
    .Select(u => u.ToDto())
    .ToList();

// AutoMapper ProjectTo — acceptable for simple 1:1 read-only projections at DB level
var dtos = await _context.Users
    .Where(u => u.IsActive)
    .ProjectTo<UserSimpleDto>(_mapper.ConfigurationProvider)
    .AsNoTracking()
    .ToListAsync();
```

> **Note on `ProjectTo`**: translates mappings to SQL. Only use when the AutoMapper profile
> is 1:1 with no `ForMember` and all target properties exist in the database schema.

### UPDATE — Manual Field Merge (Service Pattern)

For `PUT`/`PATCH` operations, do NOT call `request.ToEntity()` on an existing entity —
that would create a new object and lose EF tracking. Merge fields manually in the Service:

```csharp
public async Task<UserDto> UpdateAsync(Guid id, UpdateUserRequest request)
{
    var user = await _repository.GetByIdAsync(id);
    if (user is null) throw new NotFoundException(nameof(User), id);

    // Manual merge — explicit, auditable, EF-tracking-safe
    user.FirstName = request.FirstName;
    user.LastName = request.LastName;
    user.Email = request.Email.ToLowerInvariant();
    user.UpdatedAt = DateTime.UtcNow;

    await _repository.UpdateAsync(user);
    return user.ToDto();
}
```

---

## Resources

- **Standards**: See [../../../../rules-to-skills/Standardized_NET_Rules.md](../../../../rules-to-skills/Standardized_NET_Rules.md)
- **General conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **Requests**: See [../requests/SKILL.md](../requests/SKILL.md) — Request DTOs que mapean a Entity via `.ToEntity()`
- **Responses**: See [../responses/SKILL.md](../responses/SKILL.md) — Service retorna DTO (via `.ToDto()`) que el Controller retorna como HTTP response
- **Data Access**: See [../dataaccess/SKILL.md](../dataaccess/SKILL.md) — Repository retorna Entity; NEVER DTOs; mapping en Service
- **Validations**: See [../validations/SKILL.md](../validations/SKILL.md) — Validation antes del mapping: validated Request → `.ToEntity()`
- **Unit Testing**: See [../testing-unit/SKILL.md](../testing-unit/SKILL.md) — Extension methods son pure functions, testeables sin mocks
- **Performance**: See [../performance/SKILL.md](../performance/SKILL.md) — Manual mapping > AutoMapper en hot paths

---

## Changelog

### v1.2 — 2026-03-28
- Removed: "Testing Guidance" section (already in `testing-unit` skill — not mapping-specific)
- Removed: Commands block (trivial dotnet add package — no project decision)

### v1.1 — 2026-03-25

**Mapping Strategy Precedence (D-31, D-32)**
- Extension methods promoted as **default**; AutoMapper limited to 1:1 trivial `Entity → DTO` without `ForMember`.
- Explicit 4-step precedence rule added (Request→Entity always manual; in doubt → manual).
- Updated Decision Tree into a precedence table with rationale.

**Layer Ownership (D-33)**
- New section declaring Service as the **sole owner** of all mapping.
- Explicit ownership table: what each layer may and must not call.

**Canonical Data Flow**
- New ASCII diagram: full `HTTP Request → Controller → Service → Repository → Service → Controller → HTTP Response` flow with mapping points annotated.

**Anti-Patterns Expanded**
- +4 new entries: mapping in Repository, mapping in Controller, AutoMapper `ForMember` for `Request → Entity`, `IEntityMapper<T>` Singleton with stateful dependencies.
- Null-returning `.ToDto()` added as anti-pattern.

**`IEntityMapper<T>` Warning (D-34)**
- Added explicit Singleton vs Scoped warning with code examples.
- Explains captured stale reference bug under concurrency.

**UPDATE Pattern**
- New section: manual field merge for `PUT`/`PATCH` (EF-tracking-safe; do NOT call `.ToEntity()` on existing entities).

**Testing Guidance**
- New section: test table for `.ToDto()`, `.ToEntity()`, `.ToDtoList()`, AutoMapper profiles, and edge cases.
- Full xUnit example with three test methods.

**Cross-References**
- Added: `requests`, `responses`, `dataaccess`, `validations`, `testing-unit`.
- AutoMapper command block now includes note about conditional use per D-31.

### v1.0 — Initial release

- Decision Tree AutoMapper vs Manual.
- File Location, AutoMapper Profile, Extension Method Pattern, Generic Mapper.
- Basic Anti-Patterns table.
