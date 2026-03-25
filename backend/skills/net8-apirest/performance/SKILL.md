---
name: net8-apirest-performance
description: >
  Performance, concurrency and resilience rules for ASP.NET Core 8 REST APIs covering EF Core optimization, async patterns, Polly resilience pipelines, and parallelism best practices.
  Trigger: When optimizing API performance, configuring EF Core, implementing resilience with Polly, or working with concurrent operations in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.1"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Configuring EF Core for high-concurrency scenarios
- Implementing Polly resilience pipelines (retry, circuit breaker, timeout)
- Parallelizing independent async operations
- Adding rate limiting to the API
- Reviewing or optimizing database queries

---

## Critical Patterns

### EF Core — DbContextPool (Mandatory)

**Never** use `AddDbContext`. Always use `AddDbContextPool` to reuse context instances across requests and reduce allocation overhead:

```csharp
// CORRECT
builder.Services.AddDbContextPool<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")),
    poolSize: 128);

// WRONG — creates a new instance per request with higher overhead
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")));
```

> **Production config**: This snippet shows pool sizing only. For robust connection setup (null-check, `SqlConnectionStringBuilder`, MARS) → see [`../dataaccess/SKILL.md`](../dataaccess/SKILL.md)

### EF Core — AsNoTracking for Read-Only Queries (Mandatory)

Every query that does NOT need to modify data must use `.AsNoTracking()`:

```csharp
// CORRECT — read-only query, no change tracking overhead
var users = await _context.Users
    .Where(u => u.IsActive)
    .AsNoTracking()
    .ToListAsync();

// For queries with includes that may return duplicates
var orders = await _context.Orders
    .Include(o => o.Items)
    .AsNoTrackingWithIdentityResolution()
    .ToListAsync();

// WRONG — change tracking enabled unnecessarily
var users = await _context.Users.Where(u => u.IsActive).ToListAsync();
```

### EF Core — Compiled Queries for Hot Paths

For frequently executed queries, use compiled queries to bypass LINQ translation overhead:

```csharp
private static readonly Func<AppDbContext, Guid, Task<User?>> _getUserById =
    EF.CompileAsyncQuery((AppDbContext ctx, Guid id) =>
        ctx.Users.AsNoTracking().FirstOrDefault(u => u.Id == id));

public async Task<User?> GetByIdAsync(Guid id) =>
    await _getUserById(_context, id);
```

### Parallel Independent Operations — Task.WhenAll

For multiple independent async operations, execute them in parallel. **Never** use sequential `await` for independent calls:

```csharp
// CORRECT — parallel execution
var (users, orders, products) = await (
    _userRepository.GetAllAsync(),
    _orderRepository.GetByUserAsync(userId),
    _productRepository.GetActiveAsync()
).WhenAll();

// Or explicit Task.WhenAll
var userTask = _userRepository.GetAllAsync();
var orderTask = _orderRepository.GetByUserAsync(userId);
await Task.WhenAll(userTask, orderTask);
var users = await userTask;
var orders = await orderTask;

// WRONG — sequential, wastes time
var users = await _userRepository.GetAllAsync();
var orders = await _orderRepository.GetByUserAsync(userId);
```

### Polly — Resilience Pipeline (Per-Request, Not Global)

Use Polly v8 `ResiliencePipelineBuilder`. Circuit breakers must be **per-request scoped**, not global, to prevent false positives under concurrency:

```csharp
// Program.cs
builder.Services.AddResiliencePipeline("external-api", pipeline =>
{
    pipeline
        .AddRetry(new RetryStrategyOptions
        {
            ShouldHandle = new PredicateBuilder().Handle<HttpRequestException>(),
            MaxRetryAttempts = 3,
            Delay = TimeSpan.FromMilliseconds(500),
            BackoffType = DelayBackoffType.Exponential
        })
        .AddCircuitBreaker(new CircuitBreakerStrategyOptions
        {
            ShouldHandle = new PredicateBuilder().Handle<HttpRequestException>(),
            FailureRatio = 0.5,
            MinimumThroughput = 10,          // prevent false positives
            SamplingDuration = TimeSpan.FromSeconds(30),
            BreakDuration = TimeSpan.FromSeconds(15)
        })
        .AddTimeout(TimeSpan.FromSeconds(10));
});
```

Usage in service via injection:
```csharp
public class ExternalPaymentService
{
    private readonly ResiliencePipeline _pipeline;
    private readonly HttpClient _httpClient;

    public ExternalPaymentService(ResiliencePipelineProvider<string> provider, HttpClient httpClient)
    {
        _pipeline = provider.GetPipeline("external-api");
        _httpClient = httpClient;
    }

    public async Task<PaymentResult> ProcessAsync(PaymentRequest request)
    {
        return await _pipeline.ExecuteAsync(async ct =>
        {
            var response = await _httpClient.PostAsJsonAsync("/payments", request, ct);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<PaymentResult>(ct);
        });
    }
}
```

### Rate Limiting — ASP.NET Core Built-in (Not AspNetCoreRateLimit)

Use ASP.NET Core 7+ built-in rate limiting. **Never** use `AspNetCoreRateLimit` (MemoryCache dependency causes concurrency issues):

