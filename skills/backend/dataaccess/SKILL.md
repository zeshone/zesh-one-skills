---
name: net8-apirest-dataaccess
description: >
  Entity Framework Core data access standards for ASP.NET Core 8 REST APIs. Covers DbContext setup, entity configuration with Fluent API, repository pattern, audit automation via interceptors, soft delete with global query filters, raw SQL security, transactions, and all anti-patterns confirmed in the real codebase audit.
  Trigger: When working with EF Core, repositories, DbContext, raw SQL stored procedures, entity models, migrations, transactions, or audit/soft-delete patterns in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.8"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Creating or modifying a DbContext class
- Implementing a new repository (interface + class)
- Writing EF Core queries (LINQ, raw SQL, stored procedures)
- Configuring entity relationships, audit fields, or soft delete
- Automating audit stamps via `ISaveChangesInterceptor`
- Managing database transactions across multiple contexts
- Reviewing data access code for anti-patterns

---

## Critical Patterns

### [DEFICIENCY FIX] DbContext Registration — Always AddDbContextPool

**Never** use `AddDbContext`. Standard is `AddDbContextPool` with `poolSize: 128`.  
**Confirmed broken in**: Program.cs registration code — must be corrected.

```csharp
// CORRECT — mandatory standard
builder.Services.AddDbContextPool<AppDbContext>(options =>
{
    var connectionString = builder.Configuration.GetConnectionString("DefaultConnection")
        ?? throw new InvalidOperationException("Connection string not configured.");
    var connBuilder = new SqlConnectionStringBuilder(connectionString)
    {
        MultipleActiveResultSets = true,
        TrustServerCertificate = true
    };
    options.UseSqlServer(connBuilder.ConnectionString);
}, poolSize: 128);

// WRONG — no pool, no reuse, higher allocation pressure
builder.Services.AddDbContext<AppDbContext>(options => ...);
```

Multiple contexts — all must use pool in `Program.cs`:
```csharp
builder.Services.AddDbContextPool<AppDbContext>(..., poolSize: 128);
builder.Services.AddDbContextPool<InventoryDbContext>(..., poolSize: 128);
builder.Services.AddDbContextPool<ReportingDbContext>(..., poolSize: 128);
```

### Two-Level Pool — ADO.NET + EF Core

`AddDbContextPool` reuses `DbContext` **instances** (avoids object allocation). The ADO.NET connection pool reuses **TCP connections** to SQL Server. Both levels must be configured for high-concurrency APIs.

```csharp
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection")
    ?? throw new InvalidOperationException("Connection string not configured.");

// Level 1 — ADO.NET connection pool (TCP connections to SQL Server)
var connBuilder = new SqlConnectionStringBuilder(connectionString)
{
    MaxPoolSize = 200,          // max concurrent connections to SQL Server
    MinPoolSize = 10,           // connections kept alive between requests
    ConnectTimeout = 30,        // seconds to wait before failing
    Pooling = true,
    MultipleActiveResultSets = true,
    TrustServerCertificate = true
};

// Level 2 — EF Core context pool (DbContext instances, not connections)
builder.Services.AddDbContextPool<AppDbContext>(options =>
    options.UseSqlServer(connBuilder.ConnectionString),
    poolSize: 128);
```

> **Non-obvious**: These are two independent pools. `MaxPoolSize` (ADO.NET) controls how many SQL connections exist. `poolSize` (EF) controls how many `DbContext` objects are reused. A `DbContext` from the pool acquires a connection from the ADO.NET pool on demand — they are not 1:1.

### Connection Strings — Standard ASP.NET Core

Always retrieve connection strings using standard ASP.NET Core configuration:

```csharp
// CORRECT — standard IConfiguration pattern
var cs = builder.Configuration.GetConnectionString("DefaultConnection");

// CORRECT — complex settings via IOptions<T>
builder.Services.Configure<DatabaseOptions>(
    builder.Configuration.GetSection("DatabaseOptions"));

// WRONG — proprietary singleton or hardcoded paths
var cs = SomeConfigurations.Instance.Config.ConnectionStrings.SomeKey;
```

> For complex database settings (timeouts, retry policies), use `IOptions<T>` bound to a dedicated config section in `appsettings.json`.

---

### [DEFICIENCY FIX] AsNoTracking — Mandatory on Every Read Query

