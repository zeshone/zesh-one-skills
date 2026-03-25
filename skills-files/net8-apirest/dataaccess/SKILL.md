---
name: net8-apirest-dataaccess
description: >
  Entity Framework Core data access standards for ASP.NET Core 8 REST APIs. Covers DbContext setup, entity configuration with Fluent API, repository pattern, audit automation via interceptors, soft delete with global query filters, raw SQL security, transactions, and all anti-patterns confirmed in the real codebase audit.
  Trigger: When working with EF Core, repositories, DbContext, raw SQL stored procedures, entity models, migrations, transactions, or audit/soft-delete patterns in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.1"
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
// CORRECT — simple include
var orders = await _context.Orders
    .Include(o => o.Customer)
    .AsNoTracking()
    .ToListAsync(ct);

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
        builder.Property(u => u.UpdatedAt).ValueGeneratedOnUpdate();

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
            if (entry.State is EntityState.Added or EntityState.Modified)
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
    Task DeleteAsync(Guid id, CancellationToken ct = default);
}
```

**Layer contract — never violate**:

| Layer | Returns | Never |
|---|---|---|
| Repository | Domain entities, `null`, empty `List<T>` | DTOs, `ResponseDTO<T>` |
| Service | DTOs (after mapping) | Raw entities |
| Controller | `ResponseDTO<T>` wrapped in `IActionResult` | Raw entities |

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

### Keyless Entities — Existing Pattern (Maintain)

Used for SPs and view projections — already established in existing contexts.

```csharp
// Register — in OnModelCreating or IEntityTypeConfiguration
modelBuilder.Entity<PolicySummaryResult>().HasNoKey();

// Query
var results = await _context.Set<PolicySummaryResult>()
    .FromSqlRaw("EXEC sp_GetPolicySummary @AgentId", new SqlParameter("@AgentId", agentId))
    .AsNoTracking()
    .ToListAsync(ct);
```

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

**Multi-context transaction — generic pattern** (maintain when needed):
```csharp
await using var primaryTx = await _primaryContext.Database.BeginTransactionAsync(ct);
await using var secondaryTx = await _secondaryContext.Database.BeginTransactionAsync(ct);
try
{
    await _primaryContext.SaveChangesAsync(ct);
    await _secondaryContext.SaveChangesAsync(ct);
    await primaryTx.CommitAsync(ct);
    await secondaryTx.CommitAsync(ct);
}
catch
{
    await primaryTx.RollbackAsync(ct);
    await secondaryTx.RollbackAsync(ct);
    throw;
}
```

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

*Anti-patterns confirmed in production legacy codebases.*

---

## Code Examples

### Complete Repository — All Patterns Applied

```csharp
public class UserRepository : IUserRepository
{
    private readonly AppDbContext _context;

    // Compiled query for hot path (auth / session checks)
    private static readonly Func<AppDbContext, Guid, Task<User?>> _getById =
        EF.CompileAsyncQuery((AppDbContext ctx, Guid id) =>
            ctx.Users.AsNoTracking().FirstOrDefault(u => u.Id == id));

    public UserRepository(AppDbContext context)
    {
        _context = context;
    }

    public Task<User?> GetByIdAsync(Guid id, CancellationToken ct = default) =>
        _getById(_context, id);

    public async Task<List<User>> GetAllActiveAsync(CancellationToken ct = default) =>
        await _context.Users
            .Where(u => u.IsActive)
            .AsNoTracking()
            .ToListAsync(ct);

    public async Task<bool> ExistsByEmailAsync(string email, CancellationToken ct = default) =>
        await _context.Users.AnyAsync(u => u.Email == email, ct);

    public async Task AddAsync(User user, CancellationToken ct = default)
    {
        _context.Users.Add(user);
        await _context.SaveChangesAsync(ct);
    }

    public async Task UpdateAsync(User user, CancellationToken ct = default)
    {
        _context.Users.Update(user);
        await _context.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(Guid id, CancellationToken ct = default)
    {
        var user = await _context.Users.FindAsync([id], ct);
        if (user is null) return;
        _context.Users.Remove(user); // soft delete intercepted by SaveChanges override if IsDeleted pattern active
        await _context.SaveChangesAsync(ct);
    }
}
```

### SP Result — Keyless Entity (Safe Pattern)

```csharp
// Keyless model
public class PolicySummaryResult
{
    public int PolicyId { get; set; }
    public string HolderName { get; set; } = string.Empty;
    public decimal Premium { get; set; }
}

// Context: modelBuilder.Entity<PolicySummaryResult>().HasNoKey();

public async Task<List<PolicySummaryResult>> GetPolicySummaryAsync(int agentId, CancellationToken ct)
{
    var param = new SqlParameter("@AgentId", agentId);
    return await _context.Set<PolicySummaryResult>()
        .FromSqlRaw("EXEC sp_GetPolicySummary @AgentId", param)
        .AsNoTracking()
        .ToListAsync(ct);
}
```

---

## Commands

```bash
# Add migration (descriptive name mandatory)
dotnet ef migrations add Add_Users_Table --project src/YourProject.Infrastructure --startup-project src/YourProject.Api

# Apply to DB
dotnet ef database update --project src/YourProject.Infrastructure --startup-project src/YourProject.Api

# Undo last migration (before applying to DB)
dotnet ef migrations remove --project src/YourProject.Infrastructure

# Generate idempotent SQL script for production deployment
dotnet ef migrations script --output migration.sql --idempotent --project src/YourProject.Infrastructure
```

---

## Changelog

### v2.1 — 2026-03-23
- Generalized all DbContext names to `AppDbContext`, `InventoryDbContext`, `ReportingDbContext`; removed `CPostalesContext`
- Replaced proprietary configuration singleton with standard `builder.Configuration.GetConnectionString("DefaultConnection")` / `IOptions<T>`
- Rewrote advisory block for legacy contexts to be fully generic
- Converted anti-patterns table column from "Confirmed location" to "Where to look" with generic descriptors
- Replaced domain-specific examples (`Cotizaciones`, `MedicalServices`, `Policies`) with e-commerce domain (`Order`, `Customer`, `Product`)
- Updated migration examples and folder references from `Main_Database_Library` to `src/YourProject.Infrastructure`
- Removed broken reference to non-existent `Reporte EFCORE.md`
- Bumped version to 2.1

---

## Resources

- **Performance**: See [../performance/SKILL.md](../performance/SKILL.md)
- **Responses (layer contract)**: See [../responses/SKILL.md](../responses/SKILL.md)
- **Security**: See [../security/SKILL.md](../security/SKILL.md)
- **Standards**: See [../../rules-to-skills/Standardized_NET_Rules.md](../../rules-to-skills/Standardized_NET_Rules.md)