```csharp
builder.Services.AddRateLimiter(options =>
{
    options.AddFixedWindowLimiter("api", limiterOptions =>
    {
        limiterOptions.PermitLimit = 100;
        limiterOptions.Window = TimeSpan.FromMinutes(1);
        limiterOptions.QueueProcessingOrder = QueueProcessingOrder.OldestFirst;
        limiterOptions.QueueLimit = 10;
    });
    options.RejectionStatusCode = StatusCodes.Status429TooManyRequests;
});

// In pipeline
app.UseRateLimiter();
```

> **Algorithm selection**: `performance` covers queue/window tuning. For choosing the right algorithm per scenario (SlidingWindow for auth, FixedWindow for API, TokenBucket for spikes) → see [`../security/SKILL.md`](../security/SKILL.md)

Apply to specific endpoints:
```csharp
[HttpPost]
[EnableRateLimiting("api")]
public async Task<IActionResult> Create([FromBody] CreateUserRequest request) { ... }
```

### GUIDs for Public IDs (Security + Performance)

Use GUIDs for all publicly exposed resource identifiers. Avoid sequential `int` IDs on public endpoints:

```csharp
public class User
{
    public Guid Id { get; set; } = Guid.NewGuid();
    // ...
}
```

> GUIDs prevent enumeration and reduce clustered index fragmentation vs sequential IDs. For full OWASP API1 (BOLA/IDOR) coverage → see [`../security/SKILL.md`](../security/SKILL.md)

### Response Caching

For endpoints with stable, infrequently changing data:

```csharp
[HttpGet]
[ResponseCache(Duration = 60, VaryByHeader = "Accept")]
public async Task<IActionResult> GetProducts() { ... }
```

For output caching (ASP.NET Core 7+):
```csharp
builder.Services.AddOutputCache(options =>
{
    options.AddPolicy("products", b => b.Expire(TimeSpan.FromMinutes(5)));
});

app.UseOutputCache();

[HttpGet]
[OutputCache(PolicyName = "products")]
public async Task<IActionResult> GetProducts() { ... }
```

### Anti-Patterns

| Anti-pattern | Problem | Fix |
|---|---|---|
| `AddDbContext` | No pooling, high allocation | Use `AddDbContextPool` |
| No `.AsNoTracking()` on reads | Memory waste, slow | Always add on read queries |
| Global Circuit Breaker | False positives under concurrency | Per-request pipeline |
| `AspNetCoreRateLimit` | MemoryCache concurrency conflicts | Use built-in rate limiter |
| Sequential `await` for independent calls | Wasted latency | Use `Task.WhenAll` |
| `MemoryCache` for shared state | Concurrency interference | Use distributed cache or proper scoping |
| Stateless helpers as `Scoped` | Unnecessary allocations | Register as `Singleton` |

---

## Code Examples

### Kestrel — Production Configuration

```csharp
builder.WebHost.ConfigureKestrel(options =>
{
    options.Limits.MaxConcurrentConnections = 1000;
    options.Limits.MaxRequestBodySize = 10 * 1024 * 1024; // 10 MB
    options.Limits.RequestHeadersTimeout = TimeSpan.FromSeconds(15);
});
```

### Parallel Repository Calls

```csharp
public async Task<DashboardDto> GetDashboardAsync(Guid userId)
{
    var userTask = _userRepository.GetByIdAsync(userId);
    var ordersTask = _orderRepository.GetRecentByUserAsync(userId, limit: 10);
    var statsTask = _statsRepository.GetUserStatsAsync(userId);
    await Task.WhenAll(userTask, ordersTask, statsTask);
    return new DashboardDto
    {
        User = (await userTask).ToDto(),
        RecentOrders = (await ordersTask).ToDtoList(),
        Stats = (await statsTask).ToDto()
    };
}
```

---

## Commands

```bash
# Add Polly v8 packages
dotnet add package Polly
dotnet add package Microsoft.Extensions.Http.Resilience

# Add output caching
dotnet add package Microsoft.AspNetCore.OutputCaching
```

---

## Resources

- **Standards**: See [../../../../rules-to-skills/Standardized_NET_Rules.md](../../../../rules-to-skills/Standardized_NET_Rules.md)
- **General conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **Security**: See [../security/SKILL.md](../security/SKILL.md)
- **Data access**: See [`../dataaccess/SKILL.md`](../dataaccess/SKILL.md)
- **Requests (pagination, rate limiting context)**: See [`../requests/SKILL.md`](../requests/SKILL.md)

---

## Changelog

### v1.1 — 2026-03-24
- **Fixed**: `GetConnectionString("Default")` → `"DefaultConnection"` — aligned with `dataaccess` v2.1 canonical key
- **Added**: Delegation note to `dataaccess` for robust connection setup (null-check, SqlConnectionStringBuilder)
- **Added**: Cross-ref to `security` in Rate Limiting section for algorithm decision table
- **Added**: Cross-ref to `security` in GUIDs section for OWASP API1 coverage
- **Added**: `dataaccess` and `requests` to `## Resources` cross-references
- **Updated**: GUIDs callout — split performance concern (index fragmentation) from security concern (BOLA)
- **Updated**: Frontmatter version `"1.0"` → `"1.1"`