**Confirmed missing** in: read-heavy repository queries and others.  
`AsNoTracking` must be applied to **every** query that does not modify entities.  
EF Core will not set up change tracking info — less memory, faster materialization.

```csharp
// CORRECT — single entity read
public async Task<User?> GetByIdAsync(Guid id, CancellationToken ct = default) =>
    await _context.Users
        .AsNoTracking()
        .FirstOrDefaultAsync(u => u.Id == id, ct);

// CORRECT — list read
public async Task<List<User>> GetAllActiveAsync(CancellationToken ct = default) =>
    await _context.Users
        .Where(u => u.IsActive)
        .AsNoTracking()
        .ToListAsync(ct);

// CORRECT — identity resolution for complex includes (avoids duplicate entities)
var orders = await _context.Orders
    .Include(o => o.Items)
        .ThenInclude(i => i.Product)
    .AsNoTrackingWithIdentityResolution()
    .ToListAsync(ct);

// WRONG — tracking enabled with no intent to save
var users = await _context.Users.ToListAsync();
```

> **Exception**: only omit `AsNoTracking()` when you need to modify and call `SaveChangesAsync` on the same entity within the same scope.

**Performance option for read-heavy repositories** — set default at context level:
```csharp
// In constructor or OnConfiguring — applies to all queries in this instance
_context.ChangeTracker.QueryTrackingBehavior = QueryTrackingBehavior.NoTracking;
```

### Eager Loading — Always Explicit (Lazy Loading Disabled)

Lazy loading is **off** and must stay off. Navigation properties are `virtual` in some models but proxies are not enabled — they will always be `null` unless explicitly loaded.

```csharp
// CORRECT — nested (ThenInclude — introduce going forward)
var orders = await _context.Orders
    .Include(o => o.Customer)
        .ThenInclude(c => c.Contacts)
    .AsNoTracking()
    .ToListAsync(ct);

// WRONG — navigation property accessed without Include = always null
var orders = await _context.Orders.AsNoTracking().ToListAsync(ct);
var name = orders[0].Customer.Name; // NullReferenceException
```

---

### [DEFICIENCY FIX] Entity Configuration — IEntityTypeConfiguration<T>

**Current state**: large `OnModelCreating` blocks in existing contexts + Data Annotations mixed.  
**Going forward**: one `IEntityTypeConfiguration<T>` class per entity, auto-discovered via assembly scanning.

```csharp
// Database/Context/AppDbContext.cs — new projects only
protected override void OnModelCreating(ModelBuilder modelBuilder)
{
    modelBuilder.ApplyConfigurationsFromAssembly(typeof(AppDbContext).Assembly);
}
```

```csharp
// Features/Users/Models/UserConfiguration.cs
public class UserConfiguration : IEntityTypeConfiguration<User>
{
    public void Configure(EntityTypeBuilder<User> builder)
    {
        builder.ToTable("Users");
        builder.HasKey(u => u.Id);
        builder.Property(u => u.Email).HasMaxLength(200).IsRequired();
        builder.HasIndex(u => u.Email).IsUnique();
        builder.Property(u => u.CreatedAt).HasDefaultValueSql("GETUTCDATE()");
        // UpdatedAt managed by AuditInterceptor — do NOT add ValueGeneratedOnUpdate()

        // Relationship config — explicit cascade behavior
        builder.HasMany(u => u.Orders)
            .WithOne(o => o.User)
            .HasForeignKey(o => o.UserId)
            .OnDelete(DeleteBehavior.ClientSetNull); // existing project standard
    }
}
```

> **Existing legacy contexts**: If your project has established DbContexts with large `OnModelCreating` blocks, keep the current approach to avoid unnecessary migration churn.  
> Apply `IEntityTypeConfiguration<T>` **only on new entities and new projects**.

### [DEFICIENCY FIX] Audit Fields — Standardized + Automated

**Confirmed problem**: inconsistent audit field names across projects.  
**Standard for all new entities**:

```csharp
// Shared/Models/BaseEntity.cs
public abstract class BaseEntity
{
    public Guid Id { get; set; } = Guid.NewGuid();
    public DateTime CreatedAt { get; set; }
    public DateTime? UpdatedAt { get; set; }
}
```

**Automate audit stamps via ISaveChangesInterceptor** (eliminates manual assignment):

