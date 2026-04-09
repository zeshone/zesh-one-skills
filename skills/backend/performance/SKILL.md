---
name: net8-apirest-performance
description: >
  Performance, concurrency and resilience rules for ASP.NET Core 8 REST APIs covering EF Core optimization, async patterns, Polly resilience pipelines, and parallelism best practices.
  Trigger: When optimizing API performance, configuring EF Core, implementing resilience with Polly, or working with concurrent operations in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Configuring EF Core for high-concurrency scenarios
- Implementing resilience pipelines with Polly (retry, circuit breaker, timeout)
- Parallelizing independent async operations
- Optimizing database queries

---

## Critical Patterns

### Circuit Breaker — Per-Request, No Global

The circuit breaker must be **per-request scoped**, not global. A global circuit breaker generates false positives under concurrency — a slow request opens the circuit for all others.

```csharp
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
            MinimumThroughput = 10,          // prevents false positives
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

    public async Task<PaymentResult> ProcessAsync(PaymentRequest request) =>
        await _pipeline.ExecuteAsync(async ct =>
        {
            var response = await _httpClient.PostAsJsonAsync("/payments", request, ct);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<PaymentResult>(ct);
        });
}
```

### Kestrel — Production Configuration

```csharp
builder.WebHost.ConfigureKestrel(options =>
{
    options.Limits.MaxConcurrentConnections = 1000;
    options.Limits.MaxRequestBodySize = 10 * 1024 * 1024; // 10 MB
    options.Limits.RequestHeadersTimeout = TimeSpan.FromSeconds(15);
});
```

### Response Caching — For stable data

```csharp
// Output caching (ASP.NET Core 7+) — preferred over ResponseCache
builder.Services.AddOutputCache(options =>
{
    options.AddPolicy("products", b => b.Expire(TimeSpan.FromMinutes(5)));
});
app.UseOutputCache();

[HttpGet]
[OutputCache(PolicyName = "products")]
public async Task<IActionResult> GetProducts() { ... }
```

---

## Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| Global Circuit Breaker | False positives under concurrency — a slow request affects everyone |
| `AspNetCoreRateLimit` | Concurrency conflicts with MemoryCache — use built-in rate limiter |
| `MemoryCache` for shared state | Concurrency interference — use distributed cache or correct scoping |

---

## Resources

- **General conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **Security / Rate Limiting**: See [../security/SKILL.md](../security/SKILL.md)
- **Data access / DbContextPool / AsNoTracking / Compiled Queries**: See [`../dataaccess/SKILL.md`](../dataaccess/SKILL.md)

---

## Changelog

### v1.2 — 2026-03-28
- **Removed**: `AddDbContextPool` and `AsNoTracking` (already covered in `dataaccess/SKILL.md`)
- **Removed**: Rate limiting (already covered in `security/SKILL.md`)
- **Removed**: GUIDs for IDs (already in `security/SKILL.md`)
- **Removed**: Compiled Queries example (already in `dataaccess/SKILL.md`)
- **Kept**: Circuit breaker per-request (design decision, not standard doc), Kestrel config, response caching
- **Added**: Emphasis on the reason for per-request circuit breaker (false positives under concurrency)

### v1.1 — 2026-03-24
- Fixed: `GetConnectionString("Default")` → `"DefaultConnection"`
- Added: Cross-refs to `security`, `dataaccess`, `requests`
