---
name: net8-apirest-performance
description: >
  Performance, concurrency and resilience rules for ASP.NET Core 8 REST APIs covering EF Core optimization, async patterns, Polly resilience pipelines, and parallelism best practices.
  Trigger: When optimizing API performance, configuring EF Core, implementing resilience with Polly, or working with concurrent operations in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.3"
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

### Ratio-Based Circuit Breaker — Production Pattern

The count-based circuit breaker (N consecutive failures = open) is fragile at low traffic: 1 failure out of 1 request = 100% failure rate = circuit opens. Use the **ratio-based** variant with a minimum throughput threshold.

```csharp
// IHttpClientFactory registration — production recipe
builder.Services.AddHttpClient("ExternalApi", client =>
{
    client.BaseAddress = new Uri(configuration["ExternalApi:BaseUrl"]!);
    client.Timeout = TimeSpan.FromSeconds(30);
})
.ConfigurePrimaryHttpMessageHandler(() => new SocketsHttpHandler
{
    // Connection pool at the TCP level — independent from IHttpClientFactory's handler lifetime
    PooledConnectionLifetime = TimeSpan.FromMinutes(2),  // recycle connections to avoid DNS stale
    MaxConnectionsPerServer = 100,
    EnableMultipleHttp2Connections = true
})
.AddPolicyHandler(HttpPolicyExtensions
    .HandleTransientHttpError()
    .WaitAndRetryAsync(
        retryCount: 3,
        sleepDurationProvider: attempt => TimeSpan.FromSeconds(Math.Pow(2, attempt)), // 2s, 4s, 8s
        onRetry: (outcome, timespan, attempt, _) =>
            Log.Warning("Retry {Attempt} after {Delay}s — {Reason}",
                attempt, timespan.TotalSeconds, outcome.Exception?.Message ?? outcome.Result.StatusCode.ToString())))
.AddPolicyHandler(HttpPolicyExtensions
    .HandleTransientHttpError()
    .AdvancedCircuitBreakerAsync(
        failureThreshold: 0.5,                        // open when 50%+ of requests fail
        samplingDuration: TimeSpan.FromSeconds(30),   // measured over a 30-second window
        minimumThroughput: 20,                        // requires at least 20 requests before opening
        durationOfBreak: TimeSpan.FromSeconds(15),    // stay open for 15s before half-opening
        onBreak: (outcome, duration) =>
            Log.Error("Circuit OPEN for {Duration}s — {Reason}",
                duration.TotalSeconds, outcome.Exception?.Message ?? outcome.Result.StatusCode.ToString()),
        onReset: () => Log.Information("Circuit CLOSED — service recovered"),
        onHalfOpen: () => Log.Warning("Circuit HALF-OPEN — probing service")));
```

> **`PooledConnectionLifetime`**: without this, `HttpClient` holds TCP connections indefinitely. DNS changes (failovers, load balancer updates) are invisible until the connection is recycled. 2 minutes is a safe balance between connection reuse efficiency and DNS freshness.
>
> **`minimumThroughput`**: the most important parameter. Without it, a circuit breaker on a low-traffic endpoint opens on the first transient failure. Set to at least 10–20 requests — enough to establish a statistically meaningful failure ratio.
>
> **Polly v7 vs v8**: The snippet above uses Polly v7 (`HttpPolicyExtensions`). Polly v8 uses `AddResiliencePipeline` (see the count-based example above). Both are valid — choose based on your Polly version.

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

### v1.3 — 2026-04-09
- **Added**: Ratio-based circuit breaker production recipe — `AdvancedCircuitBreakerAsync` with `failureThreshold`, `samplingDuration`, `minimumThroughput`, and `onBreak`/`onReset`/`onHalfOpen` callbacks. Explains why `minimumThroughput` is critical to prevent false positives at low traffic.
- **Added**: `SocketsHttpHandler` configuration — `PooledConnectionLifetime`, `MaxConnectionsPerServer`, `EnableMultipleHttp2Connections`. Documents why `PooledConnectionLifetime` is needed for DNS freshness. Pattern derived from production microservices calling external insurance APIs.
- **Added**: Polly v7 vs v8 note.

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