```csharp
// Shared/Interceptors/AuditInterceptor.cs
public class AuditInterceptor : SaveChangesInterceptor
{
    public override ValueTask<InterceptionResult<int>> SavingChangesAsync(
        DbContextEventData eventData,
        InterceptionResult<int> result,
        CancellationToken ct = default)
    {
        if (eventData.Context is null) return base.SavingChangesAsync(eventData, result, ct);
        var now = DateTime.UtcNow;
        foreach (var entry in eventData.Context.ChangeTracker.Entries<BaseEntity>())
        {
            if (entry.State == EntityState.Added)
                entry.Entity.CreatedAt = now;
            if (entry.State == EntityState.Modified)
                entry.Entity.UpdatedAt = now;
        }
        return base.SavingChangesAsync(eventData, result, ct);
    }
}
```

Register in `Program.cs`:
```csharp
builder.Services.AddSingleton<AuditInterceptor>();

builder.Services.AddDbContextPool<AppDbContext>((sp, options) =>
{
    options.UseSqlServer(connectionString)
           .AddInterceptors(sp.GetRequiredService<AuditInterceptor>());
}, poolSize: 128);
```

---

### [NEW — NOT IN CODEBASE] Soft Delete with Global Query Filters

The project has **no soft delete pattern**. When needed, implement it this way:

> **Modified version of `BaseEntity` when soft delete is needed**: add `IsDeleted` to the existing class shown above. This is not a second competing base class.

Add `IsDeleted` to the entity:
```csharp
public abstract class BaseEntity
{
    public Guid Id { get; set; } = Guid.NewGuid();
    public bool IsDeleted { get; set; } = false;
    public DateTime CreatedAt { get; set; }
    public DateTime? UpdatedAt { get; set; }
}
```

Register `HasQueryFilter` — automatically excludes soft-deleted records from ALL queries:
```csharp
// In IEntityTypeConfiguration<T> or OnModelCreating
builder.HasQueryFilter(e => !e.IsDeleted);
```

Override `SaveChangesAsync` in `DbContext` to intercept physical deletes:
```csharp
public override async Task<int> SaveChangesAsync(CancellationToken ct = default)
{
    foreach (var entry in ChangeTracker.Entries<BaseEntity>())
    {
        if (entry.State == EntityState.Deleted)
        {
            entry.State = EntityState.Modified;
            entry.Entity.IsDeleted = true;
        }
    }
    return await base.SaveChangesAsync(ct);
}
```

To query deleted records explicitly:
```csharp
var all = await _context.Users.IgnoreQueryFilters().ToListAsync();
```

---

### [DEFICIENCY FIX] Repository Pattern — Entities Only, Never DTOs

**Confirmed anti-pattern**: returning DTOs directly from repository classes.  
New repositories must follow this contract:

```
Features/{Feature}/
  Repositories/
    I{Feature}Repository.cs   ← interface
    {Feature}Repository.cs    ← implementation
```

```csharp
public interface IUserRepository
{
    Task<User?> GetByIdAsync(Guid id, CancellationToken ct = default);
    Task<List<User>> GetAllActiveAsync(CancellationToken ct = default);
    Task<bool> ExistsByEmailAsync(string email, CancellationToken ct = default);
    Task AddAsync(User user, CancellationToken ct = default);
    Task UpdateAsync(User user, CancellationToken ct = default);
    /// <summary>
    /// Deletes the entity with the given <paramref name="id"/>.
    /// CONTRACT: The caller (service layer) MUST call <see cref="GetByIdAsync"/> first
    /// and throw <see cref="NotFoundException"/> if the entity is not found.
    /// Passing a stub entity to <c>Remove</c> without first verifying existence may silently delete the wrong row or cause FK violations.
    /// The service contract (calling <see cref="GetByIdAsync"/> first and throwing <see cref="NotFoundException"/> if absent) is the only safety guard.
    /// </summary>
    Task DeleteAsync(Guid id, CancellationToken ct = default);
}
```

**Layer contract — never violate**:

| Layer | Returns | Never |
|---|---|---|
| Repository | Domain entities, `null`, empty `List<T>` | DTOs, `ResponseDTO<T>` |
| Service | `Result<T>` wrapping DTO (preferred) / bare DTO (legacy `ResponseDTO<T>` path) | Raw entities |
| Controller | `result.ToHttpResponse()` / raw resource JSON (2xx) / `ProblemDetails` (4xx–5xx) | Raw entities |

---

### [DEFICIENCY FIX] Raw SQL & Stored Procedures — SQL Injection Prevention

**Heavy SP usage confirmed** across legacy codebases. Risk of SQL injection via string concat exists.  
EF Core 10 introduced an analyzer that warns on concatenation in `FromSqlRaw` — adopt safe patterns now:

| Method | When to use | Safe? |
|---|---|---|
| `FromSqlInterpolated($"EXEC sp {param}")` | Entity result from SP | ✅ Auto-parameterized |
| `FromSqlRaw("EXEC sp @p", sqlParam)` | Entity result, legacy | ✅ With `SqlParameter` |
| `Database.SqlQueryRaw<T>("...", sqlParam)` | Non-entity result | ✅ With `SqlParameter` |
| `ExecuteSqlRaw("...", sqlParam)` | CUD via SP | ✅ With `SqlParameter` |
| `FromSqlRaw("EXEC sp '" + id + "'")` | **NEVER** | ❌ SQL injection |
| `FromSqlRaw($"EXEC sp {id}")` with `string` type | **NEVER** | ❌ SQL injection |

```csharp
// CORRECT — non-entity result set
var rows = await _context.Database
    .SqlQueryRaw<UserSummaryDto>("EXEC sp_GetSummary @Id", new SqlParameter("@Id", id))
    .ToListAsync(ct);
```

> **Important for FromSqlRaw + SPs**: add `.AsEnumerable()` immediately after when calling SPs to prevent EF from attempting to compose the result into a subquery, which generates invalid SQL.

### Keyless Entities — SP Result Sets Pattern

Use keyless entities to map stored procedure result sets as typed classes. This avoids raw `DataTable` usage and keeps the result strongly typed without creating a real domain entity.

```csharp
// Features/Reports/Models/PolicySummaryResult.cs
// This is NOT a domain entity — it is a SP result projection
public class PolicySummaryResult
{
    public Guid PolicyId { get; set; }
    public string HolderName { get; set; } = string.Empty;
    public decimal Premium { get; set; }
    public DateTime ExpiresAt { get; set; }
}
```

**Register in `OnModelCreating` or `IEntityTypeConfiguration<T>`:**
```csharp
// HasNoKey() — no PK, EF won't try to track or update these
// ToView(null) — CRITICAL: tells EF there is no backing table or view.
//                Without this, migrations try to CREATE a table for it.
modelBuilder.Entity<PolicySummaryResult>().HasNoKey().ToView(null);
```

**Query:**
```csharp
var results = await _context.Set<PolicySummaryResult>()
    .FromSqlRaw("EXEC sp_GetPolicySummary @AgentId", new SqlParameter("@AgentId", agentId))
    .AsNoTracking()
    .ToListAsync(ct);
```

> **Why `.ToView(null)`**: without it, EF's migration tooling assumes the entity needs a table and generates a `CREATE TABLE` migration on the next `dotnet ef migrations add`. `.ToView(null)` explicitly marks it as unmapped — migrations ignore it.

> **File location**: Place SP result classes in `Features/{Feature}/Models/` with a `Result` suffix (e.g., `PolicySummaryResult.cs`) to distinguish them from domain entities. Never suffix them `Dto` — they are not DTOs; they don't cross service layer boundaries.

---

### [DEFICIENCY FIX] Transactions — Single SaveChanges or Explicit Transaction

**Confirmed anti-pattern**: calling `SaveChangesAsync` multiple times without a wrapping transaction. If the second save fails, the first is already committed — inconsistent state.

```csharp
// CORRECT — single save for atomic operation
public async Task CreateOrderAsync(Order order, List<OrderItem> items, CancellationToken ct = default)
{
    _context.Orders.Add(order);
    _context.OrderItems.AddRange(items);
    await _context.SaveChangesAsync(ct); // one atomic operation
}

// CORRECT — multiple saves wrapped in explicit transaction
public async Task CreateOrderWithNotificationAsync(Order order, Notification notif, CancellationToken ct = default)
{
    await using var transaction = await _context.Database.BeginTransactionAsync(ct);
    try
    {
        _context.Orders.Add(order);
        await _context.SaveChangesAsync(ct);
        _context.Notifications.Add(notif);
        await _context.SaveChangesAsync(ct);
        await transaction.CommitAsync(ct);
    }
    catch
    {
        await transaction.RollbackAsync(ct);
        throw;
    }
}

// WRONG — multiple uncommitted saves without transaction
_context.Orders.Add(order);
await _context.SaveChangesAsync(); // committed, no rollback if next fails
_context.OrderItems.AddRange(items);
await _context.SaveChangesAsync(); // if this fails, order is orphaned
```

### `CreateExecutionStrategy()` — Required When Retry Policies Are Active

> ⚠️ **Gotcha**: If EF Core is configured with a retry-on-failure execution strategy (e.g., `EnableRetryOnFailure()` for SQL Server), opening a manual transaction directly with `BeginTransactionAsync` throws `InvalidOperationException`:
> *"The configured execution strategy does not support user-initiated transactions."*

Wrap the entire transaction block inside `CreateExecutionStrategy().ExecuteAsync()`:

```csharp
public async Task CreateOrderWithRetryAsync(Order order, Notification notif, CancellationToken ct = default)
{
    // Required when UseSqlServer(..., o => o.EnableRetryOnFailure()) is configured
    var strategy = _context.Database.CreateExecutionStrategy();

    await strategy.ExecuteAsync(async () =>
    {
        await using var transaction = await _context.Database.BeginTransactionAsync(ct);
        try
        {
            _context.Orders.Add(order);
            await _context.SaveChangesAsync(ct);
            _context.Notifications.Add(notif);
            await _context.SaveChangesAsync(ct);
            await transaction.CommitAsync(ct);
        }
        catch
        {
            await transaction.RollbackAsync(ct);
            throw;
        }
    });
}
```

> **Why it's needed**: EF's retry strategy needs to control the entire operation so it can retry from the start on transient failure. A manually opened transaction breaks that contract — `CreateExecutionStrategy` restores it by giving EF a re-entrant execution scope.

> ⚠️ **WARNING — Multi-context transactions cannot be made atomic without a distributed coordinator.**
>
> Two independent `IDbContextTransaction` instances (one per `DbContext`) are two separate SQL transactions on the database. Committing `primaryTx` then `secondaryTx` sequentially means:
 > - If `secondaryTx.CommitAsync` throws, `primaryTx` is **already committed** — there is no way to undo it — calling `RollbackAsync` on an already-committed transaction throws `InvalidOperationException`.
> - You cannot guarantee both commits succeed or both fail without MS DTC (Distributed Transaction Coordinator) or a saga/outbox pattern.
>
> **Do NOT write code that calls `CommitAsync` on two independent transactions and claims it is atomic — it is not.**
>
> **Recommended approach for cross-context consistency: Outbox Pattern.**
> Instead of committing to two databases atomically, write the cross-context event into an `OutboxMessages` table in the **same** transaction as the primary write. A background worker (e.g., Hangfire, a hosted service) reads the outbox and delivers to the secondary context with at-least-once semantics. This provides eventual consistency without distributed locks.

---

### [NEW — NOT IN CODEBASE] Compiled Queries for Hot Paths

Not used in the codebase. Introduce for frequently called queries (auth checks, session lookups):

```csharp
// Define once — skips LINQ translation on every call
private static readonly Func<AppDbContext, Guid, Task<User?>> _getUserById =
    EF.CompileAsyncQuery((AppDbContext ctx, Guid id) =>
        ctx.Users.AsNoTracking().FirstOrDefault(u => u.Id == id));

// Use in repository
public Task<User?> GetByIdAsync(Guid id) => _getUserById(_context, id);
```

---

### [DEFICIENCY FIX] Migrations — Descriptive Naming

**Confirmed problem**: migration `AutoMigration_20240101_000000` in production — zero meaning.

```bash
# CORRECT — descriptive, intent-driven names
dotnet ef migrations add Add_Users_Table
dotnet ef migrations add Add_IsDeleted_To_Orders
dotnet ef migrations add Fix_Email_MaxLength_Constraint

# WRONG — auto-generated, meaningless in history
dotnet ef migrations add AutoMigration_20240101_000000
```

Migration folder: `src/YourProject.Infrastructure/Migrations/`

---

## Anti-Patterns — Confirmed in Codebase

| Anti-pattern | Where to look | Required fix |
|---|---|---|
| `AddDbContext` without pool | Program.cs registration code | `AddDbContextPool` with `poolSize: 128` |
| Missing `AsNoTracking` on reads | Read-heavy repository queries | Add to every read — no exceptions |
| Repository returning DTOs | Shared/generic repository classes | Return entities only; map in service |
| Multiple `SaveChangesAsync` without transaction | Feature repository with multiple saves | Single save or `BeginTransactionAsync` |
| Inconsistent audit field names | Legacy DbContext classes | `CreatedAt` / `UpdatedAt` on new entities |
| No automated audit stamp | Entire codebase | `ISaveChangesInterceptor` (AuditInterceptor) |
| No soft delete standard | Entire codebase | `HasQueryFilter(e => !e.IsDeleted)` when needed |
| No compiled queries | Entire codebase | `EF.CompileAsyncQuery` on hot read paths |
| Auto-generated migration names | Infrastructure/Migrations project | Descriptive PascalCase names |
| Large `OnModelCreating` | Legacy DbContext classes | Keep existing; `IEntityTypeConfiguration<T>` on new |
| `FromSqlRaw` with string concat potential | Raw SQL / SP execution code | Always use `SqlParameter` or `FromSqlInterpolated` |
| Virtual nav props without lazy loading enabled | Multiple models | Always use explicit `.Include()` — never rely on nav props auto-loading |
| Calling `DeleteAsync` without prior existence check | Service layer delete flows | Service MUST call `GetByIdAsync` first; throw `NotFoundException` if missing; repository assumes entity exists |

*Anti-patterns confirmed in production legacy codebases.*

---

## Changelog

### v2.8 — 2026-04-09
- **Fixed (Round 3)**: Layer contract table updated — controller row no longer mentions `ResponseDTO<T>` as the return type. Canonical controller output is now `result.ToHttpResponse()` or raw HTTP (`2xx` JSON / `ProblemDetails` for `4xx–5xx`), aligned with `responses/SKILL.md`.

### v2.7 — 2026-04-09
- **Added**: Two-Level Pool section — documents the distinction between ADO.NET connection pool (`MaxPoolSize`, `MinPoolSize`, `ConnectTimeout` via `SqlConnectionStringBuilder`) and EF Core context pool (`poolSize`). Both must be configured for high-concurrency APIs.
- **Added**: `CreateExecutionStrategy()` pattern — mandatory wrapper for manual transactions when EF retry-on-failure policies are active. Without it, `BeginTransactionAsync` throws `InvalidOperationException`. Explains why EF needs to control the execution scope.
- **Updated**: Keyless entities section expanded with `ToView(null)` explanation (critical — without it, migrations generate a spurious `CREATE TABLE`), file naming convention (`Result` suffix), and clear distinction from domain entities and DTOs.

### v2.6 — 2026-04-09
- **Fixed (W-09)**: Removed simple `.Include(o => o.Customer)` example from Eager Loading section — agent already knows the basic syntax. Kept only `ThenInclude` (non-trivial chaining) and the WRONG NullRef example.

### v2.5 — 2026-03-28
- Fixed false statement in `DeleteAsync` XML doc (stub entity / silent delete warning).
- Removed non-portable path reference to `rules-to-skills/Standardized_NET_Rules.md`.

### v2.4 — 2026-03-28
- Reverted `FirstOrDefaultAsync` → `FirstOrDefault` inside `EF.CompileAsyncQuery` (async dispatch handled by infrastructure).
- Fixed: `RollbackAsync` on committed transaction throws `InvalidOperationException` — not a no-op.
- Fixed: Removed `ValueGeneratedOnUpdate()` from `UserConfiguration` — conflicts with `AuditInterceptor`.

### v2.3 — 2026-03-28
- Replaced broken multi-context transaction block with explicit warning + Outbox Pattern recommendation.
- Fixed `EntityState.Added` removed from `UpdatedAt` assignment in `AuditInterceptor`.
- Added `DeleteAsync` XML doc contract.

### v2.2 — 2026-03-28
- Removed full repository code example and duplicate SP example.
- Removed Commands block (standard dotnet ef CLI).

### v2.1 — 2026-03-23
- Generalized all context/domain names; removed proprietary references.

---

## Resources

- **Performance**: See [../performance/SKILL.md](../performance/SKILL.md)
- **Responses (layer contract)**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Security**: See [../security/SKILL.md](../security/SKILL.md)
